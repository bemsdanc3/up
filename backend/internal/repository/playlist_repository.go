package repository

import (
	"backend/internal/entities"
	"database/sql"
	"fmt"
	"log"
)

type PlaylistRepository interface {
	CreatePlaylist(playlist *entities.Playlist) error
	GetPlaylistByID(playlistID int) (entities.Playlist, error)
	DeletePlaylistByID(playlistID int) error
	EditPlaylistByID(playlist *entities.Playlist, playlistID int) error
}

type playlistRepository struct {
	db *sql.DB
}

func NewPlaylistRepository(db *sql.DB) PlaylistRepository {
	return &playlistRepository{db: db}
}

func (r *playlistRepository) CreatePlaylist(playlist *entities.Playlist) error {
	query := `INSERT INTO playlists (title, description, cover, ispublic) VALUES ($1,$2,$3,$4)`
	_, err := r.db.Exec(query, playlist.Title, playlist.Description, playlist.Cover, playlist.IsPublic)
	if err != nil {
		log.Printf("error creating playlist on database lvl: %v", err)
	}
	return err
}

func (r *playlistRepository) GetPlaylistByID(playlistID int) (entities.Playlist, error) {
	var playlist entities.Playlist

	queryPlaylist := `SELECT id, title, creation_date, author_id, cover, description, ispublic FROM playlists WHERE id = $1`

	err := r.db.QueryRow(queryPlaylist, playlistID).Scan(
		&playlist.ID,
		&playlist.Title,
		&playlist.CreationDate,
		&playlist.AuthorID,
		&playlist.Cover,
		&playlist.Description,
		&playlist.IsPublic,
	)
	if err != nil {
		log.Printf("error querying playlist: %v", err)
		return playlist, fmt.Errorf("error fetching playlist: %v", err)
	}

	queryTracks := `
		SELECT 
			t.id, t.duration, t.album_id, t.genre_id, t.tracklink, t.listen_count
		FROM 
			playlist_tracks pt
		JOIN 
			tracks t ON pt.track_id = t.id
		WHERE 
			pt.playlist_id = $1
	`

	rows, err := r.db.Query(queryTracks, playlistID)
	if err != nil {
		log.Printf("error querying tracks: %v", err)
		return playlist, fmt.Errorf("error fetching tracks for playlists: %v", err)
	}
	defer rows.Close()

	var tracks []entities.Track
	for rows.Next() {
		var track entities.Track
		if err := rows.Scan(
			&track.ID,
			&track.Duration,
			&track.AlbumID,
			&track.GenreID,
			&track.TrackLink,
			&track.ListenCount,
		); err != nil {
			log.Printf("error adding tracks to playlist: %v", err)
			return playlist, fmt.Errorf("error scanning track row: %v", err)
		}

		tracks = append(tracks, track)
	}

	playlist.Tracks = tracks

	return playlist, nil
}

func (r *playlistRepository) DeletePlaylistByID(playlistID int) error {
	query := `DELETE from playlists WHERE id = $1`
	_, err := r.db.Exec(query, playlistID)
	if err != nil {
		log.Printf("error deleting playlist on database lvl: %v", err)
		return err
	}
	return err

}

func (r *playlistRepository) EditPlaylistByID(playlist *entities.Playlist, playlistID int) error {
	query := `UPDATE playlists SET title = $1, cover = $2, description = $3, ispublic = $4 WHERE id = $5`
	_, err := r.db.Exec(query, playlist.Title, playlist.Cover, playlist.Description, playlist.IsPublic, playlistID)
	if err != nil {
		log.Printf("cant edit playlist on database lvl: %v", err)
		return err
	}
	return err
}
