package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}

	_, err = DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return err
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS sexes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	)`)

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS titles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	)`)

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS artists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		nationality TEXT NOT NULL,
		birth_date DATE NOT NULL,
		age INT NOT NULL,
		alive BOOLEAN NOT NULL,
		sex_id INT NOT NULL,
		title_id INT NOT NULL,
		band_id INT,
		FOREIGN KEY (sex_id) REFERENCES sexes(id),
	    FOREIGN KEY (title_id) REFERENCES titles(id),
	    FOREIGN KEY (band_id) REFERENCES bands(id)
	)`)

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS bands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		nationality TEXT,
		number_of_members INT NOT NULL,
		date_formed DATE NOT NULL,
		age INT,
		active BOOLEAN
	)`)

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS albums (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		price REAL NOT NULL,
	    artist_id INT,
	    band_id INT,
	    FOREIGN KEY (artist_id) REFERENCES artists(id),
	    FOREIGN KEY (band_id) REFERENCES bands(id)
	)`)

	_, err = DB.Exec(`
	 CREATE TABLE IF NOT EXISTS songs (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  title TEXT NOT NULL,
	  length INT NOT NULL,
	  price REAL NOT NULL,
	  album_id INT,
	  artist_id INT,
	  band_id INT,
	  FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE,
	  FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
	  FOREIGN KEY (band_id) REFERENCES bands(id) ON DELETE CASCADE
	)`)

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS album_songs (
		album_id INTEGER NOT NULL,
		song_id INTEGER NOT NULL,
		PRIMARY KEY (album_id, song_id),
		FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE,
		FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
	)`)

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS artist_songs (
		artist_id INTEGER NOT NULL,
		song_id INTEGER NOT NULL,
		PRIMARY KEY (artist_id, song_id),
		FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
		FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
	)`)

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS band_songs (
		band_id INTEGER NOT NULL,
		song_id INTEGER NOT NULL,
		PRIMARY KEY (band_id, song_id),
		FOREIGN KEY (band_id) REFERENCES bands(id) ON DELETE CASCADE,
		FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
	)`)

	return err
}

