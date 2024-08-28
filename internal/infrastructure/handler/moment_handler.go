package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ViniNepo/church-player-backend/internal/domain"
	"github.com/ViniNepo/church-player-backend/internal/domain/dto"
	"github.com/ViniNepo/church-player-backend/internal/infrastructure/usecase"
	"github.com/gorilla/mux"
)

type MomentHandler struct {
	momentUseCase usecase.MomentUseCase
}

func NewMomentHandler(useCase usecase.MomentUseCase) *MomentHandler {
	return &MomentHandler{
		momentUseCase: useCase,
	}
}

func (m *MomentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var momentDTO dto.CreateMomentDTO

	err := json.NewDecoder(r.Body).Decode(&momentDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	moment, err := domain.NewMoment(momentDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.momentUseCase.CreateMoment(moment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	json.NewEncoder(w)
}

func (m *MomentHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	moment, err := m.momentUseCase.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.NewEncoder(w).Encode(moment)
	if err != nil {
		return
	}
}

func (m *MomentHandler) FindAllBySectionID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	moment, err := m.momentUseCase.FindAllBySectionID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	err = json.NewEncoder(w).Encode(moment)
	if err != nil {
		return
	}
}

func (m *MomentHandler) Update(w http.ResponseWriter, r *http.Request) {
	var moment domain.Moment

	err := json.NewDecoder(r.Body).Decode(&moment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.momentUseCase.Update(&moment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}

func (m *MomentHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var momentOrder []dto.MomentOrderDTO

	err := json.NewDecoder(r.Body).Decode(&momentOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.momentUseCase.UpdateOrder(&momentOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}

func (m *MomentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := m.momentUseCase.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}
