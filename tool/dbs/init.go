package dbs

import (
	"github.com/itooling/gox/tool"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
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
		kind:   tool.String("db.rdbms.kind"),
		host:   tool.String("db.rdbms.host"),
		port:   tool.Int("db.rdbms.port"),
		user:   tool.String("db.rdbms.user"),
		pass:   tool.String("db.rdbms.pass"),
		dbname: tool.String("db.rdbms.dbname"),
		prefix: tool.String("db.rdbms.prefix"),
	}

	conf = &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
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
