package filestorage

import (
	"bytes"
	"errors"
	"io"
)

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
