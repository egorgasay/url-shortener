package dockerdb

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

func (ddb *VDB) Setup(strConn string) (*sql.DB, string) {
	ctx := context.TODO()
	if ddb.ID == "" {
		err := ddb.Init(ctx)
		if err != nil {
			log.Fatal("Init: ", err)
		}
	}

	ctx = context.TODO()
	err := ddb.Run(ctx)
	if err != nil {
		log.Fatal("Run: ", err)
	}

	if strConn == "" {
		strConn = Build(ddb.Conf)
	}

	db, err := ddb.getDB(strConn)

	if err != nil {
		log.Fatal("ping: ", err)
	}

	return db, strConn
}

func (ddb *VDB) getDB(connStr string) (*sql.DB, error) {
	after := time.After(maxWaitTime)
	ticker := time.NewTicker(tryInterval)
	for {
		select {
		case <-after:
			return nil, errors.New("timeout")
		default:
			db, err := sql.Open(ddb.Conf.Vendor, connStr)
			pingErr := db.Ping()
			if pingErr == nil && err == nil {
				return db, nil
			}
			<-ticker.C
		}
	}
}
