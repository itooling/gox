package dbs

import (
	"log"
	"os"
	"time"

	"github.com/itooling/gox/oth"
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
		oth.CopyStruct(cn, conn)
	}

	if config != nil {
		oth.CopyStruct(cf, config)
	}

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
