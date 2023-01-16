package dbs

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresInit() *gorm.DB {
	if db != nil {
		return db
	}

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		conn.Host, conn.Port, conn.User, conn.Pass, conn.Dbname)

	if d, err := gorm.Open(postgres.Open(dsn), conf); err != nil {
		panic(err)
	} else {
		db = d
	}
	return db
}
