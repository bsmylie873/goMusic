package main

import (
	"fmt"
	"goMusic/controllers"
	"goMusic/db"
	"log"
	"net/http"
	"os"
)

func main() {
	if os.Getenv("JWT_SECRET_KEY") == "" {
		os.Setenv("JWT_SECRET_KEY", "secret")
	}
	err := db.InitDB("music.db")
	if err != nil {
		log.Fatal("Database initialization failed: ", err)
	}

	err = db.SeedDB()
	if err != nil {
		log.Fatal("Database seeding failed: ", err)
	}

	mux := http.NewServeMux()

	controllers.RegisterAuthRoutes(mux)
	controllers.RegisterAlbumRoutes(mux)
	controllers.RegisterArtistRoutes(mux)
	controllers.RegisterBandRoutes(mux)

	fmt.Println("Server starting on :8082")
	http.ListenAndServe("localhost:8082", mux)
}
