package db

import (
	"database/sql"

	"github.com/wmentor/log"

	_ "github.com/lib/pq"
)

var (
	conStr string
)

func Get() (*sql.DB, error) {

	if con, err := sql.Open("postgres", conStr); err == nil {
		if err = con.Ping(); err != nil {
			con.Close()
			return nil, err
		}
		return con, nil
	} else {
		log.Error(err.Error())
		return nil, err
	}
}

func SetConnectString(cfg string) {
	conStr = cfg
}
