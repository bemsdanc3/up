package repository

import (
	"backend/internal/entities"
	"database/sql"
	"log"
)

type AlbumRepository interface {
	CreateAlbum(album *entities.Album) error
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
	query := `INSERT INTO albums (type_of, cover, label) VALUES ($1,$2,$3)`
	_, err := r.db.Exec(query, album.TypeOf, album.Cover, album.Label)
	if err != nil {
		log.Printf("error creating album on database lvl: %v", err)
		return err
	}
	return err
}
