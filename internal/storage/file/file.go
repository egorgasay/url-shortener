package filestorage

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
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

func (fs *FileStorage) MarkAsDeleted(shortURL, cookie string) {
	err := fs.OpenForWriteAt()
	if err != nil {
		log.Println("can't open a file ", err)
	}
	defer fs.Close()

	scanner := bufio.NewScanner(fs.File)
	var i int64 = 1
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " - ")
		if len(split) > 2 && split[1] == shortURL {
			lineWithDeletedMark := "0" + line[1:] + "\n"
			_, err = fs.File.WriteAt([]byte(lineWithDeletedMark), i-1)
			if err != nil {
				log.Println(err)
				return
			}
		}
		i += int64(1 + len(line))
	}
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

func (fs *FileStorage) OpenForWriteAt() error {
	fs.Mu.Lock()
	file, err := os.OpenFile(fs.Path, os.O_RDWR|os.O_CREATE, 0777)
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

		if len(split) > 2 && split[1] == shortURL {
			return split[2], nil
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

		if len(split) == 4 && split[3] == cookie && split[0] == "1" {
			URLs = append(URLs, schema.URL{LongURL: split[2], ShortURL: baseURL + split[1]})
		}
	}

	return URLs, nil
}

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
