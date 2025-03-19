package viewModels

import (
	"database/sql"
	"goMusic/db"
	"goMusic/models"
)

type ArtistViewModel struct {
	Id          *int   `json:"id,omitempty"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Nationality string `json:"nationality"`
	BirthDate   string `json:"birth_date"`
	Age         int    `json:"age"`
	Alive       bool   `json:"alive"`
	Sex         string `json:"sex"`
	Title       string `json:"title"`
}

func GetArtistViewModels(artists []models.Artist) ([]ArtistViewModel, error) {
	result := make([]ArtistViewModel, 0, len(artists))

	for _, artist := range artists {
		vm := ArtistViewModel{
			Id:          &artist.Id,
			FirstName:   artist.FirstName,
			LastName:    artist.LastName,
			Nationality: artist.Nationality,
			BirthDate:   artist.BirthDate,
			Age:         artist.Age,
			Alive:       artist.Alive,
		}

		if artist.SexId != nil {
			var sexName string
			err := db.DB.QueryRow("SELECT name FROM sexes WHERE id = ?", *artist.SexId).Scan(&sexName)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
			if err != sql.ErrNoRows {
				vm.Sex = sexName
			}
		}

		if artist.TitleId != nil {
			var titleName string
			err := db.DB.QueryRow("SELECT name FROM titles WHERE id = ?", *artist.TitleId).Scan(&titleName)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
			if err != sql.ErrNoRows {
				vm.Title = titleName
			}
		}

		result = append(result, vm)
	}

	return result, nil
}

func GetArtistViewModel(artist models.Artist) (ArtistViewModel, error) {
	vm := ArtistViewModel{
		Id:          &artist.Id,
		FirstName:   artist.FirstName,
		LastName:    artist.LastName,
		Nationality: artist.Nationality,
		BirthDate:   artist.BirthDate,
		Age:         artist.Age,
		Alive:       artist.Alive,
	}

	if artist.SexId != nil {
		var sexName string
		err := db.DB.QueryRow("SELECT name FROM sexes WHERE id = ?", *artist.SexId).Scan(&sexName)
		if err != nil && err != sql.ErrNoRows {
			return ArtistViewModel{}, err
		}
		if err != sql.ErrNoRows {
			vm.Sex = sexName
		}
	}

	if artist.TitleId != nil {
		var titleName string
		err := db.DB.QueryRow("SELECT name FROM titles WHERE id = ?", *artist.TitleId).Scan(&titleName)
		if err != nil && err != sql.ErrNoRows {
			return ArtistViewModel{}, err
		}
		if err != sql.ErrNoRows {
			vm.Title = titleName
		}
	}

	return vm, nil
}
