package component

import (
	"database/sql"
	"os"
	"strconv"
	"time"
)

type DriverName string

const (
	Postgres DriverName = "postgres"
	Mysql               = "mysql"
)

func openSql(dn DriverName, host string) (*sql.DB, error) {
	db, err := sql.Open(string(dn), host)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(getEnvAsInt("DB_MAX_IDLE_CONN", 3))

	db.SetMaxOpenConns(getEnvAsInt("DB_MAX_OPEN_CONN", 3))

	maxLife := getEnvAsInt("DB_MAX_LIFETIME_CONN", 1)
	db.SetConnMaxLifetime(time.Second * time.Duration(maxLife))

	return db, nil
}

func NewPSQL(host string) (*sql.DB, error) {
	return openSql(Postgres, host)
}

func NewMYSQL(host string) (*sql.DB, error) {
	return openSql(Mysql, host)
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
