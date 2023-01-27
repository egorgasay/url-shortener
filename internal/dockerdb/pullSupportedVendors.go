package dockerdb

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types"
	"log"
)

var DownloadableVendors = []string{
	"postgres",
	"mysql",
}

func Pull(vendor string) error {
	ctx := context.TODO()
	ddb, err := New(CustomDB{})
	if err != nil {
		return err
	}

	pull, err := ddb.cli.ImagePull(ctx, vendor, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(pull)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}

	err = pull.Close()
	if err != nil {
		return err
	}

	return nil
}

// PullAll скачивает все докер образы поддерживаемых вендоров
func PullAll() (err error) {
	for _, ven := range DownloadableVendors {
		err = Pull(ven)
		if err != nil {
			return err
		}
	}

	return nil
}
