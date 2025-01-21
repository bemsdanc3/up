package usecase

import (
	"backend/internal/entities"
	"backend/internal/repository"
)

type AlbumUsecase interface {
	CreateAlbum(album *entities.Album) error
}

type albumUsecase struct {
	repo repository.AlbumRepository
}

func NewAlbumUsecase(repo repository.AlbumRepository) AlbumUsecase {
	return &albumUsecase{
		repo: repo,
	}
}

func (u *albumUsecase) CreateAlbum(album *entities.Album) error {
	return u.repo.CreateAlbum(album)
}
