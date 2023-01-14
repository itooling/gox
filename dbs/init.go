package dbs

import (
	"log"
	"os"
	"time"

	"github.com/itooling/gox"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db   *gorm.DB
	conf *gorm.Config
	conn Connection
)

const (
	Sqlite   = "lite"
	Mysql    = "mysql"
	Postgres = "pgsql"
)

type Connection struct {
	kind   string
	host   string
	port   int
	user   string
	pass   string
	dbname string
	prefix string
}

func init() {
	conn = Connection{
		kind:   gox.String("db.rdbms.kind"),
		host:   gox.String("db.rdbms.host"),
		port:   gox.Int("db.rdbms.port"),
		user:   gox.String("db.rdbms.user"),
		pass:   gox.String("db.rdbms.pass"),
		dbname: gox.String("db.rdbms.dbname"),
		prefix: gox.String("db.rdbms.prefix"),
	}

	conf = &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
		}),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   conn.prefix,
		},
		QueryFields: true,
	}

}

func DB() *gorm.DB {
	switch conn.kind {
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
