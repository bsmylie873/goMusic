package main

import (
	"fmt"
	"goMusic/controllers"
	"net/http"
	"os"
)

func main() {
	if os.Getenv("JWT_SECRET_KEY") == "" {
		os.Setenv("JWT_SECRET_KEY", "secret")
	}

	Setup("music.db")
	defer CloseDB()

	mux := http.NewServeMux()

	controllers.RegisterAuthRoutes(mux)
	controllers.RegisterAlbumRoutes(mux)
	controllers.RegisterArtistRoutes(mux)
	controllers.RegisterBandRoutes(mux)
	controllers.RegisterSongRoutes(mux)

	fmt.Println("Server starting on :8082")
	http.ListenAndServe("localhost:8082", mux)
}
