package models

type Album struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	ArtistID *int    `json:"artist_id,omitempty"`
	BandID   *int    `json:"band_id,omitempty"`
}
