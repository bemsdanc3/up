package repository

import (
	"backend/internal/entities"
	"database/sql"
	"log"
)

type TrackRepository interface {
	CreateTrack(track *entities.Track) error
	DeleteTrackByID(trackID int) error
	GetTrackByID(trackID int) (*entities.Track, error)
	GetTrackLinkByTrackID(trackID int) (*entities.Track, error)
	GetAllTracks() ([]entities.Track, error)
	GetTrackCount() (int, error)
	GetRandomTrack(offset int) (*entities.Track, error)
}

type trackRepository struct {
	db *sql.DB
}

func NewTrackRepository(db *sql.DB) TrackRepository {
	return &trackRepository{
		db: db,
	}
}

func (r *trackRepository) CreateTrack(track *entities.Track) error {
	query := `INSERT INTO tracks (duration, title, album_id, genre_id, tracklink, listen_count) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.Exec(query, track.Duration, track.Title, track.AlbumID, track.GenreID, track.TrackLink, track.ListenCount)
	if err != nil {
		log.Printf("error creating track on database lvl: %v", err)
		return err
	}
	return err
}

func (r *trackRepository) DeleteTrackByID(trackID int) error {
	query := `DELETE FROM tracks WHERE id = $1`
	_, err := r.db.Exec(query, trackID)
	if err != nil {
		log.Printf("Cant delete track on database lvl: %v", err)
		return err
	}

	return err
}

func (r *trackRepository) GetTrackByID(trackID int) (*entities.Track, error) {
	query := `
		SELECT 
			t.duration, t.title, t.album_id, t.genre_id, t.tracklink, t.listen_count,
			a.cover, a.author_id, u.login
		FROM 
			tracks t
		LEFT JOIN 
			albums a ON t.album_id = a.id
		LEFT JOIN 
			users u ON a.author_id = u.id
		WHERE 
			t.id = $1
	`
	var track entities.Track
	err := r.db.QueryRow(query, trackID).Scan(
		&track.Duration,
		&track.Title,
		&track.AlbumID,
		&track.GenreID,
		&track.TrackLink,
		&track.ListenCount,
		&track.Cover,
		&track.AuthorID,
		&track.AuthorLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("error querying track by ID: %v", err)
		return nil, err
	}
	return &track, nil
}

func (r *trackRepository) GetTrackLinkByTrackID(trackID int) (*entities.Track, error) {
	query := `SELECT tracklink FROM tracks WHERE id = $1`
	var trackLink string
	err := r.db.QueryRow(query, trackID).Scan(&trackLink)
	if err != nil {
		log.Printf("Cant get tracklink from database: %v", err)
		return nil, err
	}
	return &entities.Track{
		TrackLink: trackLink,
	}, nil
}

func (r *trackRepository) GetAllTracks() ([]entities.Track, error) {
	query := `
		SELECT 
			t.id, t.duration, t.title, t.album_id, t.genre_id, t.tracklink, t.listen_count, 
			a.cover, a.author_id, u.login
		FROM 
			tracks t
		LEFT JOIN 
			albums a ON t.album_id = a.id
		LEFT JOIN 
			users u ON a.author_id = u.id;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []entities.Track
	for rows.Next() {
		var track entities.Track

		err = rows.Scan(
			&track.ID,
			&track.Duration,
			&track.Title,
			&track.AlbumID,
			&track.GenreID,
			&track.TrackLink,
			&track.ListenCount,
			&track.Cover,
			&track.AuthorID,
			&track.AuthorLogin,
		)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (r *trackRepository) GetTrackCount() (int, error) {
	var trackCount int
	err := r.db.QueryRow("SELECT COUNT(*) FROM tracks").Scan(&trackCount)
	if err != nil {
		return 0, err
	}
	return trackCount, nil
}

func (r *trackRepository) GetRandomTrack(offset int) (*entities.Track, error) {
	query := `
		SELECT 
			t.id, 
			t.title, 
			t.duration, 
			t.album_id, 
			t.genre_id, 
			t.tracklink, 
			t.listen_count, 
			a.cover AS album_cover, 
			a.author_id AS album_author_id, 
			u.login AS album_author_login
		FROM 
			tracks t
		LEFT JOIN 
			albums a ON t.album_id = a.id
		LEFT JOIN 
			users u ON a.author_id = u.id
		LIMIT 1 OFFSET $1
	`
	var track entities.Track
	err := r.db.QueryRow(query, offset).Scan(
		&track.ID,
		&track.Title,
		&track.Duration,
		&track.AlbumID,
		&track.GenreID,
		&track.TrackLink,
		&track.ListenCount,
		&track.Cover,
		&track.AuthorID,
		&track.AuthorLogin,
	)
	if err != nil {
		return nil, err
	}
	return &track, nil
}
