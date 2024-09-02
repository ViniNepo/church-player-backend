package usecase

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/filemanager"
	"github.com/ViniNepo/church-player-backend/internal/gateway"
)

type AlbumUseCase interface {
	CreateAlbum(album *domain.Album) error
	FindByID(id string) (*dto.AlbumDTO, error)
	FindAll() (*[]domain.Album, error)
	Update(album *domain.Album, fileHeader *multipart.FileHeader) error
	Delete(id string) error
}

type albumUseCase struct {
	albumGateway gateway.Album
	songGateway  gateway.Song
	fileManager  filemanager.FileManager
}

func NewAlbumUseCase(albumGateway gateway.Album, songGateway gateway.Song, manager filemanager.FileManager) AlbumUseCase {
	return &albumUseCase{albumGateway: albumGateway, songGateway: songGateway, fileManager: manager}
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

func (u *albumUseCase) Update(album *domain.Album, fileHeader *multipart.FileHeader) error {
	if fileHeader != nil {
		destinationPath := filepath.Join("./dist/church-player/assets/images", album.ImageUrl)
		if album.ImageUrl != "church-player-2.png" && u.fileManager.Exists(destinationPath) {
			return errors.New("image exists")
		}

		err := u.fileManager.Copy(fileHeader, "./dist/church-player/assets/images")
		if err != nil {
			return err
		}
	}

	time.Sleep(2 * time.Second)

	return u.albumGateway.Update(album)
}

func (u *albumUseCase) Delete(id string) error {
	songs, err := u.songGateway.FindAllByAlbumID(id)
	if err != nil {
		return err
	}

	album, err := u.albumGateway.FindByID(id)
	if err != nil {
		return err
	}

	go func() {
		if album.ImageUrl != "church-player-2.png" {
			destinationPath := filepath.Join("./dist/church-player/assets/images", album.ImageUrl)
			_ = u.fileManager.Delete(destinationPath)
		}

		for _, song := range *songs {
			destinationPath := filepath.Join("./dist/church-player/assets/files", song.File)
			_ = u.fileManager.Delete(destinationPath)
		}
	}()

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
