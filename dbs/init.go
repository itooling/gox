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
		Kind:   sys.String("rdbms.kind"),
		Host:   sys.String("rdbms.host"),
		Port:   sys.Int("rdbms.port"),
		User:   sys.String("rdbms.user"),
		Pass:   sys.String("rdbms.pass"),
		Dbname: sys.String("rdbms.dbname"),
		Prefix: sys.String("rdbms.prefix"),
	}

	slow := sys.Int("rdbms.slow")
	if slow == 0 {
		slow = 300
	}
	cf = &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Duration(slow) * time.Millisecond,
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

func DBWithOption(conn *Connection, config *gorm.Config) *gorm.DB {
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
