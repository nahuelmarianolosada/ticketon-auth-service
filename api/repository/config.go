package repository

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"ticketon-auth-service/api/model"
	"time"
)

var DB *gorm.DB

type DBConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

func loadEnv() (*DBConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}

	if config.Host == "" || config.User == "" || config.Password == "" || config.Name == "" || config.Port == "" {
		return nil, fmt.Errorf("incomplete database configuration")
	}

	return config, nil
}

func Connect() {
	config, err := loadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Name)

	log.Println("============================== dbUrl " + DBURL)

	DB, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure connection pooling
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database object from GORM: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to Database!")
}

func Migrate() {
	err := DB.AutoMigrate(&model.User{}, &model.Account{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database Migration Completed!")
}
