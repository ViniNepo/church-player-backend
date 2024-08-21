package dto

type (
	CreateSongDTO struct {
		FileName string `json:"file_name"`
		Order    int    `json:"order_song"`
		AlbumID  string `json:"album_id"`
	}

	SongWithAlbumDTO struct {
		ID            string `json:"id"`
		Title         string `json:"title"`
		Number        int    `json:"number"`
		File          string `json:"file"`
		Order         int    `json:"order_song"`
		AlbumID       string `json:"album_id"`
		AlbumName     string `json:"album_name"`
		AlbumImageURL string `json:"album_image_url"`
	}

	SongUpdateOrderDTO struct {
		ID    string `json:"id"`
		Order int    `json:"order_song"`
	}

	SongDTO struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Number int    `json:"number"`
		File   string `json:"file"`
		Order  int    `json:"order_song"`
	}
)
