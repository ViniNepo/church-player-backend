package gateway

import (
	"github.com/ViniNepo/church-player-backend/internal/domain"
)

type Worship interface {
	Create(album *domain.Worship) error
	FindByID(id string) (*domain.Worship, error)
	FindAll() (*[]domain.Worship, error)
	Update(album *domain.Worship) error
	Delete(id string) error
}
