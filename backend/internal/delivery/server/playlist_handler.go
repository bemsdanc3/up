package server

import (
	"backend/internal/entities"
	"backend/internal/usecase"
	"backend/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
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

func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	// Ограничение на размер файла (например, 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	// Получение файла из `form-data`
	file, handler, err := r.FormFile("cover")
	if err != nil {
		log.Printf("error retrieving the file: %v", err)
		http.Error(w, "error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "./uploads/cover-playlist/"
	playlistCoverPath, err := utils.UploadFile(file, handler, uploadDir)
	if err != nil {
		log.Printf("error uploading track file: %v", err)
		http.Error(w, "failed to upload track", http.StatusInternalServerError)
		return
	}
	playlistCoverLink := "http://localhost:8080/uploads/cover-playlist/" + filepath.Base(playlistCoverPath)

	authorID, err := getUserIDFromCookie(r)
	if err != nil {
		log.Printf("error retrieving user ID: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var newPlaylist entities.Playlist
	newPlaylist.Title = r.FormValue("title")
	newPlaylist.Description = r.FormValue("description")
	newPlaylist.AuthorID = authorID
	isPublicStr := r.FormValue("is_public")
	isPublic := false // Значение по умолчанию

	if isPublicStr == "true" {
		isPublic = true
	} else if isPublicStr == "false" {
		isPublic = false
	}
	newPlaylist.IsPublic = isPublic
	newPlaylist.Cover = playlistCoverLink

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

	playlistID, err := getIDFromRouter(r)
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

	playlistID, err := getIDFromRouter(r)
	if err != nil {
		log.Printf("invalid playlist id: %v", err)
		http.Error(w, "invalid route parameter", http.StatusBadRequest)
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

func (h *PlaylistHandler) GetAllPlaylistsByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserIDFromCookie(r)
	if err != nil {
		log.Printf("Error extracting userID: %v", err)
		http.Error(w, "invalid id in cookie", http.StatusBadRequest)
		return
	}
	log.Printf("Extracted userID: %d", userID)

	playlists, err := h.useCase.GetAllPlaylistsByUserID(userID)
	if err != nil {
		log.Printf("Cant get playlists by user id: %v", err)
		http.Error(w, "error fetching playlists by user id", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlists)
}

func (h *PlaylistHandler) AddTrackToPlaylist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		PlaylistID int `json:"playlist_id,omitempty"`
		TrackID    int `json:"track_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("error reading body: %v", err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if requestBody.PlaylistID <= 0 || requestBody.TrackID <= 0 {
		http.Error(w, "someone of ID is less then zero", http.StatusBadRequest)
		return
	}

	if err := h.useCase.AddTrackToPlaylist(requestBody.PlaylistID, requestBody.TrackID); err != nil {
		log.Printf("error adding track to playlist: %v", err)
		http.Error(w, "failed to add track to playlist", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "track added to playlist"})
}

func (h *PlaylistHandler) GetAllPlaylistsByAuthorID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authorID, err := getIDFromRouter(r)
	if err != nil {
		log.Printf("Error getting ID from route: %v", err)
		http.Error(w, "invalid id in route", http.StatusBadRequest)
		return
	}

	playlists, err := h.useCase.GetAllPlaylistsByAuthorID(authorID)
	if err != nil {
		log.Printf("Cant get playlists by author id: %v", err)
		http.Error(w, "error fetching playlists by author id", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlists)
}

func (h *TrackHandler) IsTrackInPlaylist(w http.ResponseWriter, r *http.Request) {

	userID, err := getUserIDFromCookie(r)
	if err != nil {
		log.Printf("unauthorized: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var requestBody struct {
		TrackID int `json:"track_id"`
	}
	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("error reading request body: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	favoritePlaylist, err := h.playlistUsecase.GetFavoritePlaylist(userID)
	if err != nil {
		log.Printf("error getting favorite playlist: %v", err)
		http.Error(w, "error fetching favorite playlist", http.StatusBadRequest)
		return
	}

	isInPlaylist, err := h.playlistUsecase.IsTrackInPlaylist(favoritePlaylist.ID, requestBody.TrackID)
	if err != nil {
		log.Printf("error checking if track is in playlist: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"is_in_playlist": isInPlaylist})
}
