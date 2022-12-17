package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	Connect() (*gorm.DB, error)
	Disconnect() error
	Migrate(entities ...interface{}) error
}

type database struct {
	Instance *gorm.DB
}

var DB database

func NewDatabase() Database {
	if DB.Instance == nil {
		DB = database{}
	}
	return &DB
}

func (d *database) Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/sandbox?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	d.Instance = db
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return db, nil
}

func (d *database) Disconnect() error {
	db, err := d.Instance.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (d *database) Migrate(entities ...interface{}) error {
	return d.Instance.AutoMigrate(entities...)
}
