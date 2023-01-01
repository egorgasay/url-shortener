package storage

type IStorage interface {
	FindMaxID() (int, error)
	AddLink(longURL string, id int) (string, error)
	GetLongLink(shortURL string) (longURL string, err error)
}

type Storage struct {
	IStorage
}
