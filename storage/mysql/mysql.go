// Package mysql storage is designed to give lazy load singleton access to mysql connections
// it doesn't provide any cluster nor balancing support, assuming it is handled
// in lower level infra, i.e. proxy, cluster etc.
package mysql

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // use by sqlx
	"github.com/jmoiron/sqlx"
	"github.com/maps90/nucleus/config"
	"github.com/maps90/nucleus/util"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var dbConnection *SQLConn

// SQLConn struct
type SQLConn struct {
	sqlx   map[string]*SQLDB
	Mock   sqlmock.Sqlmock
	mux    sync.Mutex
	sqlDB  *SQLDB
	Multi  bool
	isMock bool
}

// Config struct
type Config struct {
	User        string
	Password    string
	Address     string
	DB          string
	LogMode     bool
	MaxOpen     int
	MaxIdle     int
	MaxLifetime int
	Fallback    string
}

// Get connection in thread-safe fashion
func (db *SQLConn) Get(id string) *SQLDB {
	db.mux.Lock()
	defer db.mux.Unlock()
	if conn, ok := db.sqlx[id]; ok {
		return conn
	}
	return nil
}

// Set connection in thread-safe fashion
func (db *SQLConn) Set(id string, sqlx *SQLDB) {
	db.mux.Lock()
	db.sqlx[id] = sqlx
	db.mux.Unlock()
}

// New : initialize new DB connections, if dbConnection is already exists, returns that instead
func New(isMock bool) *SQLConn {
	if dbConnection == nil {
		dbConnection = &SQLConn{sqlx: make(map[string]*SQLDB)}
	}

	if isMock {
		d, mock, err := sqlmock.New()
		if err != nil {
			return nil
		}

		db := sqlx.NewDb(d, "mysql")
		return &SQLConn{Mock: mock, sqlDB: &SQLDB{DB: db}, isMock: true}
	}

	return dbConnection
}

// MustConnect retrieve MySQL established connection client (sqlx) and panic if error
func (db *SQLConn) MustConnect(id string) *SQLDB {
	if d, err := db.Connect(id); err != nil {
		panic(err)
	} else {
		return d
	}
}

// Connect retrieve MySQL established connection client (sqlx)
func (db *SQLConn) Connect(id string) (*SQLDB, error) {
	mysqlConfig := config.GetStringMap("mysql")
	if _, ok := mysqlConfig[id]; !ok {
		return nil, fmt.Errorf("mysql configuration for [%s] does not exists", id)
	}

	if db.Multi {
		return db.openConnection(id)
	}

	if db.isMock {
		return db.sqlDB, nil
	}

	// if previously established, reuse and ping
	if con := db.Get(id); con != nil {
		return con, nil
	}

	con, err := db.openConnection(id)
	if err != nil {
		return nil, err
	}

	// otherwise establish new connection through centralized component config
	db.Set(id, con)
	return db.Get(id), nil
}

// Shutdown disconnecting all established mysql client connection
func Shutdown() (err error) {
	if dbConnection == nil {
		return nil
	}
	for _, c := range dbConnection.sqlx {
		err = c.Close()
	}
	return err
}

func (db *SQLConn) openConnection(id string) (*SQLDB, error) {
	opt := setupConfig(id)
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local",
		opt.User,
		opt.Password,
		opt.Address,
		opt.DB,
	)

	if db.Multi {
		dsn += "&multiStatements=true"
	}

	con, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		if opt.Fallback == id { // prevent endless loop
			return db.Get(id), err
		}
		return db.Connect(opt.Fallback)
	}

	con.SetConnMaxLifetime(time.Duration(opt.MaxLifetime))
	con.SetMaxOpenConns(opt.MaxOpen)
	con.SetMaxIdleConns(opt.MaxIdle)

	db.Set(id, &SQLDB{DB: con, logMode: opt.LogMode})
	return db.Get(id), nil
}

func setupConfig(id string) *Config {
	option := &Config{
		User:     config.GetString(getKey(id, "user")),
		Password: config.GetString(getKey(id, "password")),
		Address:  config.GetString(getKey(id, "address")),
		DB:       config.GetString(getKey(id, "db")),
		Fallback: config.GetString(getKey(id, "fallback_to")),
	}

	option.MaxLifetime = util.DefaultInt(config.GetInt("mysql.max_lifetime"), 30)
	option.MaxOpen = util.DefaultInt(config.GetInt("mysql.max_open"), 30)
	option.LogMode = config.GetBool("mysql.log")

	return option
}

func getKey(id, types string) string {
	return fmt.Sprintf("mysql.%s.%s", id, types)
}
