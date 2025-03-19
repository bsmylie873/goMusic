package viewModels

import (
	"database/sql"
	"goMusic/db"
	"goMusic/models"
)

type AlbumViewModel struct {
	Id     *int             `json:"id,omitempty"`
	Title  string           `json:"title"`
	Price  float64          `json:"price"`
	Artist *ArtistViewModel `json:"artist,omitempty"`
	Band   *BandViewModel   `json:"band,omitempty"`
	Songs  []SongViewModel  `json:"songs,omitempty"`
}

func GetAlbumViewModels(albums []models.Album) ([]AlbumViewModel, error) {
	result := make([]AlbumViewModel, 0, len(albums))

	for _, album := range albums {
		vm := AlbumViewModel{
			Id:    &album.ID,
			Title: album.Title,
			Price: album.Price,
			Songs: []SongViewModel{},
		}

		if album.ArtistID != nil {
			var artist ArtistViewModel
			var bandName sql.NullString
			err := db.DB.QueryRow(`
			SELECT 
			  a.first_name, 
			  a.last_name, 
			  a.nationality, 
			  a.birth_date, 
			  a.age, 
			  a.alive, 
			  s.name, 
			  t.name, 
			  b.name
			FROM artists a
			LEFT JOIN sexes s ON a.sex_id = s.id
			LEFT JOIN titles t ON a.title_id = t.id
			LEFT JOIN bands b ON a.band_id = b.id
			WHERE a.id = ?`, *album.ArtistID,
			).Scan(
				&artist.FirstName,
				&artist.LastName,
				&artist.Nationality,
				&artist.BirthDate,
				&artist.Age,
				&artist.Alive,
				&artist.Sex,
				&artist.Title,
				&bandName,
			)

			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}

			if err != sql.ErrNoRows {
				vm.Artist = &artist
			}
		}

		if album.BandID != nil {
			var band BandViewModel
			err := db.DB.QueryRow(
				"SELECT name, nationality, number_of_members, date_formed, age, active FROM bands WHERE id = ?",
				*album.BandID,
			).Scan(&band.Name, &band.Nationality, &band.NumberOfMembers, &band.DateFormed, &band.Age, &band.Active)

			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}

			if err != sql.ErrNoRows {
				vm.Band = &band
			}
		}

		result = append(result, vm)
	}

	return result, nil
}

func GetAlbumViewModel(album models.Album) (AlbumViewModel, error) {
	vm := AlbumViewModel{
		Id:    &album.ID,
		Title: album.Title,
		Price: album.Price,
	}

	if album.ArtistID != nil {
		var artist ArtistViewModel
		var bandName sql.NullString
		err := db.DB.QueryRow(`
		SELECT
			a.first_name,
			a.last_name,
			a.nationality,
			a.birth_date,
			a.age,
			a.alive,
			s.name,
			t.name,
			b.name
		FROM artists a
		LEFT JOIN sexes s ON a.sex_id = s.id
		LEFT JOIN titles t ON a.title_id = t.id
		LEFT JOIN bands b ON a.band_id = b.id
		WHERE a.id = ?`, *album.ArtistID,
		).Scan(
			&artist.FirstName,
			&artist.LastName,
			&artist.Nationality,
			&artist.BirthDate,
			&artist.Age,
			&artist.Alive,
			&artist.Sex,
			&artist.Title,
			&bandName,
		)

		if err != nil && err != sql.ErrNoRows {
			return AlbumViewModel{}, err
		}

		if err != sql.ErrNoRows {
			vm.Artist = &artist
		}
	}

	if album.BandID != nil {
		var band BandViewModel
		err := db.DB.QueryRow(
			"SELECT name, nationality, number_of_members, date_formed, age, active FROM bands WHERE id = ?",
			*album.BandID,
		).Scan(&band.Name, &band.Nationality, &band.NumberOfMembers, &band.DateFormed, &band.Age, &band.Active)

		if err != nil && err != sql.ErrNoRows {
			return AlbumViewModel{}, err
		}

		if err != sql.ErrNoRows {
			vm.Band = &band
		}
	}

	rows, err := db.DB.Query("SELECT id, title, length, price FROM songs WHERE album_id = ?", album.ID)
	if err != nil {
		return AlbumViewModel{}, err
	}
	defer rows.Close()

	vm.Songs = []SongViewModel{}
	for rows.Next() {
		var song SongViewModel
		var songID int

		if err := rows.Scan(&songID, &song.Title, &song.Length, &song.Price); err != nil {
			return AlbumViewModel{}, err
		}

		song.ID = &songID

		vm.Songs = append(vm.Songs, song)
	}

	if err := rows.Err(); err != nil {
		return AlbumViewModel{}, err
	}

	return vm, nil
}
