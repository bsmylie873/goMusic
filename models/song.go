package models

type Song struct {
	Id       int     `json:"id"`
	Title    string  `json:"title" validate:"required,min=1,max=1000"`
	Length   int     `json:"length" validate:"required,min=0"`
	Price    float64 `json:"price" validate:"required,min=0"`
	AlbumId  *int    `json:"album_id,omitempty"`
	ArtistId *int    `json:"artist_id,omitempty"`
	BandId   *int    `json:"band_id,omitempty"`
}
