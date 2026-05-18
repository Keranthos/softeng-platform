package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"softeng-platform/internal/config"
)

func main() {
	loadEnv()
	cfg := config.LoadConfig()
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		fmt.Println("打开数据库失败:", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("【连接失败】无法 Ping MySQL:", err)
		fmt.Println("请确认：1) MySQL 已启动 2) 已创建库与用户 3) .env 或环境变量中 DB_HOST / DB_PORT / DB_USER / DB_PASSWORD / DB_NAME 与 config 一致（默认库名 softeng、用户 softeng_app）。")
		os.Exit(1)
	}

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		fmt.Println("【查询失败】SHOW TABLES:", err)
		os.Exit(1)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			fmt.Println("读取结果失败:", err)
			os.Exit(1)
		}
		tables = append(tables, t)
	}

	if len(tables) == 0 {
		fmt.Println("【无表】已连上数据库，但当前库中没有任何表。")
		fmt.Println("请用有权限的账号执行本目录下 database/schema.sql，例如：")
		fmt.Println(`  mysql -u root -p < database/schema.sql`)
		os.Exit(2)
	}

	fmt.Printf("【正常】已连接，当前库共有 %d 张表：\n", len(tables))
	for _, t := range tables {
		fmt.Println("  -", t)
	}
	fmt.Println("演示数据：在项目 database 目录执行 seed_demo.sql（旧库可先执行 patch_existing_to_demo.sql）。")
}

func loadEnv() {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	_ = godotenv.Load(filepath.Join(wd, ".env"))
	_ = godotenv.Load(".env")
}
