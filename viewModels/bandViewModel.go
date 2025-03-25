package viewModels

import (
	"goMusic/models"
)

type BandViewModel struct {
	Id              *int   `json:"id,omitempty"`
	Name            string `json:"name"`
	Nationality     string `json:"nationality"`
	NumberOfMembers int    `json:"number_of_members"`
	DateFormed      string `json:"date_formed"`
	Age             int    `json:"age"`
	Active          bool   `json:"active"`
}

type BasicBandViewModel struct {
	Id   *int   `json:"id,omitempty"`
	Name string `json:"name"`
}

func GetBandViewModels(bands []models.Band) ([]BandViewModel, error) {
	result := make([]BandViewModel, 0, len(bands))
	for _, band := range bands {
		vm := BandViewModel{
			Id:              &band.Id,
			Name:            band.Name,
			Nationality:     band.Nationality,
			NumberOfMembers: band.NumberOfMembers,
			DateFormed:      band.DateFormed,
			Age:             band.Age,
			Active:          band.Active,
		}
		result = append(result, vm)
	}
	return result, nil
}

func GetBandViewModel(band models.Band) (BandViewModel, error) {
	vm := BandViewModel{
		Id:              &band.Id,
		Name:            band.Name,
		Nationality:     band.Nationality,
		NumberOfMembers: band.NumberOfMembers,
		DateFormed:      band.DateFormed,
		Age:             band.Age,
		Active:          band.Active,
	}
	return vm, nil
}

func GetBasicBandViewModel(band models.Band) BasicBandViewModel {
	return BasicBandViewModel{
		Id:   &band.Id,
		Name: band.Name,
	}
}
