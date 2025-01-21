package server

import (
	"backend/internal/entities"
	"backend/internal/usecase"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TrackHandler struct {
	usecase usecase.TrackUsecase
}

func NewTrackHandler(usecase usecase.TrackUsecase) *TrackHandler {
	return &TrackHandler{usecase: usecase}
}

func getTrackIDFromRouter(r *http.Request) (int, error) {
	vars := mux.Vars(r)

	idStr, ok := vars["track_id"]
	if !ok {
		return 0, fmt.Errorf("parameter 'id' not found in route")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid 'id' format in route: %s", idStr)
	}

	if id <= 0 {
		return 0, fmt.Errorf("id can be less then zero")
	}

	return id, nil
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

	// Генерация уникального имени файла
	uniqueID := uuid.New().String()
	fileExt := filepath.Ext(handler.Filename)
	uniqueFileName := uniqueID + fileExt

	// Генерация пути для сохранения файла
	uploadDir := "./uploads/tracks/"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("error creating upload directory: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(uploadDir, uniqueFileName)

	// Сохранение файла на сервер
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("error saving file: %v", err)
		http.Error(w, "error saving file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := outFile.ReadFrom(file); err != nil {
		log.Printf("error writing to file: %v", err)
		http.Error(w, "error writing to file", http.StatusInternalServerError)
		return
	}

	// Генерация ссылки на файл
	trackLink := "localhost:8080/uploads/tracks/" + uniqueFileName

	// Извлечение остальных данных трека из формы
	duration, _ := strconv.Atoi(r.FormValue("duration"))
	albumID, _ := strconv.Atoi(r.FormValue("album_id"))
	genreID, _ := strconv.Atoi(r.FormValue("genre_id"))

	// Создание объекта трека
	newTrack := entities.Track{
		Duration:    duration,
		AlbumID:     albumID,
		GenreID:     genreID,
		TrackLink:   trackLink,
		ListenCount: 0, // Начальное значение
	}

	// Сохранение трека через usecase
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

	trackID, err := getTrackIDFromRouter(r)
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

	trackID, err := getTrackIDFromRouter(r)
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
