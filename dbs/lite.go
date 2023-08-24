package dbs

import (
	"fmt"

	"github.com/itooling/gorm-sqlite-cipher"
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
	if cn.Dbname != "" {
		path = cn.Prefix + cn.Dbname + ".db"
	}
	if cn.Pass != "" {
		path += fmt.Sprintf("?_pragma_key=%s", cn.Pass)
	}

	if d, err := gorm.Open(sqlcipher.Open(path), cf); err != nil {
		panic(err)
	} else {
		db = d
	}
	return db
}
