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
	"time"
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

func generateJWT(user entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(jwtSecret)
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
		Cover:       "",
		Description: "Автоматически созданный плейлист для любимых треков",
		IsPublic:    false,
	}
	if err := h.playlistUsecase.CreatePlaylist(favoritePlaylist); err != nil {
		log.Printf("Error creating favorite playlist: %v", err)
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "user created successfully"})
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

	token, err := generateJWT(*user)
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

	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(user.ID),
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	})

	json.NewEncoder(w).Encode(map[string]string{"message": "login successful"})
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
