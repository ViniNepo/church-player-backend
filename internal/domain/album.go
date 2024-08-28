package domain

import (
	"errors"

	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/google/uuid"
)

type Album struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	ImageUrl   string `json:"image_url"`
	AlbumOrder int    `json:"album_order"`
}

func NewAlbum(dto dto.CreateAlbumDTO) (*Album, error) {
	if dto.Title == "" {
		return nil, errors.New("album must have a title")
	}

	return &Album{
		ID:       uuid.New().String(),
		Title:    dto.Title,
		ImageUrl: dto.ImageUrl,
	}, nil
}
