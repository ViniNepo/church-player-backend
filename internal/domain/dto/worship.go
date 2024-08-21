package dto

type (
	CreateWorshipDTO struct {
		Title    string `json:"title"`
		ImageUrl string `json:"image_url"`
	}

	WorshipDTO struct {
		ID       string        `json:"id"`
		Title    string        `json:"title"`
		ImageUrl string        `json:"image_url"`
		Sections *[]SectionDTO `json:"sections"`
	}
)
