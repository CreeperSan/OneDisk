package database

import (
	"OneDisk/lib/definition"
	"database/sql"
)

func Initialize() error {
	db, err := sql.Open("sqlite3", definition.PathDatabase)
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