func SeedDB() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM sexes").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		sexes := []struct {
			ID   int
			Name string
		}{
			{1, "Male"},
			{2, "Female"},
			{3, "Non-binary"},
		}
		for _, s := range sexes {
			_, err := DB.Exec("INSERT INTO sexes (id, name) VALUES (?, ?)", s.ID, s.Name)
			if err != nil {
				return err
			}
		}

		titles := []struct {
			ID   int
			Name string
		}{
			{1, "Mr."},
			{2, "Mrs."},
			{3, "Ms."},
			{4, "Dr."},
			{5, "Prof."},
		}
		for _, t := range titles {
			_, err := DB.Exec("INSERT INTO titles (id, name) VALUES (?, ?)", t.ID, t.Name)
			if err != nil {
				return err
			}
		}

		bands := []struct {
			ID              int
			Name            string
			Nationality     string
			NumberOfMembers int
			DateFormed      string
			Age             int
			Active          bool
		}{
			{1, "Pink Floyd", "British", 5, "1965-01-01", 59, false},
			{2, "The Beatles", "British", 4, "1960-08-01", 64, false},
			{3, "Radiohead", "British", 5, "1985-09-01", 39, true},
		}
		for _, b := range bands {
			_, err := DB.Exec(
				"INSERT INTO bands (id, name, nationality, number_of_members, date_formed, age, active) VALUES (?, ?, ?, ?, ?, ?, ?)",
				b.ID, b.Name, b.Nationality, b.NumberOfMembers, b.DateFormed, b.Age, b.Active)
			if err != nil {
				return err
			}
		}

		artists := []struct {
			ID          int
			FirstName   string
			LastName    string
			Nationality string
			BirthDate   string
			Age         int
			Alive       bool
			SexID       int
			TitleID     int
			BandID      *int
		}{
			{1, "John", "Coltrane", "American", "1926-09-23", 40, false, 1, 1, nil},
			{2, "Gerry", "Mulligan", "American", "1927-04-06", 68, false, 1, 1, nil},
			{3, "Sarah", "Vaughan", "American", "1924-03-27", 66, false, 2, 2, nil},
			{4, "Thom", "Yorke", "British", "1968-10-07", 56, true, 1, 1, intPtr(3)},
		}
		for _, a := range artists {
			_, err := DB.Exec(
				"INSERT INTO artists (id, first_name, last_name, nationality, birth_date, age, alive, sex_id, title_id, band_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				a.ID, a.FirstName, a.LastName, a.Nationality, a.BirthDate, a.Age, a.Alive, a.SexID, a.TitleID, a.BandID)
			if err != nil {
				return err
			}
		}

		albums := []struct {
			ID       int
			Title    string
			Price    float64
			ArtistID *int
			BandID   *int
		}{
			{1, "Blue Train", 56.99, intPtr(1), nil},
			{2, "Jeru", 17.99, intPtr(2), nil},
			{3, "Sarah Vaughan and Clifford Brown", 39.99, intPtr(3), nil},
			{4, "OK Computer", 45.99, nil, intPtr(3)},
			{5, "Abbey Road", 42.99, nil, intPtr(2)},
		}
		for _, a := range albums {
			_, err := DB.Exec(
				"INSERT INTO albums (id, title, price, artist_id, band_id) VALUES (?, ?, ?, ?, ?)",
				a.ID, a.Title, a.Price, a.ArtistID, a.BandID)
			if err != nil {
				return err
			}
		}

		songs := []struct {
			ID       int
			Title    string
			Length   int
			Price    float64
			AlbumID  *int
			ArtistID *int
			BandID   *int
		}{
			{1, "Blue Train", 543, 9.99, intPtr(1), intPtr(1), nil},
			{2, "Lazy Bird", 434, 8.99, intPtr(1), intPtr(1), nil},
			{3, "Jeru", 294, 5.99, intPtr(2), intPtr(2), nil},
			{4, "Paranoid Android", 387, 7.99, intPtr(4), nil, intPtr(3)},
			{5, "Karma Police", 264, 5.99, intPtr(4), nil, intPtr(3)},
			{6, "Come Together", 259, 6.99, intPtr(5), nil, intPtr(2)},
		}
		for _, s := range songs {
			_, err := DB.Exec(
				"INSERT INTO songs (id, title, length, price, album_id, artist_id, band_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
				s.ID, s.Title, s.Length, s.Price, s.AlbumID, s.ArtistID, s.BandID)
			if err != nil {
				return err
			}
		}

		albumSongs := []struct {
			AlbumID int
			SongID  int
		}{
			{1, 1},
			{1, 2},
			{2, 3},
			{4, 4},
			{4, 5},
			{5, 6},
		}
		for _, as := range albumSongs {
			_, err := DB.Exec(
				"INSERT INTO album_songs (album_id, song_id) VALUES (?, ?)",
				as.AlbumID, as.SongID)
			if err != nil {
				return err
			}
		}

		artistSongs := []struct {
			ArtistID int
			SongID   int
		}{
			{1, 1},
			{1, 2},
			{2, 3},
		}
		for _, as := range artistSongs {
			_, err := DB.Exec(
				"INSERT INTO artist_songs (artist_id, song_id) VALUES (?, ?)",
				as.ArtistID, as.SongID)
			if err != nil {
				return err
			}
		}

		bandSongs := []struct {
			BandID int
			SongID int
		}{
			{3, 4},
			{3, 5},
			{2, 6},
		}
		for _, bs := range bandSongs {
			_, err := DB.Exec(
				"INSERT INTO band_songs (band_id, song_id) VALUES (?, ?)",
				bs.BandID, bs.SongID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func intPtr(i int) *int {
	return &i
}
