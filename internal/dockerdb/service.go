// Package dockerdb allows user to create virtual databases using docker.
package dockerdb

import (
	"github.com/docker/docker/client"
	"time"
)

const maxWaitTime = 20 * time.Second
const tryInterval = 1 * time.Second

type VDB struct {
	ID   string
	cli  *client.Client
	Conf CustomDB
}

type DB struct {
	Name     string
	User     string
	Password string
}

type CustomDB struct {
	DB     DB
	Port   string
	Vendor string
}

func New(conf CustomDB) (*VDB, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv,
		client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &VDB{cli: cli, Conf: conf}, nil
}
