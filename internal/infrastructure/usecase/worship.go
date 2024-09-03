package usecase

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/filemanager"
	"github.com/ViniNepo/church-player-backend/internal/gateway"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/sqlite"
)

type WorshipUseCase interface {
	Create(album *domain.Worship) error
	FindByID(id string) (*dto.WorshipDTO, error)
	FindAll() (*[]domain.Worship, error)
	Update(album *domain.Worship, fileHeader *multipart.FileHeader) error
	Delete(id string) error
}

type worshipUseCase struct {
	worshipGateway gateway.Worship
	momentGateway  gateway.Moment
	songGateway    gateway.Song
	albumGateway   gateway.Album
	sectionGateway gateway.Section
	fileManager    filemanager.FileManager
}

func NewWorshipUseCase(worshipGateway gateway.Worship, momentGateway gateway.Moment, songGateway gateway.Song, albumGateway gateway.Album, sectionGateway *sqlite.SectionDB, manager filemanager.FileManager) WorshipUseCase {
	return &worshipUseCase{
		worshipGateway: worshipGateway,
		momentGateway:  momentGateway,
		songGateway:    songGateway,
		albumGateway:   albumGateway,
		sectionGateway: sectionGateway,
		fileManager:    manager,
	}
}

func (w *worshipUseCase) Create(worship *domain.Worship) error {
	return w.worshipGateway.Create(worship)
}

func (w *worshipUseCase) FindByID(id string) (*dto.WorshipDTO, error) {
	var worshipDTO dto.WorshipDTO
	var sectionsDTO []dto.SectionDTO

	worship, err := w.worshipGateway.FindByID(id)
	if err != nil {
		return nil, err
	}

	sections, err := w.sectionGateway.FindAllByWorshipID(worship.ID)
	if err != nil {
		return nil, err
	}

	for _, section := range *sections {
		var sectionDTO dto.SectionDTO
		var momentsDTO []dto.MomentDTO

		moments, err := w.momentGateway.FindAllBySectionID(section.ID)
		if err != nil {
			return nil, err
		}

		for _, moment := range *moments {
			momentDTO := dto.MomentDTO{}

			if moment.SongID != "" {
				song, err := w.songGateway.FindByID(moment.SongID)
				if err != nil {
					momentDTO.ID = moment.ID
					momentDTO.Title = moment.Title
					momentDTO.SectionID = moment.SectionID

					momentsDTO = append(momentsDTO, momentDTO)
					continue
				}

				album, err := w.albumGateway.FindByID(song.AlbumID)
				if err != nil {
					return nil, err
				}

				songDTO := dto.SongWithAlbumDTO{
					ID:            song.ID,
					Title:         song.Title,
					File:          song.File,
					Order:         song.Order,
					AlbumID:       song.AlbumID,
					AlbumName:     album.Title,
					AlbumImageURL: album.ImageUrl,
				}

				momentDTO.Song = &songDTO
			}

			momentDTO.ID = moment.ID
			momentDTO.Title = moment.Title
			momentDTO.SectionID = moment.SectionID

			momentsDTO = append(momentsDTO, momentDTO)
		}

		sectionDTO.ID = section.ID
		sectionDTO.Name = section.SectionName
		sectionDTO.Moments = &momentsDTO

		sectionsDTO = append(sectionsDTO, sectionDTO)
	}

	worshipDTO.ID = worship.ID
	worshipDTO.Title = worship.Title
	worshipDTO.ImageUrl = worship.ImageUrl
	worshipDTO.Sections = &sectionsDTO

	return &worshipDTO, nil
}

func (w *worshipUseCase) FindAll() (*[]domain.Worship, error) {
	return w.worshipGateway.FindAll()
}

func (w *worshipUseCase) Update(worship *domain.Worship, fileHeader *multipart.FileHeader) error {
	if fileHeader != nil {
		if worship.ImageUrl == "church-player-2.png" {
			return fmt.Errorf("image cannot be use this name: %s", worship.ImageUrl)
		}

		err := w.fileManager.Copy(fileHeader, "dist/church-player/assets/images")
		if err != nil {
			return err
		}
	}

	time.Sleep(2 * time.Second)

	return w.worshipGateway.Update(worship)
}

func (w *worshipUseCase) Delete(id string) error {
	worship, err := w.worshipGateway.FindByID(id)
	if err != nil {
		return err
	}

	go func() {
		if worship.ImageUrl != "church-player-2.png" {
			destinationPath := filepath.Join("dist/church-player/assets/images", worship.ImageUrl)
			_ = w.fileManager.Delete(destinationPath)
		}
	}()

	sections, err := w.sectionGateway.FindAllByWorshipID(id)
	if err != nil {
		return err
	}

	for _, section := range *sections {
		err = w.sectionGateway.DeleteByWorshipID(id, section.ID)
		if err != nil {
			return err
		}
	}

	return w.worshipGateway.Delete(id)
}
