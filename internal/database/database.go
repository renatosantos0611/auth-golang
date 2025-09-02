package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type Service struct {
	DB *sql.DB
}

var (
	dbHost     = os.Getenv("DB_HOST")
	dbPort     = os.Getenv("DB_PORT")
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName     = os.Getenv("DB_NAME")
)

func New() *Service {
	// Valores padrão se não estiverem definidos nas variáveis de ambiente
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "admin"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "123"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "auth_golang"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL")

	return &Service{
		DB: db,
	}
}
