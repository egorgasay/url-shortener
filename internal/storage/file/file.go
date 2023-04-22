package filestorage

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
)

var (
	_ storage.IStorage = (*FileStorage)(nil)
)

// FileStorage struct with Mu for concurrent uses, File instance and Path variable.
// It has methods for working with URLs.
type FileStorage struct {
	Path string
	File *os.File
	Mu   sync.Mutex
}

// FileStorageType type for file storage.
const FileStorageType storage.Type = "file"

// NewFileStorage FileStorage struct constructor.
func NewFileStorage(path string) (storage.IStorage, error) {
	fs := &FileStorage{Path: path}
	err := fs.Open()
	if err != nil {
		return nil, err
	}

	err = fs.Close()
	if err != nil {
		return nil, err
	}

	return fs, nil
}

// Open opens the file.
func (fs *FileStorage) Open() error {
	fs.Mu.Lock()
	file, err := os.OpenFile(fs.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}

	fs.File = file
	return nil
}

// OpenForWriteAt opens the file for write at.
func (fs *FileStorage) OpenForWriteAt() error {
	fs.Mu.Lock()
	file, err := os.OpenFile(fs.Path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	fs.File = file
	return nil
}

// Close closes the file.
func (fs *FileStorage) Close() error {
	fs.Mu.Unlock()
	return fs.File.Close()
}

// AddLink adds a link to the file.
func (fs *FileStorage) AddLink(longURL, shortURL, cookie string) (string, error) {
	err := fs.Open()
	if err != nil {
		return "", err
	}

	defer fs.Close()

	writer := bufio.NewWriter(fs.File)

	_, err = writer.Write([]byte("1" + " - " + shortURL + " - " + longURL + " - " + cookie + "\n"))
	if err != nil {
		return "", err
	}

	err = writer.Flush()
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

// FindMaxID gets len of the file.
func (fs *FileStorage) FindMaxID() (int, error) {
	err := fs.Open()
	if err != nil {
		return 0, err
	}

	defer fs.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := fs.File.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case errors.Is(err, io.EOF):
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// GetLongLink gets a long link from the file.
func (fs *FileStorage) GetLongLink(shortURL string) (longURL string, err error) {
	err = fs.Open()
	if err != nil {
		return "", err
	}

	defer fs.Close()

	scanner := bufio.NewScanner(fs.File)

	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " - ")

		if len(split) > 2 && split[1] == shortURL {
			return split[2], nil
		}
	}

	return longURL, errors.New("not found")
}

// GetAllLinksByCookie gets all links ([]schema.URL) by cookie.
func (fs *FileStorage) GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error) {
	err := fs.Open()
	if err != nil {
		return nil, err
	}

	defer fs.Close()

	scanner := bufio.NewScanner(fs.File)
	var link []schema.URL

	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " - ")

		if len(split) == 4 && split[3] == cookie && split[0] == "1" {
			link = append(link, schema.URL{LongURL: split[2], ShortURL: baseURL + split[1]})
		}
	}

	return link, nil
}

// Ping check for the presence of a file.
func (fs *FileStorage) Ping() error {
	err := fs.Open()
	if err != nil {
		return err
	}

	err = fs.Close()
	if err != nil {
		return err
	}

	return nil
}

// MarkAsDeleted finds a URL and marks it as deleted.
func (fs *FileStorage) MarkAsDeleted(shortURL, cookie string) error {
	err := fs.OpenForWriteAt()
	if err != nil {
		return fmt.Errorf("can't open a file %w", err)
	}
	defer fs.Close()

	scanner := bufio.NewScanner(fs.File)
	var i int64 = 1
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " - ")
		if len(split) > 2 && split[1] == shortURL && split[3] == cookie {
			lineWithDeletedMark := "0" + line[1:] + "\n"
			_, err = fs.File.WriteAt([]byte(lineWithDeletedMark), i-1)
			if err != nil {
				return fmt.Errorf("can't write a file %w", err)
			}
		}
		i += int64(1 + len(line))
	}
	return nil
}

// Shutdown closes the file.
func (fs *FileStorage) Shutdown() error {
	err := fs.Close()
	if err != nil {
		return fmt.Errorf("can't close a file %w", err)
	}

	return nil
}
