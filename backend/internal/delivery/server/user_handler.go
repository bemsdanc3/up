package server

import (
	"backend/internal/entities"
	"backend/internal/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type UserHandler struct {
	usecase         usecase.UserUseCase
	playlistUsecase usecase.PlaylistUseCase
}

func NewUserHandler(usecase usecase.UserUseCase, playlistUsecase usecase.PlaylistUseCase) *UserHandler {
	return &UserHandler{
		usecase:         usecase,
		playlistUsecase: playlistUsecase,
	}
}

func generateTokens(user entities.User) (accessToken, refreshToken string, err error) {
	// Генерация access-токена (жизнь короткая, например, 1 час)
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour).Unix(), // 1 час
	}
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJWT.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Генерация refresh-токена (жизнь долгая, например, 7 дней)
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 дней
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJWT.SignedString(jwtSecret)
	return
}

func getUserIDFromCookie(r *http.Request) (int, error) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		log.Printf("Error getting cookie: %v", err)
		return 0, errors.New("user_id cookie not found")
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Printf("invalid user_id format: %v", err)
		return 0, errors.New("invalid user_id format in cookie")
	}
	return userID, nil
}

func getIDFromRouter(r *http.Request) (int, error) {
	vars := mux.Vars(r)

	idStr, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("parameter 'id' not found in route")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid 'id' format in route: %s", idStr)
	}

	if id <= 0 {
		return 0, fmt.Errorf("id cant be less then zero")
	}

	return id, nil
}

func ValidatePassword(password string) error {
	// Проверяем длину пароля
	if len(password) < 4 || len(password) > 16 {
		return errors.New("password must be between 4 and 16 characters")
	}

	// Проверяем наличие запрещенных символов
	forbiddenSymbols := "*&{}|+"
	for _, char := range password {
		if strings.ContainsRune(forbiddenSymbols, char) {
			return errors.New("password contains forbidden symbols")
		}
	}

	// Проверяем наличие заглавных букв и цифр
	var hasUpper, hasDigit bool
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	return nil
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}

	var newUser entities.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := ValidatePassword(newUser.Pass); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.usecase.CreateUser(&newUser); err != nil {
		if err.Error() == "email is already in use" {
			http.Error(w, "email is already in use", http.StatusBadRequest)
		} else if err.Error() == "login is already in use" {
			http.Error(w, "login is already in use", http.StatusBadRequest)
		} else {
			http.Error(w, "error creating user", http.StatusInternalServerError)
			log.Printf("something went wrong: %v", err)
		}
		return
	}
	favoritePlaylist := &entities.Playlist{
		Title:       "Любимые треки",
		AuthorID:    newUser.ID,
		Cover:       "http://localhost:8080/uploads/cover-playlist/rnqa7yhv4il71.webp",
		Description: "Автоматически созданный плейлист для любимых треков",
		IsPublic:    false,
	}
	if err := h.playlistUsecase.CreatePlaylist(favoritePlaylist); err != nil {
		log.Printf("Error creating favorite playlist: %v", err)
	}

	token, refreshToken, err := generateTokens(newUser)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	//userId := newUser.ID

	http.SetCookie(w, &http.Cookie{ // токен в куки
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{ // токен в куки
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(newUser.ID),
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "user created successfully",
		"role":    "user",
		"user_id": strconv.Itoa(newUser.ID),
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method", http.StatusBadRequest)
		return
	}

	var input entities.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Erorr decoding request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("login attempt for email: %v", input.Email)
	user, err := h.usecase.GetUserByEmail(input.Email)
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		http.Error(w, "Invalid email", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(input.Pass)); err != nil {
		log.Printf("Password missmatch for user: %s", input.Email)
		log.Printf("stored hash: %s, Input password: %s", user.Pass, input.Pass)
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}

	token, refreshToken, err := generateTokens(*user)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{ // токен в куки
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{ // токен в куки
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(user.ID),
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "login successful",
		"role":    user.Role,
		"user_id": user.ID,
	})
}

func (h *UserHandler) GetUserDetailsByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "wrong method", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromCookie(r)
	if err != nil {
		log.Printf("error getting user id from cookie: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.usecase.GetUserDetails(userID)
	if err != nil {
		log.Printf("error getting user details: %v", err)
		http.Error(w, "error fetching user details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	var updatedUser entities.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		log.Printf("cant read updated user body: %v", err)
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromCookie(r)
	if err != nil {
		log.Printf("Cant get user ID from cookie: %v", err)
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	if userID <= 0 {
		http.Error(w, "user ID can be less then zero", http.StatusBadRequest)
		return
	}

	err = h.usecase.UpdateUserProfile(&updatedUser, userID)
	if err != nil {
		log.Printf("cant update user: %v", err)
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "user updated successfully"})
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getIDFromRouter(r)
	if err != nil {
		log.Printf("cant get user ID from cookie: %v", err)
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.GetUserProfile(userID)
	if err != nil {
		log.Printf("error getting user by ID: %v", err)
		http.Error(w, "error fetching user by ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "No refresh token provided", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		http.Error(w, "Invalid refresh token claims", http.StatusUnauthorized)
		return
	}

	userID := int(claims["user_id"].(float64))

	user := entities.User{ID: userID}
	accessToken, refreshToken, err := generateTokens(user)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tokens refreshed successfully"})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Время в прошлом
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.usecase.GetAllUsers()
	if err != nil {
		log.Printf("error getting users: %v", err)
		http.Error(w, "error getting users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Role   string `json:"role,omitempty"`
		UserID int    `json:"user_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("error reading req body: %v", err)
		return
	}

	if requestBody.UserID == 0 || requestBody.Role == "" {
		http.Error(w, "Invalid user ID or role", http.StatusBadRequest)
		return
	}

	if err := h.usecase.UpdateUsersRole(requestBody.UserID, requestBody.Role); err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Error updating user role: %v", err)
			http.Error(w, "Error updating user role", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "user role updated successfully",
	})
}
