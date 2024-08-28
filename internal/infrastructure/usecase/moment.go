package usecase

import (
	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/gateway"
)

type MomentUseCase interface {
	CreateMoment(moment *domain.Moment) error
	FindByID(id string) (*domain.Moment, error)
	FindAllBySectionID(id string) (*[]domain.Moment, error)
	Update(moment *domain.Moment) error
	UpdateOrder(momentOrderDTO *[]dto.MomentOrderDTO) error
	Delete(id string) error
}

type momentUseCase struct {
	momentGateway gateway.Moment
}

func NewMomentUseCase(momentGateway gateway.Moment) MomentUseCase {
	return &momentUseCase{momentGateway: momentGateway}
}

func (m *momentUseCase) CreateMoment(moment *domain.Moment) error {
	return m.momentGateway.Create(moment)
}

func (m *momentUseCase) FindByID(id string) (*domain.Moment, error) {
	return m.momentGateway.FindByID(id)
}

func (m *momentUseCase) FindAllBySectionID(id string) (*[]domain.Moment, error) {
	return m.momentGateway.FindAllBySectionID(id)
}

func (m *momentUseCase) Update(moment *domain.Moment) error {
	return m.momentGateway.Update(moment)
}

func (m *momentUseCase) UpdateOrder(momentOrderDTO *[]dto.MomentOrderDTO) error {
	return m.momentGateway.UpdateOrder(momentOrderDTO)
}

func (m *momentUseCase) Delete(id string) error {
	return m.momentGateway.Delete(id)
}
