package filestorage

import (
	"bufio"
	"url-shortener/internal/storage/shortenAlgorithm"
)

func (fs *FileStorage) AddLink(longURL string, id int) (string, error) {
	shortURL := shortenalgorithm.GetShortName(id)

	err := fs.Open()
	if err != nil {
		return "", err
	}

	defer fs.Close()

	writer := bufio.NewWriter(fs.File)

	_, err = writer.Write([]byte(shortURL + " - " + longURL + "\n"))
	if err != nil {
		return "", err
	}

	err = writer.Flush()
	if err != nil {
		return "", err
	}

	return shortURL, nil
}
