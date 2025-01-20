package server

import (
	"backend/internal/entities"
	"backend/internal/usecase"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type PlaylistHandler struct {
	useCase usecase.PlaylistUseCase
}

type RequestBody struct {
	PlaylistID int `json:"playlist_id,omitempty"`
}

func NewPlaylistHandler(useCase usecase.PlaylistUseCase) *PlaylistHandler {
	return &PlaylistHandler{useCase: useCase}
}

func getPlaylistIDFromRouter(r *http.Request) (int, error) {
	vars := mux.Vars(r)

	idStr, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("parameter 'id' not found in route")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid 'id' format in route: %s", idStr)
	}

	return id, nil
}

// ПОФИКСИТЬ ЧТОБЫ СОЗДАВАТЬ МОГ ТОЛЬКО АВТОРИЗИРОВАННЫЙ ПОЛЬЗОВАТЕЛЬ
func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	var newPlaylist entities.Playlist
	if err := json.NewDecoder(r.Body).Decode(&newPlaylist); err != nil {
		log.Printf("Error decoding playlist body in json: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.useCase.CreatePlaylist(&newPlaylist); err != nil {
		log.Printf("Error creating playlist: %v", err)
		http.Error(w, "error creating playlist", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "playlist created successfully"})
}

func (h *PlaylistHandler) GetPlaylistByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	playlistID, err := getPlaylistIDFromRouter(r)
	if err != nil {
		log.Printf("invalid playlist id: %v", err)
		http.Error(w, "invalid route parameter", http.StatusBadRequest)
		return
	}

	if playlistID <= 0 {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	playlist, err := h.useCase.GetPlaylistByID(playlistID)
	if err != nil {
		http.Error(w, "failed to get playlist details"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlist)
}

// ПОФИКСИТЬ ЧТОБЫ ПОЛЬЗОВАТЕЛИ МОГЛИ УДАЛЯТЬ ТОЛЬКО СВОИ ПЛЕЙЛИСТЫ
func (h *PlaylistHandler) DeletePlaylistByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody RequestBody

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("error fertching request body: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.PlaylistID <= 0 {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	err := h.useCase.DeletePlaylistByID(requestBody.PlaylistID)
	if err != nil {
		log.Printf("error deleting playlist: %v", err)
		http.Error(w, "failed to delete playlist", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "playlist deleted successfully"})
}

func (h *PlaylistHandler) EditPlaylistByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var updatedPlaylist entities.Playlist
	if err := json.NewDecoder(r.Body).Decode(&updatedPlaylist); err != nil {
		log.Printf("error fetching playlist: %v", err)
		http.Error(w, "Invalid playlist", http.StatusBadRequest)
		return
	}

	playlistID, err := getPlaylistIDFromRouter(r)
	if err != nil {
		log.Printf("invalid playlist id: %v", err)
		http.Error(w, "invalid route parameter", http.StatusBadRequest)
		return
	}

	if playlistID <= 0 {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	err = h.useCase.EditPlaylistByID(&updatedPlaylist, playlistID)
	if err != nil {
		log.Printf("cant edit playlist: %v", err)
		http.Error(w, "failed to edit playlist", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "playlist edited successfully"})
}
