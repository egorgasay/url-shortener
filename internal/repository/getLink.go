package repository

func (r Repository) GetLink(shortURL string) (longURL string, err error) {
	return r.repo.GetLongLink(shortURL)
}
