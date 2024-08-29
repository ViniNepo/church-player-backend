package usecase

import (
	"errors"
	"path/filepath"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/filemanager"
	"github.com/ViniNepo/church-player-backend/internal/gateway"
)

type SongUseCase interface {
	Create(song *domain.Song) error
	FindAll() (*[]dto.SongWithAlbumDTO, error)
	FindByID(id string) (*domain.Song, error)
	FindAllByAlbumID(id string) (*[]domain.Song, error)
	Update(song *domain.Song) error
	UpdateOrder(songUpdateOrderDTO *[]dto.SongUpdateOrderDTO) error
	Delete(id string) error
}

type songUseCase struct {
	songGateway gateway.Song
	fileManager filemanager.FileManager
}

func NewSongUseCase(songGateway gateway.Song, manager filemanager.FileManager) SongUseCase {
	return &songUseCase{songGateway: songGateway, fileManager: manager}
}

func (s *songUseCase) Create(song *domain.Song) error {
	if s.fileManager.Exists(song.File) {
		return errors.New("file exists")
	}

	err := s.songGateway.Create(song)
	if err != nil {
		return err
	}

	go func() {
		destinationPath := filepath.Join("", song.File)
		err = s.fileManager.Copy(song.File, destinationPath)
		if err != nil {
			_ = s.songGateway.Delete(song.ID)
		}

	}()

	return nil
}

func (s *songUseCase) FindAll() (*[]dto.SongWithAlbumDTO, error) {
	return s.songGateway.FindAll()
}

func (s *songUseCase) FindByID(id string) (*domain.Song, error) {
	return s.songGateway.FindByID(id)
}

func (s *songUseCase) FindAllByAlbumID(id string) (*[]domain.Song, error) {
	return s.songGateway.FindAllByAlbumID(id)
}

func (s *songUseCase) Update(song *domain.Song) error {
	return s.songGateway.Update(song)
}

func (s *songUseCase) UpdateOrder(songUpdateOrderDTO *[]dto.SongUpdateOrderDTO) error {
	return s.songGateway.UpdateOrder(songUpdateOrderDTO)
}

func (s *songUseCase) Delete(id string) error {
	song, err := s.FindByID(id)
	if err != nil {
		return err
	}

	go func() {
		_ = s.fileManager.Delete(song.File)
	}()

	return s.songGateway.Delete(id)
}
