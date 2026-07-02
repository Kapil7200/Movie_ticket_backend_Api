package database

import (
	"fmt"
	"movie_ticket/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)


	return db, nil
}
