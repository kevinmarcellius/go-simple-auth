package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"strconv"

	"github.com/dotenv-org/godotenvvault"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Port     int
	Postgres PostgresConfig
	JWTkey   string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadConfig() (*Config, error) {
	err := godotenvvault.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		Postgres: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
		},
		JWTkey: os.Getenv("JWT_SECRET"),
		Port: func() int {
			port, err := strconv.Atoi(os.Getenv("PORT"))
			if err != nil {
				return 0 // Default to 0 or handle the error as needed
			}
			return port
		}(),
	}
	return config, nil
}

func ConnectPostgres(cfg PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to Postgres: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get sql.DB from gorm DB: %v", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("Failed to ping Postgres: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to Postgres")
	return db, nil
}

func DBHealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
