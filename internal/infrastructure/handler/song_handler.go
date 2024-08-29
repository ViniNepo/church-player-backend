package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/usecase"
	"github.com/gorilla/mux"
)

type SongHandler struct {
	songUseCase usecase.SongUseCase
}

func NewSongHandler(useCase usecase.SongUseCase) *SongHandler {
	return &SongHandler{
		songUseCase: useCase,
	}
}

func (h *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["album_id"]

	if id == "" {
		http.Error(w, "album_id is required", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(100 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	destDir := "./files"
	os.MkdirAll(destDir, os.ModePerm)

	files := r.MultipartForm.File["files"]
	for _, fileHeader := range files {
		song, err := domain.NewSong(dto.CreateSongDTO{
			FileName: fileHeader.Filename,
			AlbumID:  id,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		err = h.songUseCase.Create(song, fileHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}
	json.NewEncoder(w)
}

func (h *SongHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	songs, err := h.songUseCase.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.NewEncoder(w).Encode(songs)
	if err != nil {
		return
	}
}

func (h *SongHandler) FindById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	song, err := h.songUseCase.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.NewEncoder(w).Encode(song)
	if err != nil {
		return
	}
}

func (h *SongHandler) Update(w http.ResponseWriter, r *http.Request) {
	var song domain.Song

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.songUseCase.Update(&song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}

func (h *SongHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var songs []dto.SongUpdateOrderDTO

	err := json.NewDecoder(r.Body).Decode(&songs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.songUseCase.UpdateOrder(&songs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}

func (h *SongHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := h.songUseCase.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}
