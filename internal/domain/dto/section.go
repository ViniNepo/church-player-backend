package dto

type (
	CreateSectionDTO struct {
		SectionName string `json:"section_name"`
		WorshipID   string `json:"worship_id"`
	}

	SectionDTO struct {
		ID      string       `json:"id"`
		Name    string       `json:"section_name"`
		Moments *[]MomentDTO `json:"moments"`
	}
)
