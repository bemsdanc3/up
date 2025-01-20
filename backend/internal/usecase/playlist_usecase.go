package usecase

import (
	"backend/internal/entities"
	"backend/internal/repository"
)

type PlaylistUseCase interface {
	CreatePlaylist(playlist *entities.Playlist) error
	GetPlaylistByID(playlistID int) (entities.Playlist, error)
	DeletePlaylistByID(playlistID int) error
	EditPlaylistByID(playlist *entities.Playlist, playlistID int) error
}

type playlistUseCase struct {
	repo repository.PlaylistRepository
}

func NewPlaylistUseCase(repo repository.PlaylistRepository) PlaylistUseCase {
	return &playlistUseCase{repo: repo}
}

func (u *playlistUseCase) CreatePlaylist(playlist *entities.Playlist) error {
	return u.repo.CreatePlaylist(playlist)
}

func (u *playlistUseCase) GetPlaylistByID(playlistID int) (entities.Playlist, error) {
	return u.repo.GetPlaylistByID(playlistID)
}

func (u *playlistUseCase) DeletePlaylistByID(playlistID int) error {
	return u.repo.DeletePlaylistByID(playlistID)
}

func (u *playlistUseCase) EditPlaylistByID(playlist *entities.Playlist, playlistID int) error {
	return u.repo.EditPlaylistByID(playlist, playlistID)
}
