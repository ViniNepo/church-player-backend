package fileoperator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/ViniNepo/church-player-backend/internal/filemanager"
)

type fileOperator struct{}

func (fo fileOperator) Copy(path string, destination string) error {
	origin, err := os.OpenFile(path, os.O_RDONLY, 0o777)
	if err != nil {
		return fmt.Errorf("opening origin file: %w", err)
	}
	defer origin.Close()

	info, err := origin.Stat()
	if err != nil {
		return fmt.Errorf("getting file info: %w", err)
	}

	dest, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, 0o777)
	if err != nil {
		return fmt.Errorf("opening destination file: %w", err)
	}
	defer dest.Close()

	// copy contents of origin blob to destination file
	b, err := io.Copy(dest, origin)
	if err != nil {
		return fmt.Errorf("copy contents of %s: %w", origin.Name(), err)
	}

	// if the amount of copied bytes is not equal to the size of the blob file then we return error
	if info.Size() != b {
		return errors.New("copied bytes not equal to origin file size")
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
