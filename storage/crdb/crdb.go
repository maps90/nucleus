package crdb

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // use by sqlx
	"github.com/spf13/cast"
)

var (
	dbConnection *SQLConn
	config       *Config
)

type Config struct {
	Host        string
	DB          string
	Params      string
	LogMode     bool
	MaxLifetime int64
	MaxOpen     int
	MaxIdle     int
}

// SQLConn struct
type SQLConn struct {
	sqlx    *SQLDB
	Multi   bool
	LogMode bool
}

// New : initialize new DB connections, if dbConnection is already exists, returns that instead
func New() *SQLConn {
	if dbConnection == nil {
		dbConnection = &SQLConn{}
	}

	if config == nil {
		config = &Config{
			os.Getenv("SQL_HOST"),
			os.Getenv("SQL_DB"),
			os.Getenv("SQL_PARAMS"),
			cast.ToBool(os.Getenv("SQL_LOG")),
			cast.ToInt64(os.Getenv("SQL_MAX_LIFETIME")),
			cast.ToInt(os.Getenv("SQL_MAX_OPEN")),
			cast.ToInt(os.Getenv("SQL_MAX_IDLE")),
		}
	}

	return dbConnection
}

// MustConnect retrieve MySQL established connection client (sqlx) and panic if error
func (db *SQLConn) MustConnect() *SQLDB {
	if d, err := db.Connect(); err != nil {
		panic(err)
	} else {
		return d
	}
}

// Connect retrieve MySQL established connection client (sqlx)
func (db *SQLConn) Connect() (*SQLDB, error) {
	con, err := db.openConnection()
	if err != nil {
		return nil, err
	}

	return con, nil
}

// Shutdown disconnecting all established mysql client connection
func Shutdown() (err error) {
	if dbConnection == nil {
		return nil
	}

	return dbConnection.sqlx.Close()
}

func (db *SQLConn) openConnection() (*SQLDB, error) {
	dsn := fmt.Sprintf(
		"postgresql://%v/%v?%v",
		config.Host,
		config.DB,
		config.Params,
	)
	if db.Multi {
		dsn += "&multiStatements=true"
	}

	con, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	newConn := &SQLDB{DB: con, logMode: config.LogMode}

	con.SetConnMaxLifetime(time.Duration(config.MaxLifetime))
	con.SetMaxOpenConns(config.MaxOpen)
	con.SetMaxIdleConns(config.MaxIdle)

	return newConn, nil
}
