package main

import "database/sql"

type PostgresDataStore struct {
	Database *sql.DB
}
