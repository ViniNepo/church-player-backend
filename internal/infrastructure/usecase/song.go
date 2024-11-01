package usecase

import (
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/filemanager"
	"github.com/ViniNepo/church-player-backend/internal/gateway"
)

type SongUseCase interface {
	Create(song *domain.Song, fileHeader *multipart.FileHeader) error
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

func (s *songUseCase) Create(song *domain.Song, fileHeader *multipart.FileHeader) error {
	fileName := ""
	for _, st := range song.File {
		if (st > 96 && st < 123) || (st > 64 && st < 91) || (st == 8226 || st == 45 || st == 46 || st == 32) || (st > 47 && st < 58) {
			fileName = fileName + string(st)
		}
	}
	song.File = fileName
	fileHeader.Filename = fileName
	destinationPath := filepath.Join("dist/church-player/assets/files", song.File)

	if s.fileManager.Exists(destinationPath) {
		return fmt.Errorf("file %s already exists", song.File)
	}

	err := s.songGateway.Create(song, nil)
	if err != nil {
		return err
	}

	err = s.fileManager.Copy(fileHeader, "dist/church-player/assets/files")
	if err != nil {
		return s.songGateway.Delete(song.ID)
	}

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
		destinationPath := filepath.Join("dist/church-player/assets/files", song.File)
		_ = s.fileManager.Delete(destinationPath)
	}()

	return s.songGateway.Delete(id)
}
