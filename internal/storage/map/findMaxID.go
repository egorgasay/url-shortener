package mapstorage

func (s MapStorage) FindMaxID() (int, error) {
	return len(s.container), nil
}
