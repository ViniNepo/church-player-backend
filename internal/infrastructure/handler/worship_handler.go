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

type WorshipHandler struct {
	worshipUseCase usecase.WorshipUseCase
}

func NewWorshipHandler(useCase usecase.WorshipUseCase) *WorshipHandler {
	return &WorshipHandler{
		worshipUseCase: useCase,
	}
}

func (h *WorshipHandler) Create(w http.ResponseWriter, r *http.Request) {
	var worshipDTO dto.CreateWorshipDTO

	err := json.NewDecoder(r.Body).Decode(&worshipDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	worship, err := domain.NewWorship(worshipDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = h.worshipUseCase.Create(worship)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}

func (h *WorshipHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	worship, err := h.worshipUseCase.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.NewEncoder(w).Encode(worship)
	if err != nil {
		return
	}
}

func (h *WorshipHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	worships, err := h.worshipUseCase.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.NewEncoder(w).Encode(worships)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
}

func (h *WorshipHandler) Update(w http.ResponseWriter, r *http.Request) {
	var worship domain.Worship

	// Verifica se o campo do álbum foi enviado
	albumData := r.FormValue("worship")
	if albumData != "" {
		err := json.Unmarshal([]byte(albumData), &worship)
		if err != nil {
			http.Error(w, "Invalid worship data", http.StatusBadRequest)
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
			err = h.worshipUseCase.Update(&worship, fileHeader)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
		}
	} else {
		// Se não houver arquivos, apenas atualiza o worship sem arquivo
		err = h.worshipUseCase.Update(&worship, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}

	json.NewEncoder(w)
}

func (h *WorshipHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := h.worshipUseCase.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	json.NewEncoder(w)
}
