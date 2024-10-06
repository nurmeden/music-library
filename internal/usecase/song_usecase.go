package usecase

import (
	"github.com/nurmeden/music-library/internal/entity"
	"github.com/nurmeden/music-library/internal/logger"
	"github.com/nurmeden/music-library/internal/repository"
)

type SongUseCase struct {
	repo   repository.SongRepository
	logger logger.Logger
}

func NewSongUseCase(repo repository.SongRepository, logger logger.Logger) *SongUseCase {
	return &SongUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *SongUseCase) FetchAll(filters map[string]interface{}, limit, offset int) ([]entity.Song, error) {
	return uc.repo.FetchAll(filters, limit, offset)
}

func (uc *SongUseCase) FetchByID(id int) (*entity.Song, error) {
	return uc.repo.FetchByID(id)
}

func (uc *SongUseCase) AddNewSong(song *entity.Song) error {
	return uc.repo.Store(song)
}

func (uc *SongUseCase) UpdateSong(song *entity.Song) error {
	return uc.repo.Update(song)
}

func (uc *SongUseCase) DeleteSong(id int) error {
	return uc.repo.Delete(id)
}
