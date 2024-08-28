package domain

import (
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/google/uuid"
)

type Moment struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	SongOrder int    `json:"song_order"`
	SongID    string `json:"song_id"`
	SectionID string `json:"section_id"`
}

func NewMoment(dto dto.CreateMomentDTO) (*Moment, error) {
	return &Moment{
		ID:        uuid.New().String(),
		Title:     dto.Title,
		SectionID: dto.SectionID,
	}, nil
}
