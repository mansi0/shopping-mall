package storage

import (
	"fmt"
	"log"
	"os"
	model "shopping-mall-customer/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

type Repository struct {
	DB *gorm.DB
}

func NewConnection(config *Config) (*gorm.DB, error) {
	// dsn := fmt.Sprintf(
	// 	"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	// 	config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	// )
	dsn := "postgresql://mansi:password@postgres:5432/shopping_mall?sslmode=disable"
	//fmt.Println(config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	err = model.MigrateCustomer(db)
	if err != nil {
		log.Print("could not migrate customer db")
		log.Fatal(err)
	}
	return db, err

}

func NewConnectionTesting(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)
	//dsn := "postgresql://mansi:password@postgres:5432/shopping_mall?sslmode=disable"
	//fmt.Println(config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	err = model.MigrateCustomer(db)
	if err != nil {
		log.Print("could not migrate customer db")
		log.Fatal(err)
	}
	return db, err

}

func FetchRepoTesting() *Repository {
	//fmt.Println(os.Getwd())
	err := godotenv.Load("storage/.env")
	if err != nil {

		log.Fatal(err)
	}

	config := &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}
	db, err := NewConnection(config)
	if err != nil {
		log.Fatal("Could not load the database")
	}
	r := &Repository{
		DB: db,
	}
	return r
}

func FetchRepo() *Repository {
	//fmt.Println(os.Getwd())
	err := godotenv.Load("../storage/.env")
	if err != nil {

		log.Fatal(err)
	}

	config := &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}
	db, err := NewConnection(config)
	if err != nil {
		log.Fatal("Could not load the database")
	}
	r := &Repository{
		DB: db,
	}
	return r
}

func InitDbForTesting() {

	log.Print(os.Getwd())

	err := godotenv.Load("../storage/.env")
	if err != nil {
		log.Fatal(err)
	}

	//log.Print("1")

	config := &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}
	db, err := NewConnectionTesting(config)
	if err != nil {
		log.Fatal(err)
	}
	db.Logger = logger.Default.LogMode(logger.Info)

	//log.Print("2")
}
