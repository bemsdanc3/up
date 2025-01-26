package server

import (
	"backend/internal/entities"
	"backend/internal/usecase"
	"backend/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TrackHandler struct {
	usecase         usecase.TrackUsecase
	playlistUsecase usecase.PlaylistUseCase
}

func NewTrackHandler(usecase usecase.TrackUsecase, playlistUsecase usecase.PlaylistUseCase) *TrackHandler {
	return &TrackHandler{
		usecase:         usecase,
		playlistUsecase: playlistUsecase,
	}
}

func (h *TrackHandler) CreateTrack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ограничение на размер файла (например, 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	// Получение файла из `form-data`
	file, handler, err := r.FormFile("trackFile")
	if err != nil {
		log.Printf("error retrieving the file: %v", err)
		http.Error(w, "error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "./uploads/tracks/"
	trackPath, err := utils.UploadFile(file, handler, uploadDir)
	if err != nil {
		log.Printf("error uploading track file: %v", err)
		http.Error(w, "failed to upload track", http.StatusInternalServerError)
		return
	}
	trackLink := "http://localhost:8080/uploads/tracks/" + filepath.Base(trackPath)

	// Извлечение остальных данных трека из формы
	var newTrack entities.Track
	newTrack.Duration, err = strconv.Atoi(r.FormValue("duration"))
	newTrack.Title = r.FormValue("title")
	newTrack.AlbumID, err = strconv.Atoi(r.FormValue("album_id"))
	newTrack.GenreID, err = strconv.Atoi(r.FormValue("genre_id"))
	newTrack.ListenCount = 0
	newTrack.TrackLink = trackLink

	if err := h.usecase.CreateTrack(&newTrack); err != nil {
		log.Printf("error creating track: %v", err)
		http.Error(w, "error creating track", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	json.NewEncoder(w).Encode(map[string]string{"message": "track created successfully", "track_link": trackLink})
}

func (h *TrackHandler) DeleteTrackByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	trackID, err := getIDFromRouter(r)
	if err != nil {
		log.Printf("cant get track id from route: %v", err)
		http.Error(w, "invalid track id", http.StatusBadRequest)
		return
	}

	track, err := h.usecase.GetTrackLinkByTrackID(trackID)
	if err != nil {
		log.Printf("Cant get track link: %v", err)
		http.Error(w, "error fetching track link", http.StatusBadRequest)
		return
	}

	if err = h.usecase.DeleteTrackByID(trackID); err != nil {
		log.Printf("cant delete track: %v", err)
		http.Error(w, "failed to delete track", http.StatusInternalServerError)
		return
	}

	filePath := "." + strings.TrimPrefix(track.TrackLink, "localhost:8080")
	if err := os.Remove(filePath); err != nil {
		log.Printf("failed to delete track file: %v", err)
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "track deleted successfully"})
}

func (h *TrackHandler) GetTrackByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	trackID, err := getIDFromRouter(r)
	if err != nil {
		log.Printf("cant get track id from route: %v", err)
		http.Error(w, "invalid track id", http.StatusBadRequest)
		return
	}

	track, err := h.usecase.GetTrackByID(trackID)
	if err != nil {
		log.Printf("cant get track: %v", err)
		http.Error(w, "error fetching track", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(track)
}

func (h *TrackHandler) GetAllTracks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tracks, err := h.usecase.GetAllTracks()
	if err != nil {
		log.Printf("cant get tracks: %v", err)
		http.Error(w, "erorr fetching tracks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

func (h *TrackHandler) LikeTrack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
		http.Error(w, "error fertching favorite playlist", http.StatusBadRequest)
		return
	}

	err = h.playlistUsecase.AddTrackToFavorite(favoritePlaylist.ID, requestBody.TrackID)
	if err != nil {
		if err.Error() == "track already in favorite playlist" {
			http.Error(w, "track already in favorite playlist", http.StatusConflict)
		} else {
			log.Printf("error adding track to favorite playlist: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "track liked successfully"})
}

func (h *TrackHandler) GetRandomTrack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	track, err := h.usecase.GetRandomTrack()
	if err != nil {
		log.Printf("error retrieving random track: %v", err)
		http.Error(w, "error retrieving random track", http.StatusInternalServerError)
		return
	}
	if track == nil {
		http.Error(w, "no tracks available", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(track)
}

func (h *TrackHandler) GetTrackByAuthorID(w http.ResponseWriter, r *http.Request) {
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

	tracks, err := h.usecase.GetTrackByAuthorID(authorID)
	if err != nil {
		log.Printf("Cant get tracks by author id: %v", err)
		http.Error(w, "error fetching tracks by author id", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

func (h *TrackHandler) GetTrackByGenreID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	genreID, err := getIDFromRouter(r)
	if err != nil {
		log.Printf("Error getting ID from route: %v", err)
		http.Error(w, "invalid id in route", http.StatusBadRequest)
		return
	}

	track, err := h.usecase.GetRandomTrackByGenre(genreID)
	if err != nil {
		http.Error(w, "Error fetching random track", http.StatusInternalServerError)
		log.Printf("Error in handler: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(track)
}

func (h *TrackHandler) GetTracksByPlaylistID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	playlistID, err := getIDFromRouter(r)
	if err != nil {
		log.Printf("Error getting ID from route: %v", err)
		http.Error(w, "invalid id in route", http.StatusBadRequest)
		return
	}

	playlist, err := h.usecase.GetTracksByPlaylistID(playlistID)
	if err != nil {
		http.Error(w, "Error fetching tracks", http.StatusInternalServerError)
		log.Printf("Error in handler: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlist)
}

func (h *TrackHandler) GetTracksByAlbumID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	albumID, err := getIDFromRouter(r)
	if err != nil {
		log.Printf("Error getting ID from route: %v", err)
		http.Error(w, "invalid id in route", http.StatusBadRequest)
		return
	}

	album, err := h.usecase.GetTracksByAlbumID(albumID)
	if err != nil {
		http.Error(w, "Error fetching tracks", http.StatusInternalServerError)
		log.Printf("Error in handler: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}
