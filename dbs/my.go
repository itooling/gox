package dbs

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlInit() *gorm.DB {
	if db != nil {
		return db
	}

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conn.User, conn.Pass, conn.Host, conn.Port, conn.Dbname)

	if d, err := gorm.Open(mysql.Open(dsn), conf); err != nil {
		panic(err)
	} else {
		db = d
	}
	return db
}
