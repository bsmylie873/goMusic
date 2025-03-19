package models

type Artist struct {
	Id          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Nationality string `json:"nationality"`
	BirthDate   string `json:"birth_date"`
	Age         int    `json:"age"`
	Alive       bool   `json:"alive"`
	SexId       *int   `json:"sex_id,omitempty"`
	TitleId     *int   `json:"title_id,omitempty"`
	BandId      *int   `json:"band_id,omitempty"`
}
