package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ViniNepo/church-player-backend/internal/infrastructure/usecase"
	"github.com/gorilla/mux"
)

type FileControlHandler struct {
	fileControlUseCase usecase.FileControlUseCase
}

func NewFileControlHandler(fileControlUseCase usecase.FileControlUseCase) *FileControlHandler {
	return &FileControlHandler{
		fileControlUseCase: fileControlUseCase,
	}
}

func (h *FileControlHandler) Play(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["path"]

	if id == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}

	err := h.fileControlUseCase.Play(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}
