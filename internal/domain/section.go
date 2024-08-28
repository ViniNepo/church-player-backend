package domain

import (
	"errors"

	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/google/uuid"
)

type Section struct {
	ID           string `json:"id"`
	SectionName  string `json:"section_name"`
	SectionOrder int    `json:"section_order"`
	WorshipID    string `json:"worship_id"`
}

func NewSection(dto dto.CreateSectionDTO) (*Section, error) {
	if dto.SectionName == "" && dto.WorshipID == "" {
		return nil, errors.New("must provide either a section name and a worship id")
	}

	return &Section{
		ID:          uuid.New().String(),
		SectionName: dto.SectionName,
		WorshipID:   dto.WorshipID,
	}, nil
}
