package entities

import "time"

type Album struct {
	ID              int        `json:"id,omitempty"`
	TypeOf          string     `json:"type_of,omitempty"`
	Cover           string     `json:"cover,omitempty"`
	AuthorID        int        `json:"author_id,omitempty"`
	Label           string     `json:"label,omitempty"`
	PublicationDate *time.Time `json:"publication_date,omitempty"`
}
