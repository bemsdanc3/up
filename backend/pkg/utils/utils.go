package utils

import (
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadFile(file multipart.File, handler *multipart.FileHeader, uploadDir string) (string, error) {
	// Генерация уникального имени файла
	uniqueID := uuid.New().String()
	fileExt := filepath.Ext(handler.Filename)
	uniqueFileName := uniqueID + fileExt

	// Создание директории, если она не существует
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("error creating upload directory: %v", err)
		return "", err
	}

	// Генерация полного пути к файлу
	filePath := filepath.Join(uploadDir, uniqueFileName)

	// Открытие файла для записи
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("error creating file: %v", err)
		return "", err
	}
	defer dst.Close()

	// Копирование содержимого загруженного файла в целевой файл
	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("error saving file: %v", err)
		return "", err
	}

	// Возвращаем путь к файлу
	return filePath, nil
}
