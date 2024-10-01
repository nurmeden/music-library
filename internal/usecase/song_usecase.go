package usecase

import (
	"github.com/nurmeden/music-library/internal/entity"
	"github.com/nurmeden/music-library/internal/repository"
)

type SongUseCase struct {
	SongRepo repository.SongRepository
}

func NewSongUseCase(repo repository.SongRepository) *SongUseCase {
	return &SongUseCase{SongRepo: repo}
}

func (uc *SongUseCase) FetchAll(filters map[string]interface{}, limit, offset int) ([]entity.Song, error) {
	return uc.SongRepo.FetchAll(filters, limit, offset)
}

func (uc *SongUseCase) FetchByID(id int) (*entity.Song, error) {
	return uc.SongRepo.FetchByID(id)
}

func (uc *SongUseCase) AddNewSong(song *entity.Song) error {
	return uc.SongRepo.Store(song)
}

func (uc *SongUseCase) UpdateSong(song *entity.Song) error {
	return uc.SongRepo.Update(song)
}

func (uc *SongUseCase) DeleteSong(id int) error {
	return uc.SongRepo.Delete(id)
}
