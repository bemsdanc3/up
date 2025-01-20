package entities

import "time"

type Playlist struct {
	ID           int        `json:"id,omitempty"`
	Title        string     `json:"title,omitempty"`
	CreationDate *time.Time `json:"creation_date,omitempty"`
	AuthorID     int        `json:"author_id,omitempty"`
	Cover        string     `json:"cover,omitempty"`
	Description  string     `json:"description,omitempty"`
	IsPublic     bool       `json:"ispublic,omitempty"`
	Tracks       []Track    `json:"tracks,omitempty"`
}

type PlaylistTracks struct {
	ID         int        `json:"id,omitempty"`
	TrackID    int        `json:"track_id,omitempty"`
	AddDate    *time.Time `json:"add_date,omitempty"`
	PlaylistID int        `json:"playlist_id,omitempty"`
}
