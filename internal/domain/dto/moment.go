package dto

type (
	CreateMomentDTO struct {
		Title        string `json:"title"`
		SectionName  string `json:"section_name"`
		SectionOrder int    `json:"section_order"`
		SongOrder    int    `json:"order"`
		WorshipID    string `json:"worship_id"`
	}

	MomentDTO struct {
		ID           string            `json:"id"`
		Title        string            `json:"title"`
		SectionName  string            `json:"section_name"`
		SectionOrder int               `json:"section_order"`
		SongOrder    int               `json:"song_order"`
		Song         *SongWithAlbumDTO `json:"song"`
	}

	SectionDTO struct {
		Name    string       `json:"section_name"`
		Moments *[]MomentDTO `json:"moments"`
	}

	MomentOrderDTO struct {
		ID        string `json:"id"`
		SongOrder int    `json:"song_order"`
	}
)
