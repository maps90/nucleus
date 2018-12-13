package migration

import (
	_ "github.com/go-sql-driver/mysql" // required
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file" // required
	storage "github.com/maps90/nucleus/storage/mysql"
)

var defaultSeedLocation = "/resources/seeds"

// Migrate run DB migration
func Migrate(fileLocation string, dbName string, status bool) (v uint, d bool, err error) {
	sqlDB := storage.New(false)
	sqlDB.Multi = true

	db := sqlDB.MustConnect("master")
	driver, err := mysql.WithInstance(db.DB.DB, new(mysql.Config))
	if err != nil {
		return v, d, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+fileLocation,
		dbName, driver)
	if err != nil {
		return v, d, err
	}

	v, d, err = m.Version()

	if status {
		if err := m.Up(); err != nil {
			return v, d, err
		}
	} else {
		if err := m.Down(); err != nil {
			return v, d, err
		}
	}

	return
}
