package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL驱动
)

type Database struct {
	*sql.DB
}

func NewDatabase(connectionString string) (*Database, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// 配置连接池参数
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	db.SetMaxOpenConns(25)
	
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	db.SetMaxIdleConns(10)
	
	// SetConnMaxLifetime 设置了连接可复用的最大时间
	// 超过这个时间的连接会在下次使用时被关闭并重新创建
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// 设置字符集为 utf8mb4，确保中文显示正常
	_, err = db.Exec("SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci")
	if err != nil {
		log.Printf("Warning: failed to set charset to utf8mb4: %v", err)
	}

	log.Println("Successfully connected to MySQL database with connection pool configured")
	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}
