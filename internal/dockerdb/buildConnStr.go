package dockerdb

import "fmt"

func Build(conf CustomDB) (connStr string) {
	if conf.Vendor == "postgres" {
		connStr = fmt.Sprintf(
			"host=localhost user=%s password='%s' dbname=%s port=%s sslmode=disable",
			conf.DB.User, conf.DB.Password, conf.DB.Name, conf.Port)
	} else if conf.Vendor == "mysql" {
		connStr = fmt.Sprintf(
			"%s:%s@tcp(127.0.0.1:%s)/%s",
			conf.DB.User, conf.DB.Password, conf.Port, conf.DB.Name)
	}

	return connStr
}
