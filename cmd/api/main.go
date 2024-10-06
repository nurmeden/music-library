package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nurmeden/music-library/internal/delivery/http"
	postgres "github.com/nurmeden/music-library/internal/infrastructure"
	"github.com/nurmeden/music-library/internal/logger"
	"github.com/nurmeden/music-library/internal/usecase"
	"github.com/nurmeden/music-library/utils"
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

	err = utils.RunMigrations(db)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	appLogger := logger.NewZapLogger()

	songRepo := postgres.NewPostgresSongRepository(db, appLogger)
	songUC := usecase.NewSongUseCase(songRepo, appLogger)

	r := gin.Default()
	http.NewSongHandler(r, songUC, appLogger)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
