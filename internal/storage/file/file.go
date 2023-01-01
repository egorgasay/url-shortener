package filestorage

import "os"

// добавить mutex?

type FileStorage struct {
	Path string
	File *os.File
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{Path: path}
}

// если добавить мьютекс то сделать возможность выбирать для чего открывать
// файл, посредством аргумента в Open(globals.Read) где Read = os.O_RDONLY

func (fs *FileStorage) Open() error {
	file, err := os.OpenFile(fs.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}

	fs.File = file
	return nil
}

func (fs *FileStorage) Close() error {
	return fs.File.Close()
}
