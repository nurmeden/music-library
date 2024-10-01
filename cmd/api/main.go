package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nurmeden/music-library/internal/delivery/http"
	postgres "github.com/nurmeden/music-library/internal/infrastructure"
	"github.com/nurmeden/music-library/internal/usecase"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DB_CONNECTION"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	songRepo := postgres.NewPostgresSongRepository(db)
	songUC := usecase.NewSongUseCase(songRepo)

	r := gin.Default()
	http.NewSongHandler(r, songUC)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
