package models

type Album struct {
	ID       int     `json:"id" validate:"required,min=0"`
	Title    string  `json:"title" validate:"required,min=1,max=100"`
	Price    float64 `json:"price" validate:"required,min=0"`
	ArtistId *int    `json:"artist_id,omitempty"`
	BandId   *int    `json:"band_id,omitempty"`
}
