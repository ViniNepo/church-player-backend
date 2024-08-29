package gateway

import (
	"mime/multipart"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
)

type Song interface {
	Create(song *domain.Song, fileHeader *multipart.FileHeader) error
	FindAll() (*[]dto.SongWithAlbumDTO, error)
	FindAllByAlbumID(id string) (*[]domain.Song, error)
	FindByID(id string) (*domain.Song, error)
	Update(song *domain.Song) error
	UpdateOrder(songUpdateOrderDTO *[]dto.SongUpdateOrderDTO) error
	Delete(id string) error
}
