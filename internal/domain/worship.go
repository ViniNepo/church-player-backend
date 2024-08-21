package domain

import (
	"errors"

	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/google/uuid"
)

type Worship struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ImageUrl string `json:"image_url"`
}

func NewWorship(dto dto.CreateWorshipDTO) (*Worship, error) {
	if dto.Title == "" {
		return nil, errors.New("worship must have a title")
	}

	return &Worship{
		ID:       uuid.New().String(),
		Title:    dto.Title,
		ImageUrl: dto.ImageUrl,
	}, nil
}
