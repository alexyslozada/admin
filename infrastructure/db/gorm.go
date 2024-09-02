package db

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbGlobal  *gorm.DB
	errGlobal error
)

func NewGorm() (*gorm.DB, error) {
	sync.OnceFunc(func() {
		localhost := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslmode := os.Getenv("DB_SSLMODE")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", localhost, user, password, dbname, port, sslmode)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			errGlobal = err
			return
		}

		dbGlobal = db
	})()

	return dbGlobal, errGlobal
}
