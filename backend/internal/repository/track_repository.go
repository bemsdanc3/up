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
	query := `INSERT INTO tracks (duration, album_id, genre_id, tracklink, listen_count) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.Exec(query, track.Duration, track.AlbumID, track.GenreID, track.TrackLink, track.ListenCount)
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
	query := `SELECT duration, album_id, genre_id, tracklink, listen_count FROM tracks WHERE id = $1`
	var track entities.Track
	err := r.db.QueryRow(query, trackID).Scan(&track.Duration, &track.AlbumID, &track.GenreID, &track.TrackLink, &track.ListenCount)
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
