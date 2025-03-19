package viewModels

type SongViewModel struct {
	ID     *int               `json:"id"`
	Title  string             `json:"title"`
	Length int                `json:"length"`
	Price  float64            `json:"price"`
	Albums *[]AlbumViewModel  `json:"albums,omitempty"`
	Artist *[]ArtistViewModel `json:"artist,omitempty"`
	Band   *[]BandViewModel   `json:"band,omitempty"`
}
