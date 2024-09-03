package fileoperator

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/ViniNepo/church-player-backend/internal/filemanager"
)

type fileOperator struct {
	osPath string
}

func (fo fileOperator) Copy(fileHeader *multipart.FileHeader, destinationPath string) error {
	//Abre o arquivo de origem
	file, err := fileHeader.Open()
	if err != nil {
		return errors.New("Could not open uploaded file")
	}
	defer file.Close()

	//Cria o arquivo de destino
	destPath := filepath.Join(fo.osPath, destinationPath, fileHeader.Filename)
	log.Print(destPath)
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
		log.Print(filepath.Join(fo.osPath, path))
		cmd = exec.Command("cmd", "/C", "start", filepath.Join(fo.osPath, path))
	case "darwin": // macOS
		log.Print(filepath.Join(fo.osPath, path))
		cmd = exec.Command("open", filepath.Join(fo.osPath, path))
	case "linux":
		log.Print(filepath.Join(fo.osPath, path))
		cmd = exec.Command("xdg-open", filepath.Join(fo.osPath, path))
	default:
		return fmt.Errorf("unsupported platform")
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

	return nil
}

func (fo fileOperator) Delete(path string) error {
	filepath.Join(fo.osPath, path)
	if err := os.Remove(filepath.Join(fo.osPath, path)); err != nil {
		return fmt.Errorf("could not delete file: %w", err)
	}
	return nil
}

func (fo fileOperator) Exists(path string) bool {
	log.Print(filepath.Join(fo.osPath, path))
	_, err := os.Stat(filepath.Join(fo.osPath, path))
	return !errors.Is(err, os.ErrNotExist)
}

func NewFileManager(dir string) filemanager.FileManager {
	return fileOperator{
		osPath: dir,
	}
}
