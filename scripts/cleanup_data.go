package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"softeng-platform/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("=========================================")
	fmt.Println("数据清理脚本")
	fmt.Println("=========================================")
	fmt.Println("此脚本将执行以下操作：")
	fmt.Println("1. 清空所有用户的所有提交记录（工具、课程、项目）")
	fmt.Println("2. 删除所有名为'豆包'的工具及其相关资源")
	fmt.Println("=========================================")
	fmt.Print("确认执行？(yes/no): ")

	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "yes" {
		fmt.Println("操作已取消")
		return
	}

	// 加载配置
	cfg := config.LoadConfig()
	
	// 打印数据库连接信息（隐藏密码）
	fmt.Printf("\n数据库连接信息: %s@%s/%s\n", 
		func() string {
			if user := os.Getenv("DB_USER"); user != "" {
				return user
			}
			return "softeng_app"
		}(),
		func() string {
			if host := os.Getenv("DB_HOST"); host != "" {
				return host
			}
			return "127.0.0.1"
		}(),
		func() string {
			if dbname := os.Getenv("DB_NAME"); dbname != "" {
				return dbname
			}
			return "softeng"
		}())

	// 连接数据库
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v\n提示: 请检查环境变量 DB_USER, DB_PASSWORD, DB_HOST, DB_NAME 是否正确", err)
	}

	fmt.Println("数据库连接成功")

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("开始事务失败: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("发生错误，已回滚: %v", r)
		}
	}()

	// 1. 删除所有名为"豆包"的工具及其相关资源
	fmt.Println("\n[1/2] 正在删除所有名为'豆包'的工具...")
	if err := deleteToolsByName(tx, "豆包"); err != nil {
		tx.Rollback()
		log.Fatalf("删除'豆包'工具失败: %v", err)
	}
	fmt.Println("✓ 已删除所有名为'豆包'的工具")

	// 2. 清空所有用户的提交记录
	fmt.Println("\n[2/2] 正在清空所有用户的提交记录...")
	if err := clearAllUserSubmissions(tx); err != nil {
		tx.Rollback()
		log.Fatalf("清空用户提交记录失败: %v", err)
	}
	fmt.Println("✓ 已清空所有用户的提交记录")

	// 提交事务
	if err := tx.Commit(); err != nil {
		log.Fatalf("提交事务失败: %v", err)
	}

	fmt.Println("\n=========================================")
	fmt.Println("数据清理完成！")
	fmt.Println("=========================================")
}

