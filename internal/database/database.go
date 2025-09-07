package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/G0tem/go-servise-entity/internal"
	"github.com/G0tem/go-servise-entity/internal/config"
	"github.com/G0tem/go-servise-entity/internal/model"
)

// Connect function
func Connect(cfg *config.Config) (*gorm.DB, error) {
	postgresPort := cfg.PostgresPort
	// because our config function returns a string, we are parsing our      str to int here
	port := internal.ParseInt(postgresPort, 5432)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDb,
		port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(cfg.LogLevel)),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "entity_",                         // table name prefix, table for `User` would be `t_users`
			SingularTable: true,                              // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   false,                             // skip the snake_casing of names
			NameReplacer:  strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("Connected")
	log.Println("running migrations")
	err = db.AutoMigrate(
		&model.Entity{},
	)
	if err != nil {
		log.Fatal("Failed to migrations. \n", err)
		return nil, err
	}
	log.Println("migrations complete")

	// Apply connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to create database connection pool. \n", err)
	}
	sqlDB.SetMaxIdleConns(cfg.PostgresMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.PostgresMaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.PostgresConnMaxLifetime)

	log.Println("Connection success")

	return db, nil
}
