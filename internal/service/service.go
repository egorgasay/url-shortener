package service

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type IService interface {
	CreateLink(string) (string, error)
	GetLink(string) (string, error)
}

type Service struct {
	IService
}

func NewService(db IService) *Service {
	if db == nil {
		panic("переменная storage равна nil")
	}

	return &Service{db}
}
