package models

type Song struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Length   int     `json:"length"`
	Price    float64 `json:"price"`
	AlbumID  *int    `json:"album_id,omitempty"`
	ArtistID *int    `json:"artist_id,omitempty"`
	BandID   *int    `json:"band_id,omitempty"`
}
