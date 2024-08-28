package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/google/uuid"
)

type Song struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Number  int    `json:"song_number"`
	File    string `json:"file"`
	Order   int    `json:"song_order"`
	AlbumID string `json:"album_id"`
}

func NewSong(dto dto.CreateSongDTO) (*Song, error) {
	if dto.FileName == "" || dto.AlbumID == "" {
		return nil, errors.New("file name, file path and order cannot be empty")
	}

	song := Song{}
	if IsValidSongFormat(dto.FileName) {
		number, title, err := GetSongDetails(dto.FileName)
		if err != nil {
			return nil, err
		}

		song.Title = title
		song.Number = number
	} else {
		song.Title = dto.FileName
	}

	return &Song{
		ID:      uuid.New().String(),
		Title:   song.Title,
		Number:  song.Number,
		File:    dto.FileName,
		AlbumID: dto.AlbumID,
	}, nil
}

func IsValidSongFormat(text string) bool {
	regex := regexp.MustCompile(`^\s*\d+\s*-\s*.+\s*$`)
	return regex.MatchString(text)
}

func GetSongDetails(text string) (int, string, error) {
	regex := regexp.MustCompile(`^\s*(\d+)\s*-\s*(.+)\s*$`)
	matches := regex.FindStringSubmatch(text)

	if len(matches) != 3 {
		return 0, "", fmt.Errorf("invalid format")
	}

	number, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", err
	}

	name := matches[2]
	return number, name, nil
}
