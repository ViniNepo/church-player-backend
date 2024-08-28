package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/usecase"
	"github.com/gorilla/mux"
)

type SectionHandler struct {
	sectionUseCase usecase.SectionUseCase
}

func NewSectionHandler(useCase usecase.SectionUseCase) *SectionHandler {
	return &SectionHandler{
		sectionUseCase: useCase,
	}
}

func (h *SectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sectionDTO dto.CreateSectionDTO

	err := json.NewDecoder(r.Body).Decode(&sectionDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	section, err := domain.NewSection(sectionDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.sectionUseCase.Create(section)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}

func (h *SectionHandler) FindAllByWorshipID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	songs, err := h.sectionUseCase.FindAllByWorshipID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.NewEncoder(w).Encode(songs)
	if err != nil {
		return
	}
}

func (h *SectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	var section domain.Section

	err := json.NewDecoder(r.Body).Decode(&section)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.sectionUseCase.Update(&section)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}

func (h *SectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := h.sectionUseCase.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}
