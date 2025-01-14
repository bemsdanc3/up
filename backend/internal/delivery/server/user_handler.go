package server

import (
	"backend/internal/entities"
	"backend/internal/usecase"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
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

func generateJWT(user entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString(jwtSecret)
}
