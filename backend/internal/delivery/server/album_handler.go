package server

import (
	"backend/internal/entities"
	"backend/internal/usecase"
	"encoding/json"
	"log"
	"net/http"
)

type AlbumHandler struct {
	usecase usecase.AlbumUsecase
}

func NewAlbumHandler(usecase usecase.AlbumUsecase) *AlbumHandler {
	return &AlbumHandler{usecase: usecase}
}

func (h *AlbumHandler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newAlbum entities.Album
	if err := json.NewDecoder(r.Body).Decode(&newAlbum); err != nil {
		log.Printf("cant read album body: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	authorID, err := getUserIDFromCookie(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	newAlbum.AuthorID = authorID
	if err = h.usecase.CreateAlbum(&newAlbum); err != nil {
		log.Printf("error creating album: %v", err)
		http.Error(w, "failed to create album", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "album created successfully"})
}
