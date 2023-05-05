package usecase

import (
	"context"
	"testing"
	"url-shortener/internal/repository"
	"url-shortener/internal/schema"
)

func TestUseCase_Ping(t *testing.T) {
	ctx := context.Background()
	cfg := &repository.Config{
		DriverName:     "map",
		DataSourcePath: "test",
	}

	repo, err := repository.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	uc := New(repo)

	err = uc.Ping(ctx)
	if err != nil {
		t.Fatal(err)
	}

}

func TestUseCase_Batch(t *testing.T) {
	ctx := context.Background()

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

	_, err = uc.Batch(ctx, urls, "test", "test")
	if err != nil {
		t.Fatal(err)
	}

}

func TestUseCase_GetAllLinksByCookie(t *testing.T) {
	ctx := context.Background()
	cfg := &repository.Config{
		DriverName:     "map",
		DataSourcePath: "test",
	}

	repo, err := repository.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	uc := New(repo)

	_, err = uc.GetAllLinksByCookie(ctx, "test", "shor.t/")
	if err != nil {
		t.Fatal(err)
	}
}
