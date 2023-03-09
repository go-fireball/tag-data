package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type EnvironmentConfiguration struct {
	ConnectionString string
	PgDb             *PostgresDataStore
}

var (
	EnvConfig     *EnvironmentConfiguration
	EnvConfigOnce sync.Once
)

func GetEnvironmentConfig() *EnvironmentConfiguration {
	EnvConfigOnce.Do(func() {
		const (
			host     = "dbHost.com"
			port     = 5432
			user     = "username"
			password = "user-password"
			dbname   = "dbName"
		)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		var connectionString = psqlInfo
		var db, err = sql.Open("postgres", connectionString)
		db.SetConnMaxIdleTime(time.Minute)
		db.SetMaxIdleConns(20)
		db.SetConnMaxLifetime(time.Minute * 10)
		if err != nil {
			panic(err)
		}
		EnvConfig = &EnvironmentConfiguration{
			ConnectionString: connectionString,
			PgDb: &PostgresDataStore{
				Database: db,
			},
		}
	})
	return EnvConfig
}
