package models

type Artist struct {
	Id          int    `json:"id"`
	FirstName   string `json:"first_name" validate:"required,min=1,max=100"`
	LastName    string `json:"last_name" validate:"required,min=1,max=100"`
	Nationality string `json:"nationality" validate:"required,min=1,max=100"`
	BirthDate   string `json:"birth_date" validate:"required"`
	Age         int    `json:"age" validate:"required,min=0,max=150"`
	Alive       bool   `json:"alive"`
	SexId       *int   `json:"sex_id,omitempty" validate:"validSex"`
	TitleId     *int   `json:"title_id,omitempty" validate:"validTitle"`
	BandId      *int   `json:"band_id,omitempty"`
}
