package storage

func (s MapStorage) FindMaxID() (int, error) {
	return len(s), nil
}

func (s RealStorage) FindMaxID() (int, error) {
	var id int

	stm := s.QueryRow("SELECT MAX(id) FROM urls")
	err := stm.Scan(&id)

	return id, err
}
