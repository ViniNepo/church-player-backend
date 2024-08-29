package filemanager

type FileManager interface {
	Copy(file string, destinationPath string) error
	Read(path string) error
	Delete(path string) error
	Exists(path string) bool
}
