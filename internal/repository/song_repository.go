package repository

import "github.com/nurmeden/music-library/internal/entity"

type SongRepository interface {
	FetchAll(filters map[string]interface{}, limit, offset int) ([]entity.Song, error)
	FetchByID(id int) (*entity.Song, error)
	Store(song *entity.Song) error
	Update(song *entity.Song) error
	Delete(id int) error
}
