package repository

import (
	"backend/internal/entities"
	"database/sql"
	"log"
)

type AlbumRepository interface {
	CreateAlbum(album *entities.Album) error
	DeleteAlbumByID(albumID int) error
	EditAlbumByID(album *entities.Album, albumID int) error
}

type albumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) AlbumRepository {
	return &albumRepository{
		db: db,
	}
}

func (r *albumRepository) CreateAlbum(album *entities.Album) error {
	query := `INSERT INTO albums (type_of, title, cover, label) VALUES ($1,$2,$3,$4)`
	_, err := r.db.Exec(query, album.TypeOf, album.Title, album.Cover, album.Label)
	if err != nil {
		log.Printf("error creating album on database lvl: %v", err)
		return err
	}
	return err
}

func (r *albumRepository) DeleteAlbumByID(albumID int) error {
	query := `DELETE FROM albums WHERE id = $1`
	_, err := r.db.Exec(query, albumID)
	if err != nil {
		log.Printf("error deleting album on database lvl: %v", err)
		return err
	}
	return err
}

func (r *albumRepository) EditAlbumByID(album *entities.Album, albumID int) error {
	query := `UPDATE albums SET title = $1, cover = $2 WHERE id = $3`
	_, err := r.db.Exec(query, album.Title, album.Cover, albumID)
	if err != nil {
		log.Printf("error updating album on database lvl: %v", err)
		return err
	}
	return err
}
