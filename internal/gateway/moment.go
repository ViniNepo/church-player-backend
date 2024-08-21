package gateway

import (
	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
)

type Moment interface {
	Create(moment *domain.Moment) error
	FindByID(id string) (*domain.Moment, error)
	FindAllByWorshipID(id string) (*[]domain.Moment, error)
	Update(moment *domain.Moment) error
	UpdateOrder(momentOrderDTO *[]dto.MomentOrderDTO) error
	Delete(id string) error
}
