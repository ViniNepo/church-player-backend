package dto

type (
	CreateSongDTO struct {
		FileName string `json:"file_name"`
		AlbumID  string `json:"album_id"`
	}

	SongWithAlbumDTO struct {
		ID            string `json:"id"`
		Title         string `json:"title"`
		File          string `json:"file"`
		SongNumber    int    `json:"song_number"`
		Order         int    `json:"song_order"`
		AlbumID       string `json:"album_id"`
		AlbumName     string `json:"album_name"`
		AlbumImageURL string `json:"album_image_url"`
	}

	SongUpdateOrderDTO struct {
		ID    string `json:"id"`
		Order int    `json:"song_order"`
	}

	SongDTO struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		File   string `json:"file"`
		Number int    `json:"song_number"`
		Order  int    `json:"song_order"`
	}
)
