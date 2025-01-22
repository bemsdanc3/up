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
	albumRepo := repository.NewAlbumRepository(database)

	// инициализация юзкейсов
	userUsecase := usecase.NewUserUsecse(userRepo)
	playlistUsecase := usecase.NewPlaylistUseCase(playlistRepo)
	trackUsecase := usecase.NewTrackUsecase(trackRepo)
	albumUsecase := usecase.NewAlbumUsecase(albumRepo)

	// инициализация обработчиков
	userHandler := server.NewUserHandler(userUsecase, playlistUsecase)
	playlistHandler := server.NewPlaylistHandler(playlistUsecase)
	trackHandler := server.NewTrackHandler(trackUsecase)
	albumHandler := server.NewAlbumHandler(albumUsecase)

	// урлы
	r := mux.NewRouter()
	r.HandleFunc("/register", userHandler.Register)
	r.HandleFunc("/login", userHandler.Login)

	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)

	//пользовательские запросы
	authRouter.HandleFunc("/user/my/profile", userHandler.GetUserDetailsByID).Methods(http.MethodGet)
	authRouter.HandleFunc("/user/my/profile/edit", userHandler.UpdateUserProfile).Methods(http.MethodPatch)
	authRouter.HandleFunc("/users/profile/{id}", userHandler.GetUserProfile).Methods(http.MethodGet)

	//запросы с плейлистами
	authRouter.HandleFunc("/playlists/create", playlistHandler.CreatePlaylist).Methods(http.MethodPost)
	authRouter.HandleFunc("/playlists/view/{id}", playlistHandler.GetPlaylistByID).Methods(http.MethodGet)
	authRouter.HandleFunc("/playlists/delete", playlistHandler.DeletePlaylistByID).Methods(http.MethodDelete)
	authRouter.HandleFunc("/playlists/update/{id}", playlistHandler.EditPlaylistByID).Methods(http.MethodPatch)
	authRouter.HandleFunc("/playlists/all", playlistHandler.GetAllPlaylistsByUserID).Methods(http.MethodGet)
	authRouter.HandleFunc("/playlists/add/track", playlistHandler.AddTrackToPlaylist).Methods(http.MethodPost)

	//запросы с треками
	authRouter.HandleFunc("/tracks/create", trackHandler.CreateTrack).Methods(http.MethodPost)
	authRouter.HandleFunc("/tracks/delete/{id}", trackHandler.DeleteTrackByID).Methods(http.MethodDelete)
	authRouter.HandleFunc("/tracks/view/{id}", trackHandler.GetTrackByID).Methods(http.MethodGet)
	authRouter.HandleFunc("/tracks/all", trackHandler.GetAllTracks).Methods(http.MethodGet)

	//запросы с альбомами
	authRouter.HandleFunc("/album/create", albumHandler.CreateAlbum).Methods(http.MethodPost)
	authRouter.HandleFunc("/album/delete/{id}", albumHandler.DeleteAlbumByID).Methods(http.MethodDelete)
	authRouter.HandleFunc("/album/update/{id}", albumHandler.EditAlbumByID).Methods(http.MethodPatch)
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
