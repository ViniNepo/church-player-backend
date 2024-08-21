package domain

import (
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/google/uuid"
)

type Moment struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	SectionName  string `json:"section_name"`
	SectionOrder int    `json:"section_order"`
	SongOrder    int    `json:"song_order"`
	SongID       string `json:"song_id"`
	WorshipID    string `json:"worship_id"`
}

func NewMoment(dto dto.CreateMomentDTO) (*Moment, error) {
	return &Moment{
		ID:           uuid.New().String(),
		Title:        dto.Title,
		SectionName:  dto.SectionName,
		SectionOrder: dto.SectionOrder,
		SongOrder:    dto.SongOrder,
		SongID:       "",
		WorshipID:    dto.WorshipID,
	}, nil
}