// deleteToolsByName 删除所有指定名称的工具及其相关资源
func deleteToolsByName(tx *sql.Tx, toolName string) error {
	// 首先查询所有名为"豆包"的工具ID
	rows, err := tx.Query("SELECT resource_id FROM tools WHERE resource_name = ?", toolName)
	if err != nil {
		return fmt.Errorf("查询工具失败: %v", err)
	}
	defer rows.Close()

	var toolIDs []int
	for rows.Next() {
		var toolID int
		if err := rows.Scan(&toolID); err != nil {
			continue
		}
		toolIDs = append(toolIDs, toolID)
	}

	if len(toolIDs) == 0 {
		fmt.Printf("  未找到名为'%s'的工具\n", toolName)
		return nil
	}

	fmt.Printf("  找到 %d 个名为'%s'的工具\n", len(toolIDs), toolName)

	// 删除相关资源（由于外键CASCADE，tool_images, tool_tags, tool_contributors会自动删除）
	// 但需要手动删除其他关联数据

	// 删除收藏记录
	for _, toolID := range toolIDs {
		_, err := tx.Exec("DELETE FROM collections WHERE resource_type = 'tool' AND resource_id = ?", toolID)
		if err != nil {
			return fmt.Errorf("删除收藏记录失败 (tool_id=%d): %v", toolID, err)
		}
	}

	// 删除点赞记录
	for _, toolID := range toolIDs {
		_, err := tx.Exec("DELETE FROM likes WHERE resource_type = 'tool' AND resource_id = ?", toolID)
		if err != nil {
			return fmt.Errorf("删除点赞记录失败 (tool_id=%d): %v", toolID, err)
		}
	}

	// 删除评论（评论的点赞会通过外键CASCADE自动删除）
	for _, toolID := range toolIDs {
		_, err := tx.Exec("DELETE FROM comments WHERE resource_type = 'tool' AND resource_id = ?", toolID)
		if err != nil {
			return fmt.Errorf("删除评论失败 (tool_id=%d): %v", toolID, err)
		}
	}

	// 删除状态日志
	for _, toolID := range toolIDs {
		_, err := tx.Exec("DELETE FROM resource_status_logs WHERE resource_type = 'tool' AND resource_id = ?", toolID)
		if err != nil {
			return fmt.Errorf("删除状态日志失败 (tool_id=%d): %v", toolID, err)
		}
	}

	// 最后删除工具本身（会级联删除 tool_images, tool_tags, tool_contributors）
	result, err := tx.Exec("DELETE FROM tools WHERE resource_name = ?", toolName)
	if err != nil {
		return fmt.Errorf("删除工具失败: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("  已删除 %d 个工具及其所有相关资源\n", rowsAffected)
	return nil
}

// clearAllUserSubmissions 清空所有用户的提交记录
func clearAllUserSubmissions(tx *sql.Tx) error {
	// 1. 删除所有工具提交记录
	fmt.Println("  正在删除工具提交记录...")
	
	// 先查询所有有提交者的工具ID
	toolRows, err := tx.Query("SELECT resource_id FROM tools WHERE submitter_id IS NOT NULL")
	if err != nil {
		return fmt.Errorf("查询工具ID失败: %v", err)
	}
	defer toolRows.Close()

	var toolIDs []int
	for toolRows.Next() {
		var toolID int
		if err := toolRows.Scan(&toolID); err != nil {
			continue
		}
		toolIDs = append(toolIDs, toolID)
	}

	if len(toolIDs) > 0 {
		// 构建IN子句
		placeholders := ""
		args := make([]interface{}, len(toolIDs))
		for i, id := range toolIDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args[i] = id
		}

		// 删除工具相关的关联数据
		_, err = tx.Exec(fmt.Sprintf("DELETE FROM collections WHERE resource_type = 'tool' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除工具收藏记录失败: %v", err)
		}

		_, err = tx.Exec(fmt.Sprintf("DELETE FROM likes WHERE resource_type = 'tool' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除工具点赞记录失败: %v", err)
		}

		_, err = tx.Exec(fmt.Sprintf("DELETE FROM comments WHERE resource_type = 'tool' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除工具评论失败: %v", err)
		}

		_, err = tx.Exec(fmt.Sprintf("DELETE FROM resource_status_logs WHERE resource_type = 'tool' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除工具状态日志失败: %v", err)
		}

		// 删除工具（会级联删除 tool_images, tool_tags, tool_contributors）
		result, err := tx.Exec(fmt.Sprintf("DELETE FROM tools WHERE resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除工具提交记录失败: %v", err)
		}
		toolsDeleted, _ := result.RowsAffected()
		fmt.Printf("  已删除 %d 条工具提交记录\n", toolsDeleted)
	} else {
		fmt.Printf("  没有找到工具提交记录\n")
	}

	// 2. 删除所有课程提交记录
	fmt.Println("  正在删除课程提交记录...")
	
	// 先查询所有有贡献者的课程ID
	courseRows, err := tx.Query("SELECT DISTINCT course_id FROM course_contributors")
	if err != nil {
		return fmt.Errorf("查询课程ID失败: %v", err)
	}
	defer courseRows.Close()

	var courseIDs []int
	for courseRows.Next() {
		var courseID int
		if err := courseRows.Scan(&courseID); err != nil {
			continue
		}
		courseIDs = append(courseIDs, courseID)
	}

	if len(courseIDs) > 0 {
		// 构建IN子句
		placeholders := ""
		args := make([]interface{}, len(courseIDs))
		for i, id := range courseIDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args[i] = id
		}

		// 删除课程相关的关联数据
		_, err = tx.Exec(fmt.Sprintf("DELETE FROM collections WHERE resource_type = 'course' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除课程收藏记录失败: %v", err)
		}

		_, err = tx.Exec(fmt.Sprintf("DELETE FROM likes WHERE resource_type = 'course' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除课程点赞记录失败: %v", err)
		}

		_, err = tx.Exec(fmt.Sprintf("DELETE FROM comments WHERE resource_type = 'course' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除课程评论失败: %v", err)
		}

		// 删除课程（会级联删除 course_teachers, course_categories, course_resources_web, course_resources_upload, course_contributors）
		result, err := tx.Exec(fmt.Sprintf("DELETE FROM courses WHERE course_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除课程提交记录失败: %v", err)
		}
		coursesDeleted, _ := result.RowsAffected()
		fmt.Printf("  已删除 %d 条课程提交记录\n", coursesDeleted)
	} else {
		fmt.Printf("  没有找到课程提交记录\n")
	}

	// 3. 删除所有项目提交记录
	fmt.Println("  正在删除项目提交记录...")
	
	// 先查询所有有作者的项目ID
	projectRows, err := tx.Query("SELECT DISTINCT project_id FROM project_authors")
	if err != nil {
		return fmt.Errorf("查询项目ID失败: %v", err)
	}
	defer projectRows.Close()

	var projectIDs []int
	for projectRows.Next() {
		var projectID int
		if err := projectRows.Scan(&projectID); err != nil {
			continue
		}
		projectIDs = append(projectIDs, projectID)
	}

	if len(projectIDs) > 0 {
		// 构建IN子句
		placeholders := ""
		args := make([]interface{}, len(projectIDs))
		for i, id := range projectIDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args[i] = id
		}

		// 删除项目相关的关联数据
		_, err = tx.Exec(fmt.Sprintf("DELETE FROM collections WHERE resource_type = 'project' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除项目收藏记录失败: %v", err)
		}

		_, err = tx.Exec(fmt.Sprintf("DELETE FROM likes WHERE resource_type = 'project' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除项目点赞记录失败: %v", err)
		}

		_, err = tx.Exec(fmt.Sprintf("DELETE FROM comments WHERE resource_type = 'project' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除项目评论失败: %v", err)
		}

		_, err = tx.Exec(fmt.Sprintf("DELETE FROM resource_status_logs WHERE resource_type = 'project' AND resource_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除项目状态日志失败: %v", err)
		}

		// 删除项目（会级联删除 project_tech_stack, project_images, project_authors）
		result, err := tx.Exec(fmt.Sprintf("DELETE FROM projects WHERE project_id IN (%s)", placeholders), args...)
		if err != nil {
			return fmt.Errorf("删除项目提交记录失败: %v", err)
		}
		projectsDeleted, _ := result.RowsAffected()
		fmt.Printf("  已删除 %d 条项目提交记录\n", projectsDeleted)
	} else {
		fmt.Printf("  没有找到项目提交记录\n")
	}

	return nil
}

