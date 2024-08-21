package usecase

import (
	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/gateway"
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
}

func NewWorshipUseCase(worshipGateway gateway.Worship, momentGateway gateway.Moment, songGateway gateway.Song, albumGateway gateway.Album) WorshipUseCase {
	return &worshipUseCase{
		worshipGateway: worshipGateway,
		momentGateway:  momentGateway,
		songGateway:    songGateway,
		albumGateway:   albumGateway,
	}
}

func (w *worshipUseCase) Create(worship *domain.Worship) error {
	return w.worshipGateway.Create(worship)
}

func (w *worshipUseCase) FindByID(id string) (*dto.WorshipDTO, error) {
	var worshipDTO dto.WorshipDTO
	var sections []dto.SectionDTO

	worship, err := w.worshipGateway.FindByID(id)
	if err != nil {
		return nil, err
	}

	moments, err := w.momentGateway.FindAllByWorshipID(worship.ID)
	if err != nil {
		return nil, err
	}

	worshipDTO.ID = worship.ID
	worshipDTO.Title = worship.Title
	worshipDTO.ImageUrl = worship.ImageUrl
	worshipDTO.Sections = &sections

	sectionMap := make(map[string][]dto.MomentDTO)
	for _, moment := range *moments {
		momentDTO := dto.MomentDTO{
			ID:           moment.ID,
			Title:        moment.Title,
			SectionName:  moment.SectionName,
			SectionOrder: moment.SectionOrder,
			SongOrder:    moment.SongOrder,
		}

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
				Number:        song.Number,
				File:          song.File,
				Order:         song.Order,
				AlbumID:       song.AlbumID,
				AlbumName:     album.Title,
				AlbumImageURL: album.ImageUrl,
			}

			momentDTO.Song = &songDTO
		}

		sectionMap[momentDTO.SectionName] = append(sectionMap[momentDTO.SectionName], momentDTO)
	}

	for sectionName, ms := range sectionMap {
		section := dto.SectionDTO{
			Name:    sectionName,
			Moments: &ms,
		}
		sections = append(sections, section)
	}

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
