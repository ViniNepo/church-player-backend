package dto

type (
	CreateMomentDTO struct {
		Title     string `json:"title"`
		SectionID string `json:"section_id"`
	}

	MomentDTO struct {
		ID        string            `json:"id"`
		Title     string            `json:"title"`
		SectionID string            `json:"section_id"`
		Song      *SongWithAlbumDTO `json:"song"`
	}

	MomentOrderDTO struct {
		ID        string `json:"id"`
		SongOrder int    `json:"song_order"`
		SectionID string `json:"section_id"`
	}
)
