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
	GetAlbumsByAuthorID(authorID int) ([]entities.Album, error)
	GetAlbumByID(albumID int) (*entities.Album, error)
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
	query := `INSERT INTO albums (type_of, title, cover, label, author_id) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.Exec(query, album.TypeOf, album.Title, album.Cover, album.Label, album.AuthorID)
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

func (r *albumRepository) GetAlbumsByAuthorID(authorID int) ([]entities.Album, error) {
	query := `SELECT id, title, cover FROM albums WHERE author_id = $1`
	rows, err := r.db.Query(query, authorID)
	if err != nil {
		log.Printf("error getting albums by author id on database lvl: %v", err)
		return nil, err
	}
	defer rows.Close()

	var albums []entities.Album
	for rows.Next() {
		var album entities.Album
		if err = rows.Scan(
			&album.ID,
			&album.Title,
			&album.Cover,
		); err != nil {
			log.Printf("error scanning row: %v", err)
			return nil, err
		}
		albums = append(albums, album)
	}
	if err = rows.Err(); err != nil {
		log.Printf("row iteration error: %v", err)
		return nil, err
	}

	return albums, nil
}

func (r *albumRepository) GetAlbumByID(albumID int) (*entities.Album, error) {
	query := `SELECT 
			a.title AS album_title,
			a.cover AS album_cover,
			a.type_of AS album_type,
			a.label AS album_label,
			a.publication_date AS publication_date,
			u.id AS author_id,
			u.login AS author_login
		FROM 
			albums a
		LEFT JOIN 
			users u ON a.author_id = u.id
		WHERE 
			a.id = $1`

	var album entities.Album
	err := r.db.QueryRow(query, albumID).Scan(
		&album.Title,
		&album.Cover,
		&album.TypeOf,
		&album.Label,
		&album.PublicationDate,
		&album.AuthorID,
		&album.AuthorLogin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("error fetching album by ID: %v", err)
		return nil, err
	}
	return &album, nil
}
