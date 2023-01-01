package filestorage

import (
	"bufio"
	"errors"
	"strings"
)

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
