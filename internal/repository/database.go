package repository

import (
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() (*gorm.DB, error) {
	dsn := viper.GetString("database.dsn")
	if dsn == "" {
		return nil, ErrDatabaseDSNNotFound
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return DB, nil
}

func GetDB() *gorm.DB {
	return DB
}