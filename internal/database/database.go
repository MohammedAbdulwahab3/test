package database

import (
 	"database/sql"
 	"fmt"
 	"time"

 	"gorm.io/driver/postgres"
 	"gorm.io/driver/sqlite"
 	"gorm.io/gorm"
 	"gorm.io/gorm/logger"
 	"gorm.io/gorm/schema"

 	"flutter_go_crud/backend/internal/config"
 	"flutter_go_crud/backend/internal/models"
)

func Connect(cfg config.Config) (*gorm.DB, func(), error) {
 	var (
 		db  *gorm.DB
 		err error
 	)

 	gcfg := &gorm.Config{
 		Logger:                 logger.Default.LogMode(logger.Warn),
 		PrepareStmt:            true,
 		SkipDefaultTransaction: true,
 		NamingStrategy: schema.NamingStrategy{
 			SingularTable: false,
 		},
 	}

 	if cfg.DBDriver == "postgres" {
 		db, err = gorm.Open(postgres.Open(cfg.DBDSN), gcfg)
 	} else {
 		db, err = gorm.Open(sqlite.Open(cfg.DBDSN), gcfg)
 	}
 	if err != nil {
 		return nil, func() {}, fmt.Errorf("open database: %w", err)
 	}

 	sqlDB, err := db.DB()
 	if err != nil {
 		return nil, func() {}, fmt.Errorf("db(): %w", err)
 	}

 	// tuned pool defaults for fast dev
 	sqlDB.SetMaxOpenConns(25)
 	sqlDB.SetMaxIdleConns(25)
 	sqlDB.SetConnMaxLifetime(30 * time.Minute)
 	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

 	closeFn := func() {
 		_ = sqlDB.Close()
 	}

 	return db, closeFn, nil
}

func AutoMigrate(db *gorm.DB) error {
 	return db.AutoMigrate(
 		&models.Item{},
 	)
}

func GetSQL(db *gorm.DB) *sql.DB {
 	s, _ := db.DB()
 	return s
}
