package entities

type Track struct {
	ID          int    `json:"id,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	Title       string `json:"title,omitempty"`
	AlbumID     int    `json:"album_id,omitempty"`
	GenreID     int    `json:"genre_id,omitempty"`
	TrackLink   string `json:"track_link,omitempty"`
	ListenCount int    `json:"listen_count,omitempty"`
	Cover       string `json:"cover,omitempty"`
	AuthorID    int    `json:"author_id,omitempty"`
	AuthorLogin string `json:"author_login,omitempty"`
	AlbumTitle  string `json:"album_title,omitempty"`
}
