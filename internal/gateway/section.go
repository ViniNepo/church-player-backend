package gateway

import (
	"github.com/ViniNepo/church-player-backend/internal/domain"
)

type Section interface {
	Create(section *domain.Section) error
	FindAllByWorshipID(id string) (*[]domain.Section, error)
	Update(section *domain.Section) error
	Delete(id string) error
}
