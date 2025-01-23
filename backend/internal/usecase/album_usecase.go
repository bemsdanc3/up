package usecase

import (
	"backend/internal/entities"
	"backend/internal/repository"
)

type AlbumUsecase interface {
	CreateAlbum(album *entities.Album) error
	DeleteAlbumByID(albumID int) error
	EditAlbumByID(album *entities.Album, albumID int) error
	GetAlbumsByAuthorID(authorID int) ([]entities.Album, error)
	GetAlbumByID(albumID int) (*entities.Album, error)
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

func (u *albumUsecase) DeleteAlbumByID(albumID int) error {
	return u.repo.DeleteAlbumByID(albumID)
}

func (u *albumUsecase) EditAlbumByID(album *entities.Album, albumID int) error {
	return u.repo.EditAlbumByID(album, albumID)
}

func (u *albumUsecase) GetAlbumsByAuthorID(authorID int) ([]entities.Album, error) {
	return u.repo.GetAlbumsByAuthorID(authorID)
}

func (u *albumUsecase) GetAlbumByID(albumID int) (*entities.Album, error) {
	return u.repo.GetAlbumByID(albumID)
}
