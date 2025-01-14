package app

import (
	"backend/internal/delivery/server"
	"backend/internal/repository"
	"backend/internal/usecase"
	"backend/middleware"
	"backend/pkg/db"
	"log"
	"net/http"
)

func Run() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// инициализация репозиториев
	userRepo := repository.NewUserRepository(database)

	// инициализация юзкейсов
	userUsecase := usecase.NewUserUsecse(userRepo)

	// инициализация обработчиков
	userHandler := server.NewUserHandler(userUsecase)

	// урлы
	mux := http.NewServeMux()
	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)

	// Добавляем middleware для CORS
	handler := middleware.CorsMiddleware(mux)

	// Запуск сервера
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
