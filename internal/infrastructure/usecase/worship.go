package usecase

import (
	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/gateway"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/sqlite"
)

type WorshipUseCase interface {
	Create(album *domain.Worship) error
	FindByID(id string) (*dto.WorshipDTO, error)
	FindAll() (*[]domain.Worship, error)
	Update(album *domain.Worship) error
	Delete(id string) error
}

type worshipUseCase struct {
	worshipGateway gateway.Worship
	momentGateway  gateway.Moment
	songGateway    gateway.Song
	albumGateway   gateway.Album
	sectionGateway gateway.Section
}

func NewWorshipUseCase(worshipGateway gateway.Worship, momentGateway gateway.Moment, songGateway gateway.Song, albumGateway gateway.Album, sectionGateway *sqlite.SectionDB) WorshipUseCase {
	return &worshipUseCase{
		worshipGateway: worshipGateway,
		momentGateway:  momentGateway,
		songGateway:    songGateway,
		albumGateway:   albumGateway,
		sectionGateway: sectionGateway,
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
					return nil, err
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

func (w *worshipUseCase) Update(album *domain.Worship) error {
	return w.worshipGateway.Update(album)
}

func (w *worshipUseCase) Delete(id string) error {
	return w.worshipGateway.Delete(id)
}
