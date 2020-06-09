package component

import (
	"database/sql"
	"github.com/kayx-org/freja/env"
	"time"
)

type DriverName string

const (
	Postgres DriverName = "postgres"
	Mysql    DriverName = "mysql"
)

func openSql(dn DriverName, host string) (*sql.DB, error) {
	db, err := sql.Open(string(dn), host)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(env.GetEnvAsInt("DB_MAX_IDLE_CONN", 3))

	db.SetMaxOpenConns(env.GetEnvAsInt("DB_MAX_OPEN_CONN", 3))

	maxLife := env.GetEnvAsInt("DB_MAX_LIFETIME_CONN", 1)
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLife))

	return db, nil
}

func NewPSQL(host string) (*sql.DB, error) {
	return openSql(Postgres, host)
}

func NewMYSQL(host string) (*sql.DB, error) {
	return openSql(Mysql, host)
}
