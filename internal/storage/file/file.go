package filestorage

import (
	"os"
	"sync"
	"url-shortener/internal/storage"
)

// добавить mutex?

type FileStorage struct {
	Path string
	File *os.File
	Mu   sync.Mutex
}

func NewFileStorage(path string) storage.IStorage {
	return &FileStorage{Path: path}
}

func (fs *FileStorage) Open() error {
	fs.Mu.Lock()
	file, err := os.OpenFile(fs.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}

	fs.File = file
	return nil
}

func (fs *FileStorage) Close() error {
	fs.Mu.Unlock()
	return fs.File.Close()
}
