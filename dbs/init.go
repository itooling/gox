package dbs

import (
	"log"
	"os"
	"reflect"
	"time"

	"github.com/itooling/gox/sys"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db *gorm.DB
	cf *gorm.Config
	cn *Connection
)

const (
	Sqlite   = "lite"
	Mysql    = "mysql"
	Postgres = "pgsql"
)

type Connection struct {
	Kind   string
	Host   string
	Port   int
	User   string
	Pass   string
	Dbname string
	Prefix string
}

func init() {
	cn = &Connection{
		Kind:   sys.String("app.dbs.rdbms.kind"),
		Host:   sys.String("app.dbs.rdbms.host"),
		Port:   sys.Int("app.dbs.rdbms.port"),
		User:   sys.String("app.dbs.rdbms.user"),
		Pass:   sys.String("app.dbs.rdbms.pass"),
		Dbname: sys.String("app.dbs.rdbms.dbname"),
		Prefix: sys.String("app.dbs.rdbms.prefix"),
	}

	cf = &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             300 * time.Millisecond,
			LogLevel:                  logger.Warn,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
		}),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   cn.Prefix,
		},
		QueryFields: true,
	}
}

func DB() *gorm.DB {
	switch cn.Kind {
	case Sqlite:
		SqliteInit()
	case Mysql:
		MysqlInit()
	case Postgres:
		PostgresInit()
	default:
		SqliteInit()
	}
	return db
}

func DBS(conn *Connection, config *gorm.Config) *gorm.DB {
	if conn != nil {
		CopyStruct(cn, conn)
	}

	if config != nil {
		CopyStruct(cf, config)
	}

	switch conn.Kind {
	case Sqlite:
		SqliteInit()
	case Mysql:
		MysqlInit()
	case Postgres:
		PostgresInit()
	default:
		SqliteInit()
	}
	return db
}

func CopyStruct(dst, src interface{}) {
	dstValue := reflect.ValueOf(dst).Elem()
	srcValue := reflect.ValueOf(src).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		srcName := srcValue.Type().Field(i).Name
		dstFieldByName := dstValue.FieldByName(srcName)

		if dstFieldByName.IsValid() {
			switch dstFieldByName.Kind() {
			case reflect.Ptr:
				switch srcField.Kind() {
				case reflect.Ptr:
					if srcField.IsNil() {
						dstFieldByName.Set(reflect.New(dstFieldByName.Type().Elem()))
					} else {
						dstFieldByName.Set(srcField)
					}
				default:
					dstFieldByName.Set(srcField.Addr())
				}
			default:
				switch srcField.Kind() {
				case reflect.Ptr:
					if srcField.IsNil() {
						dstFieldByName.Set(reflect.Zero(dstFieldByName.Type()))
					} else {
						dstFieldByName.Set(srcField.Elem())
					}
				default:
					dstFieldByName.Set(srcField)
				}
			}
		}
	}
}
