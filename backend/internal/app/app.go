package app

import (
	"backend/internal/delivery/server"
	"backend/internal/repository"
	"backend/internal/usecase"
	"backend/middleware"
	"backend/pkg/db"
	"github.com/gorilla/mux"
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
	playlistRepo := repository.NewPlaylistRepository(database)
	trackRepo := repository.NewTrackRepository(database)

	// инициализация юзкейсов
	userUsecase := usecase.NewUserUsecse(userRepo)
	playlistUsecase := usecase.NewPlaylistUseCase(playlistRepo)
	trackUsecase := usecase.NewTrackUsecase(trackRepo)

	// инициализация обработчиков
	userHandler := server.NewUserHandler(userUsecase)
	playlistHandler := server.NewPlaylistHandler(playlistUsecase)
	trackHandler := server.NewTrackHandler(trackUsecase)

	// урлы
	r := mux.NewRouter()
	r.HandleFunc("/register", userHandler.Register)
	r.HandleFunc("/login", userHandler.Login)

	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)

	//пользовательские запросы
	authRouter.HandleFunc("/user/my/profile", userHandler.GetUserDetailsByID).Methods(http.MethodGet)
	authRouter.HandleFunc("/user/my/profile/edit", userHandler.UpdateUserProfile).Methods(http.MethodPatch)
	authRouter.HandleFunc("/users/profile/{user_id}", userHandler.GetUserProfile).Methods(http.MethodGet)

	//запросы с плейлистами
	authRouter.HandleFunc("/playlists/create", playlistHandler.CreatePlaylist).Methods(http.MethodPost)
	authRouter.HandleFunc("/playlists/view/{id}", playlistHandler.GetPlaylistByID).Methods(http.MethodGet)
	authRouter.HandleFunc("/playlists/delete", playlistHandler.DeletePlaylistByID).Methods(http.MethodDelete)
	authRouter.HandleFunc("/playlists/update/{id}", playlistHandler.EditPlaylistByID).Methods(http.MethodPatch)
	authRouter.HandleFunc("/playlists/all", playlistHandler.GetAllPlaylistsByUserID).Methods(http.MethodGet)

	//запросы с треками
	authRouter.HandleFunc("/track/create", trackHandler.CreateTrack).Methods(http.MethodPost)
	authRouter.HandleFunc("/track/delete/{track_id}", trackHandler.DeleteTrackByID).Methods(http.MethodDelete)
	authRouter.HandleFunc("/track/view/{track_id}", trackHandler.GetTrackByID).Methods(http.MethodGet)

	//статическая папка для треков
	setupStaticFileServer(r)

	// Добавляем middleware для CORS
	handler := middleware.CorsMiddleware(r)

	// Запуск сервера
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupStaticFileServer(router *mux.Router) {
	// Настраиваем маршруты для статической папки
	staticDir := "./uploads"
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(staticDir))))
}
