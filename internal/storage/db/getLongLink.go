package dbstorage

func (s RealStorage) GetLongLink(shortURL string) (longURL string, err error) {
	stm := s.DB.QueryRow("SELECT long FROM urls WHERE short = ?", shortURL)
	err = stm.Scan(&longURL)

	return longURL, err
}
