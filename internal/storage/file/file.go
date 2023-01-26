package filestorage

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"url-shortener/internal/schema"
	"url-shortener/internal/storage"
)

type FileStorage struct {
	Path string
	File *os.File
	Mu   sync.Mutex
}

const FileStorageType storage.Type = "file"

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

func (fs *FileStorage) AddLink(longURL, shortURL, cookie string) (string, error) {
	err := fs.Open()
	if err != nil {
		return "", err
	}

	defer fs.Close()

	writer := bufio.NewWriter(fs.File)

	_, err = writer.Write([]byte(shortURL + " - " + longURL + " - " + cookie + "\n"))
	if err != nil {
		return "", err
	}

	err = writer.Flush()
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

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

		if len(split) > 1 && split[0] == shortURL {
			return split[1], nil
		}
	}

	return longURL, errors.New("not found")
}

func (fs *FileStorage) GetAllLinksByCookie(cookie, baseURL string) ([]schema.URL, error) {
	err := fs.Open()
	if err != nil {
		return nil, err
	}

	defer fs.Close()

	scanner := bufio.NewScanner(fs.File)
	var URLs []schema.URL

	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " - ")

		if len(split) == 3 && split[2] == cookie {
			URLs = append(URLs, schema.URL{LongURL: split[1], ShortURL: baseURL + split[0]})
		}
	}

	return URLs, nil
}
