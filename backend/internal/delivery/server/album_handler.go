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

	// Парсинг формы для обработки файлов
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("error parsing multipart form: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// Получение файла обложки из запроса
	file, handler, err := r.FormFile("cover")
	if err != nil {
		log.Printf("error retrieving file: %v", err)
		http.Error(w, "invalid cover file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "./uploads/covers-albums/"
	coverPath, err := utils.UploadFile(file, handler, uploadDir)
	if err != nil {
		log.Printf("error uploading cover file: %v", err)
		http.Error(w, "failed to upload cover", http.StatusInternalServerError)
		return
	}
	// Генерация URL для обложки
	coverLink := "http://localhost:8080/uploads/covers-albums/" + filepath.Base(coverPath)

	// Чтение остальных данных альбома из формы
	var newAlbum entities.Album
	newAlbum.Title = r.FormValue("title")
	newAlbum.TypeOf = r.FormValue("type_of")
	newAlbum.Label = r.FormValue("label")
	newAlbum.Cover = coverLink // Установить путь к обложке

	// Получение ID автора из cookie
	authorID, err := getUserIDFromCookie(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	newAlbum.AuthorID = authorID

	// Создание альбома через usecase
	if err = h.usecase.CreateAlbum(&newAlbum); err != nil {
		log.Printf("error creating album: %v", err)
		http.Error(w, "failed to create album", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "album created successfully",
		"cover":   coverPath,
	})
}

func (h *AlbumHandler) DeleteAlbumByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	albumID, err := getIDFromRouter(r)
	if err != nil {
		http.Error(w, "invalid id in route", http.StatusBadRequest)
		return
	}

	if err = h.usecase.DeleteAlbumByID(albumID); err != nil {
		log.Printf("cant delete album: %v", err)
		http.Error(w, "invalid input", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "album deleted successfully"})
}

func (h *AlbumHandler) EditAlbumByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	albumID, err := getIDFromRouter(r)
	if err != nil {
		http.Error(w, "invalid id in route", http.StatusBadRequest)
		return
	}

	var updatedAlbum entities.Album
	if err = json.NewDecoder(r.Body).Decode(&updatedAlbum); err != nil {
		log.Printf("cant read album body: %v", err)
		http.Error(w, "invalid album body", http.StatusBadRequest)
		return
	}

	if err = h.usecase.EditAlbumByID(&updatedAlbum, albumID); err != nil {
		log.Printf("error updating album: %v", err)
		http.Error(w, "invalid input", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "album updated successfully"})
}
