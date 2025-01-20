package entities

type User struct {
	ID    int     `json:"id,omitempty"`
	Login string  `json:"login,omitempty"`
	Pass  string  `json:"pass,omitempty"`
	Email string  `json:"email,omitempty"`
	Role  string  `json:"role,omitempty"`
	PFP   *string `json:"pfp,omitempty"`
}

type UserDetails struct {
	User         User
	Subscription []int
	Playlists    []Playlist
}
