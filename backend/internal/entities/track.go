package entities

type Track struct {
	ID          int    `json:"id,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	AlbumID     int    `json:"album_id,omitempty"`
	GenreID     int    `json:"genre_id,omitempty"`
	TrackLink   string `json:"track_link,omitempty"`
	ListenCount int    `json:"listen_count,omitempty"`
}
