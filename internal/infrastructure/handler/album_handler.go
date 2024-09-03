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

type AlbumHandler struct {
	albumUseCase usecase.AlbumUseCase
}

func NewAlbumHandler(albumUseCase usecase.AlbumUseCase) *AlbumHandler {
	return &AlbumHandler{
		albumUseCase: albumUseCase,
	}
}

func (h *AlbumHandler) Create(w http.ResponseWriter, r *http.Request) {
	var albumDTO dto.CreateAlbumDTO

	err := json.NewDecoder(r.Body).Decode(&albumDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	album, err := domain.NewAlbum(albumDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.albumUseCase.CreateAlbum(album)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}

func (h *AlbumHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	album, err := h.albumUseCase.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.NewEncoder(w).Encode(album)
	if err != nil {
		return
	}
}

func (h *AlbumHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	albums, err := h.albumUseCase.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.NewEncoder(w).Encode(albums)
}

func (h *AlbumHandler) Update(w http.ResponseWriter, r *http.Request) {
	var album domain.Album

	// Verifica se o campo do álbum foi enviado
	albumData := r.FormValue("album")
	if albumData != "" {
		err := json.Unmarshal([]byte(albumData), &album)
		if err != nil {
			http.Error(w, "Invalid album data", http.StatusBadRequest)
			return
		}
	}

	// Tenta analisar o formulário multipart (pode conter arquivos)
	err := r.ParseMultipartForm(50 << 20) // 100 MB
	if err != nil && err != http.ErrNotMultipart {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	destDir := "././ChurchPlayer/resources/app/dist/church-player/assets/images"
	os.MkdirAll(destDir, os.ModePerm)

	// Processa os arquivos, se houver
	files := r.MultipartForm.File["files"]
	if len(files) > 0 {
		for _, fileHeader := range files {
			err = h.albumUseCase.Update(&album, fileHeader)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
		}
	} else {
		// Se não houver arquivos, apenas atualiza o álbum sem arquivo
		err = h.albumUseCase.Update(&album, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(album)
}

func (h *AlbumHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := h.albumUseCase.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}
