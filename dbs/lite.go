package dbs

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SqliteInit() *gorm.DB {
	if db != nil {
		return db
	}

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	path := "file::memory:?cache=shared"
	if conn.dbname != "" {
		path = conn.prefix + conn.dbname + ".db"
	}

	if d, err := gorm.Open(sqlite.Open(path), conf); err != nil {
		panic(err)
	} else {
		db = d
	}
	return db
}
