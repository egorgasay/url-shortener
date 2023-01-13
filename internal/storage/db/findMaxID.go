package dbstorage

func (s RealStorage) FindMaxID() (int, error) {
	var id int

	stm := s.DB.QueryRow("SELECT MAX(id) FROM urls")
	err := stm.Scan(&id)

	return id, err
}
