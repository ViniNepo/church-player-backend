package dto

type (
	CreateAlbumDTO struct {
		Title    string `json:"title"`
		ImageUrl string `json:"image_url"`
	}

	AlbumDTO struct {
		ID        string     `json:"id"`
		Title     string     `json:"title"`
		ImageUrl  string     `json:"image_url"`
		HasNumber bool       `json:"has_numbered_songs"`
		Songs     *[]SongDTO `json:"songs"`
	}
)
