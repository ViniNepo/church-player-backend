package usecase

import (
	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/gateway"
)

type SectionUseCase interface {
	Create(section *domain.Section) error
	FindAllByWorshipID(id string) (*[]domain.Section, error)
	Update(section *domain.Section) error
	Delete(id string) error
}

type sectionUseCase struct {
	sectionGateway gateway.Section
}

func NewSectionUseCase(sectionGateway gateway.Section) SectionUseCase {
	return &sectionUseCase{sectionGateway: sectionGateway}
}

func (s *sectionUseCase) Create(section *domain.Section) error {
	return s.sectionGateway.Create(section)
}

func (s *sectionUseCase) FindAllByWorshipID(id string) (*[]domain.Section, error) {
	return s.sectionGateway.FindAllByWorshipID(id)
}

func (s *sectionUseCase) Update(section *domain.Section) error {
	return s.sectionGateway.Update(section)
}

func (s *sectionUseCase) Delete(id string) error {
	return s.sectionGateway.Delete(id)
}
