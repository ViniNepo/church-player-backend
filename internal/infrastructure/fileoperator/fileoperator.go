package fileoperator

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/ViniNepo/church-player-backend/internal/filemanager"
)

type fileOperator struct{}

func (fo fileOperator) Copy(fileHeader *multipart.FileHeader, destinationPath string) error {
	//Abre o arquivo de origem
	file, err := fileHeader.Open()
	if err != nil {
		return errors.New("Could not open uploaded file")
	}
	defer file.Close()

	//Cria o arquivo de destino
	destPath := filepath.Join(destinationPath, fileHeader.Filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		return errors.New("Could not create destination file")
	}
	defer destFile.Close()
	//
	//Copia o arquivo para o destino
	_, err = io.Copy(destFile, file)
	if err != nil {
		return errors.New("Could not copy file to destination")
	}

	return nil
}

func (fo fileOperator) Read(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", path)
	case "darwin": // macOS
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

	return nil
}

func (fo fileOperator) Delete(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("could not delete file: %w", err)
	}
	return nil
}

func (fo fileOperator) Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func NewFileManager() filemanager.FileManager {
	return fileOperator{}
}
