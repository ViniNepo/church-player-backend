package usecase

import (
	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/gateway"
)

type AlbumUseCase interface {
	CreateAlbum(album *domain.Album) error
	FindByID(id string) (*dto.AlbumDTO, error)
	FindAll() (*[]domain.Album, error)
	Update(album *domain.Album) error
	Delete(id string) error
}

type albumUseCase struct {
	albumGateway gateway.Album
	songGateway  gateway.Song
}

func NewAlbumUseCase(albumGateway gateway.Album, songGateway gateway.Song) AlbumUseCase {
	return &albumUseCase{albumGateway: albumGateway, songGateway: songGateway}
}

func (u *albumUseCase) CreateAlbum(album *domain.Album) error {
	return u.albumGateway.Create(album)
}

func (u *albumUseCase) FindByID(id string) (*dto.AlbumDTO, error) {
	album, err := u.albumGateway.FindByID(id)
	if err != nil {
		return nil, err
	}

	songs, err := u.songGateway.FindAllByAlbumID(album.ID)
	if err != nil {
		return nil, err
	}

	albumDTO := &dto.AlbumDTO{
		ID:       album.ID,
		Title:    album.Title,
		ImageUrl: album.ImageUrl,
		Songs:    toSongDTO(songs),
	}

	if albumDTO.Songs != nil {
		for _, songDTO := range *albumDTO.Songs {
			if songDTO.Number != 0 {
				albumDTO.HasNumber = true
			}
		}

	}

	return albumDTO, nil
}

func (u *albumUseCase) FindAll() (*[]domain.Album, error) {
	return u.albumGateway.FindAll()
}

func (u *albumUseCase) Update(album *domain.Album) error {
	return u.albumGateway.Update(album)
}

func (u *albumUseCase) Delete(id string) error {
	return u.albumGateway.Delete(id)
}

func toSongDTO(songs *[]domain.Song) *[]dto.SongDTO {
	var ss []dto.SongDTO

	for _, song := range *songs {
		songDTO := dto.SongDTO{
			ID:     song.ID,
			Title:  song.Title,
			Number: song.Number,
			File:   song.File,
			Order:  song.Order,
		}

		ss = append(ss, songDTO)
	}

	return &ss
}
