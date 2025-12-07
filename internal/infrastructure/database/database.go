package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"canteen/internal/infrastructure/config"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	User     string
	Password string
	Network  string
	Host     string
	Port     int
	DBName   string
}

// NewDB creates a new DB connection with provided config
func NewDB(config DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Network, config.Host, config.Port, config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Failed to open DB: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping DB: %v", err)
		return nil, err
	}

	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established")
	return db, nil
}

// LoadDBConfig 从配置中读取数据库配置
func LoadDBConfig() DBConfig {
	return DBConfig{
		User:     config.GetString("database.user"),
		Password: config.GetString("database.pwd"),
		Network:  "tcp",
		Host:     config.GetString("database.host"),
		Port:     config.GetInt("database.port"),
		DBName:   config.GetString("database.dbname"),
	}
}

func InitDb() *sql.DB {
	config := LoadDBConfig()
	db, err := NewDB(config)
	if err != nil {
		log.Fatalf("Database init failed: %v", err)
	}
	return db
}