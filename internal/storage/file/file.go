package filestorage

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"url-shortener/internal/storage"
	shortener "url-shortener/pkg/api"
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
		return fmt.Errorf("open file error: %w", err)
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
func (fs *FileStorage) AddLink(ctx context.Context, longURL, shortURL, cookie string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

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
func (fs *FileStorage) FindMaxID(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

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
func (fs *FileStorage) GetLongLink(ctx context.Context, shortURL string) (longURL string, err error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

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
func (fs *FileStorage) GetAllLinksByCookie(ctx context.Context, cookie, baseURL string) ([]*shortener.UserURL, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	err := fs.Open()
	if err != nil {
		return nil, err
	}

	defer fs.Close()

	scanner := bufio.NewScanner(fs.File)
	var link = make([]*shortener.UserURL, 0)

	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " - ")

		if len(split) == 4 && split[3] == cookie && split[0] == "1" {
			link = append(link, &shortener.UserURL{OriginalUrl: split[2], ShortUrl: baseURL + split[1]})
		}
	}

	return link, nil
}

// Ping check for the presence of a file.
func (fs *FileStorage) Ping(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

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

// URLsCount returns the number of URLs in the file.
func (fs *FileStorage) URLsCount(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	err := fs.Open()
	if err != nil {
		return 0, err
	}
	defer fs.Close()
	scanner := bufio.NewScanner(fs.File)
	var count int
	for ; scanner.Scan(); count++ {
	}
	return count, nil
}

// UsersCount returns the number of users in the file.
func (fs *FileStorage) UsersCount(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	err := fs.Open()
	if err != nil {
		return 0, err
	}

	defer fs.Close()

	scanner := bufio.NewScanner(fs.File)
	var users = make(map[string]struct{}, 100)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, " - ")
		if len(split) > 3 {
			users[split[3]] = struct{}{}
		}
	}

	return len(users), nil
}
