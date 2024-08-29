package usecase

import (
	"github.com/ViniNepo/church-player-backend/internal/filemanager"
)

type FileControlUseCase interface {
	Play(path string) error
}

type fileControlUseCase struct {
	fileManager filemanager.FileManager
}

func NewFileControlUseCase(manager filemanager.FileManager) FileControlUseCase {
	return &fileControlUseCase{fileManager: manager}
}

func (f *fileControlUseCase) Play(path string) error {
	return f.fileManager.Read(path)
}
