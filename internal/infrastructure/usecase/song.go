package usecase

import (
	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
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
}

func NewSongUseCase(songGateway gateway.Song) SongUseCase {
	return &songUseCase{songGateway: songGateway}
}

func (s *songUseCase) Create(song *domain.Song) error {
	return s.songGateway.Create(song)
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
	return s.songGateway.Delete(id)
}
