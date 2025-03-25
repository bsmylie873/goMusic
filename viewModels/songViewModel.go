package viewModels

import (
	"database/sql"
	"goMusic/db"
	"goMusic/models"
)

type SongViewModel struct {
	ID     *int               `json:"id"`
	Title  string             `json:"title"`
	Length int                `json:"length"`
	Price  float64            `json:"price"`
	Albums *[]AlbumViewModel  `json:"albums,omitempty"`
	Artist *[]ArtistViewModel `json:"artist,omitempty"`
	Band   *[]BandViewModel   `json:"band,omitempty"`
}

func GetSongViewModels(songs []models.Song) ([]SongViewModel, error) {
	result := make([]SongViewModel, 0, len(songs))

	for _, song := range songs {
		vm := SongViewModel{
			ID:     &song.Id,
			Title:  song.Title,
			Length: song.Length,
			Price:  song.Price,
		}

		albumRows, err := db.DB.Query(`
            SELECT a.id, a.title
            FROM albums a
            JOIN album_songs sa ON a.id = sa.album_id
            WHERE sa.song_id = ?`, song.Id)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if err != sql.ErrNoRows {
			var albums []AlbumViewModel
			defer albumRows.Close()

			for albumRows.Next() {
				var album AlbumViewModel
				var id int
				if err := albumRows.Scan(&id, &album.Title); err != nil {
					return nil, err
				}
				album.Id = &id
				albums = append(albums, album)
			}

			if len(albums) > 0 {
				vm.Albums = &albums
			}
		}

		artistRows, err := db.DB.Query(`
            SELECT a.id, a.first_name, a.last_name
            FROM artists a
            JOIN artist_songs sa ON a.id = sa.artist_id
            WHERE sa.song_id = ?`, song.Id)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if err != sql.ErrNoRows {
			var artists []ArtistViewModel
			defer artistRows.Close()

			for artistRows.Next() {
				var artist ArtistViewModel
				var id int
				if err := artistRows.Scan(&id, &artist.FirstName, &artist.LastName); err != nil {
					return nil, err
				}
				artist.Id = &id
				artists = append(artists, artist)
			}

			if len(artists) > 0 {
				vm.Artist = &artists
			}
		}

		bandRows, err := db.DB.Query(`
            SELECT b.id, b.name
            FROM bands b
            JOIN band_songs sb ON b.id = sb.band_id
            WHERE sb.song_id = ?`, song.Id)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		if err != sql.ErrNoRows {
			var bands []BandViewModel
			defer bandRows.Close()

			for bandRows.Next() {
				var band BandViewModel
				var id int
				if err := bandRows.Scan(&id, &band.Name); err != nil {
					return nil, err
				}
				band.Id = &id
				bands = append(bands, band)
			}

			if len(bands) > 0 {
				vm.Band = &bands
			}
		}

		result = append(result, vm)
	}

	return result, nil
}

func GetSongViewModel(song models.Song) (SongViewModel, error) {
	vm := SongViewModel{
		ID:     &song.Id,
		Title:  song.Title,
		Length: song.Length,
		Price:  song.Price,
	}

	albumRows, err := db.DB.Query(`
        SELECT a.id, a.title, a.price, a.artist_id, a.band_id
        FROM albums a
        JOIN album_songs sa ON a.id = sa.album_id
        WHERE sa.song_id = ?`, song.Id)
	if err != nil && err != sql.ErrNoRows {
		return SongViewModel{}, err
	}

	if err != sql.ErrNoRows {
		var albums []AlbumViewModel
		defer albumRows.Close()

		for albumRows.Next() {
			var album models.Album
			if err := albumRows.Scan(&album.ID, &album.Title, &album.Price, &album.ArtistId, &album.BandId); err != nil {
				return SongViewModel{}, err
			}

			albumVM, err := GetAlbumViewModel(album)
			if err != nil {
				return SongViewModel{}, err
			}

			albums = append(albums, albumVM)
		}

		if len(albums) > 0 {
			vm.Albums = &albums
		}
	}

	artistRows, err := db.DB.Query(`
        SELECT a.*
        FROM artists a
        JOIN artist_songs sa ON a.id = sa.artist_id
        WHERE sa.song_id = ?`, song.Id)
	if err != nil && err != sql.ErrNoRows {
		return SongViewModel{}, err
	}

	if err != sql.ErrNoRows {
		var artists []ArtistViewModel
		defer artistRows.Close()

		for artistRows.Next() {
			var artist models.Artist
			if err := artistRows.Scan(&artist.Id, &artist.FirstName, &artist.LastName,
				&artist.Nationality, &artist.BirthDate, &artist.Age,
				&artist.Alive, &artist.SexId, &artist.TitleId, &artist.BandId); err != nil {
				return SongViewModel{}, err
			}

			artistVM, err := GetArtistViewModel(artist)
			if err != nil {
				return SongViewModel{}, err
			}

			artists = append(artists, artistVM)
		}

		if len(artists) > 0 {
			vm.Artist = &artists
		}
	}

	bandRows, err := db.DB.Query(`
        SELECT b.*
        FROM bands b
        JOIN band_songs sb ON b.id = sb.band_id
        WHERE sb.song_id = ?`, song.Id)
	if err != nil && err != sql.ErrNoRows {
		return SongViewModel{}, err
	}

	if err != sql.ErrNoRows {
		var bands []BandViewModel
		defer bandRows.Close()

		for bandRows.Next() {
			var band models.Band
			if err := bandRows.Scan(&band.Id, &band.Name, &band.Nationality,
				&band.NumberOfMembers, &band.DateFormed,
				&band.Age, &band.Active); err != nil {
				return SongViewModel{}, err
			}

			bandVM, err := GetBandViewModel(band)
			if err != nil {
				return SongViewModel{}, err
			}

			bands = append(bands, bandVM)
		}

		if len(bands) > 0 {
			vm.Band = &bands
		}
	}

	return vm, nil
}
