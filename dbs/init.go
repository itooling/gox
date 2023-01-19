package dbs

import (
	"log"
	"os"
	"time"

	"gitee.com/itooling/gox"
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
	Kind   string
	Host   string
	Port   int
	User   string
	Pass   string
	Dbname string
	Prefix string
}

func init() {
	conn = Connection{
		Kind:   gox.String("app.dbs.rdbms.kind"),
		Host:   gox.String("app.dbs.rdbms.host"),
		Port:   gox.Int("app.dbs.rdbms.port"),
		User:   gox.String("app.dbs.rdbms.user"),
		Pass:   gox.String("app.dbs.rdbms.pass"),
		Dbname: gox.String("app.dbs.rdbms.dbname"),
		Prefix: gox.String("app.dbs.rdbms.prefix"),
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
			TablePrefix:   conn.Prefix,
		},
		QueryFields: true,
	}

}

func DB(con *Connection, cfg *gorm.Config) *gorm.DB {
	if con != nil {
		if con.Kind != "" {
			conn.Kind = con.Kind
		}
		if con.Host != "" {
			conn.Host = con.Host
		}
		if con.Port != 0 {
			conn.Port = con.Port
		}
		if con.User != "" {
			conn.User = con.User
		}
		if con.Pass != "" {
			conn.Pass = con.Pass
		}
		if con.Dbname != "" {
			conn.Dbname = con.Dbname
		}
		if con.Prefix != "" {
			conn.Prefix = con.Prefix
		}
	}

	if cfg != nil {
		conf = cfg
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
