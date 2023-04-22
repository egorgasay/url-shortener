package usecase

import (
	"testing"
	"url-shortener/internal/repository"
	"url-shortener/internal/schema"
)

func TestUseCase_Ping(t *testing.T) {
	cfg := &repository.Config{
		DriverName:     "map",
		DataSourcePath: "test",
	}

	repo, err := repository.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	uc := New(repo)

	err = uc.Ping()
	if err != nil {
		t.Fatal(err)
	}

}

func TestUseCase_Batch(t *testing.T) {
	cfg := &repository.Config{
		DriverName:     "map",
		DataSourcePath: "test",
	}

	repo, err := repository.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	uc := New(repo)

	urls := []schema.BatchURL{
		{
			Chars:    "test",
			Original: "test",
		},

		{
			Chars:    "qwd",
			Original: "vk.com/gasayminajj",
		},
	}

	_, err = uc.Batch(urls, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

}

func TestUseCase_GetAllLinksByCookie(t *testing.T) {
	cfg := &repository.Config{
		DriverName:     "map",
		DataSourcePath: "test",
	}

	repo, err := repository.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	uc := New(repo)

	_, err = uc.GetAllLinksByCookie("test", "shor.t/")
	if err != nil {
		t.Fatal(err)
	}
}
