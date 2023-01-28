package dockerdb

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"log"
)

// Init инициализирует docker контейнер с выбранной базой данных
func (ddb *DockerDB) Init(ctx context.Context) error {
	var env []string
	var portDocker nat.Port

	if ddb.Conf.Port == "" {
		return errors.New("config must be not empty")
	}

	switch ddb.Conf.Vendor {
	case "postgres":
		portDocker = "5432/tcp"
		env = []string{"POSTGRES_DB=" + ddb.Conf.DB.Name, "POSTGRES_USER=" + ddb.Conf.DB.User,
			"POSTGRES_PASSWORD=" + ddb.Conf.DB.Password}
	case "mysql":
		portDocker = "3306/tcp"
		env = []string{"MYSQL_DATABASE=" + ddb.Conf.DB.Name, "MYSQL_USER=" + ddb.Conf.DB.User,
			"MYSQL_ROOT_PASSWORD=" + ddb.Conf.DB.Password,
			"MYSQL_PASSWORD=" + ddb.Conf.DB.Password}
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			portDocker: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: ddb.Conf.Port,
				},
			},
		},
	}

	containerName := ddb.Conf.DB.Name
	r, err := ddb.cli.ContainerCreate(ctx, &container.Config{
		Image: ddb.Conf.Vendor,
		Env:   env,
	}, hostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}

	ddb.ID = r.ID
	log.Println(ddb.ID)

	return nil
}
