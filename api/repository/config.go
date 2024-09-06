package repository

import (
	"ticketon-auth-service/api/model"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB
var dbError error

func Connect() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	//root:mysql@tcp(localhost:33060)/digitalmoney-db?parseTime=true
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	fmt.Println("============================== dbUrl + ", DBURL)
	DB, dbError = gorm.Open(mysql.Open(DBURL), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB ")
	}
	log.Println("Connected to Database!")
}

func Migrate() {
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Account{})
	log.Println("Database Migration Completed!")
}
