
package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	_ "time"
)

type Database struct {
	*sql.DB
}

func NewDatabase(connectionString string) (*Database, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to MySQL database")
	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}

