package filemanager

import "mime/multipart"

type FileManager interface {
	Copy(file *multipart.FileHeader, destinationPath string) error
	Read(path string) error
	Delete(path string) error
	Exists(path string) bool
}
