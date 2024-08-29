package handler

import (
	"encoding/json"
	"net/http"
	"path/filepath"

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
	filePath := mux.Vars(r)["path"]

	if filePath == "" {
		http.Error(w, "file path is required", http.StatusBadRequest)
		return
	}

	destPath := filepath.Join("./files", filePath)
	err := h.fileControlUseCase.Play(destPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	json.NewEncoder(w)
}
