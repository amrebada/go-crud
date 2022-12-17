package database

import (
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
	instance *gorm.DB
}

var DB database

func NewDatabase() Database {
	if DB.instance == nil {
		DB = database{}
	}
	return &DB
}

func (d *database) Connect() (*gorm.DB, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/sandbox?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	d.instance = db
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return db, nil
}

func (d *database) Disconnect() error {
	db, err := d.instance.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (d *database) Migrate(entities ...interface{}) error {
	return d.instance.AutoMigrate(entities...)
}
