package gateway

import (
	"github.com/ViniNepo/church-player-backend/internal/domain"
)

type Album interface {
	Create(album *domain.Album) error
	FindByID(id string) (*domain.Album, error)
	FindAll() (*[]domain.Album, error)
	Update(album *domain.Album) error
	Delete(id string) error
}
