package repository

import (
	"context"
	"database/sql"
	"fmt"
	"softeng-platform/internal/utils"
	"strings"
	"time"
)

type ToolRepository interface {
	GetTools(ctx context.Context, category, tags []string, sort, cursor string, pageSize int) ([]map[string]interface{}, error)
	GetByID(ctx context.Context, resourceID string) (map[string]interface{}, error)
	Search(ctx context.Context, keyword, cursor string, pageSize int) ([]map[string]interface{}, error)
	Create(ctx context.Context, userID int, data map[string]interface{}) (map[string]interface{}, error)
	AddLike(ctx context.Context, userID int, resourceID string) error
	RemoveLike(ctx context.Context, userID int, resourceID string) error
	GetLikes(ctx context.Context, resourceID string) (int, error)
	CollectTool(ctx context.Context, userID int, resourceID string) error
	UncollectTool(ctx context.Context, userID int, resourceID string) error
	GetCollectionCount(ctx context.Context, resourceID string) (int, error) // 从 collections 表实时统计收藏数
	GetComments(ctx context.Context, resourceID string, cursor, limit int) ([]map[string]interface{}, error)
	AddComment(ctx context.Context, userID int, resourceID, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, resourceID, commentID string) error
	ReplyComment(ctx context.Context, userID int, resourceID, commentID, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, resourceID, commentID string) error
	LikeComment(ctx context.Context, userID int, commentID string) error
	AddView(ctx context.Context, resourceID string) error
	GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error)
	UpdateToolStatus(ctx context.Context, resourceID, status, rejectReason string) error
	CheckUserLike(ctx context.Context, userID int, resourceID string) (bool, error)
	CheckUserCollect(ctx context.Context, userID int, resourceID string) (bool, error)
	GetCommentLikes(ctx context.Context, commentID string) (int, error)
	CheckUserCommentLike(ctx context.Context, userID int, commentID string) (bool, error)
	DeleteTool(ctx context.Context, resourceID string) error
	UpdateTool(ctx context.Context, resourceID string, userID int, data map[string]interface{}) error
}

type toolRepository struct {
	db *Database
}

func NewToolRepository(db *Database) ToolRepository {
	return &toolRepository{db: db}
}

func (r *toolRepository) GetTools(ctx context.Context, category, tags []string, sort, cursor string, pageSize int) ([]map[string]interface{}, error) {
	// 构建基础查询
	// 使用子查询从 collections 表实时统计收藏数，而不是使用 tools.collections 字段
	query := `
		SELECT t.resource_id, t.resource_type, t.resource_name, t.description, 
		       t.category, t.views, 
		       COALESCE((SELECT COUNT(*) FROM collections WHERE resource_type = 'tool' AND resource_id = t.resource_id), 0) as collections,
		       t.loves, t.created_at, t.status,
		       u.username as submitter_name
		FROM tools t
		LEFT JOIN users u ON t.submitter_id = u.id
		WHERE t.status = 'approved'
	`
	
	args := []interface{}{}
	
	// 添加分类过滤
	if len(category) > 0 {
		placeholders := strings.Repeat("?,", len(category))
		placeholders = placeholders[:len(placeholders)-1]
		query += fmt.Sprintf(" AND t.category IN (%s)", placeholders)
		for _, cat := range category {
			args = append(args, cat)
		}
	}
	
	// 添加标签过滤
	if len(tags) > 0 {
		query += ` AND EXISTS (
			SELECT 1 FROM tool_tags tt 
			WHERE tt.tool_id = t.resource_id 
			AND tt.tag IN (` + strings.Repeat("?,", len(tags))[:len(strings.Repeat("?,", len(tags)))-1] + `)
		)`
		for _, tag := range tags {
			args = append(args, tag)
		}
	}
	
	// 排序
	// 注意：collections 在 SELECT 中是子查询的别名，在 ORDER BY 中可以直接使用别名
	switch sort {
	case "最新", "newest", "created_at":
		query += " ORDER BY t.created_at DESC"
	case "最多浏览", "views":
		query += " ORDER BY t.views DESC"
	case "最多收藏", "collections":
		// 使用子查询统计的收藏数进行排序，而不是 tools.collections 字段
		query += " ORDER BY (SELECT COUNT(*) FROM collections WHERE resource_type = 'tool' AND resource_id = t.resource_id) DESC"
	case "最多点赞", "loves":
		query += " ORDER BY t.loves DESC"
	default:
		query += " ORDER BY t.created_at DESC"
	}
	
	// 分页
	if cursor != "" {
		// 简单实现：使用OFFSET（生产环境建议使用cursor-based pagination）
		if cursorInt, err := time.Parse("2006-01-02 15:04:05", cursor); err == nil {
			query += " AND t.created_at < ?"
			args = append(args, cursorInt)
		}
	}
	
	query += " LIMIT ?"
	args = append(args, pageSize)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tools: %w", err)
	}
	defer rows.Close()
	
	var tools []map[string]interface{}
	for rows.Next() {
		var toolID int
		var resourceType, resourceName, description, category, submitterName sql.NullString
		var views, collections, loves int
		var createdAt time.Time
		var status string
		
		if err := rows.Scan(&toolID, &resourceType, &resourceName, &description, &category, 
			&views, &collections, &loves, &createdAt, &status, &submitterName); err != nil {
			return nil, fmt.Errorf("failed to scan tool: %w", err)
		}
		
		// 获取标签
		tagRows, err := r.db.QueryContext(ctx, "SELECT tag FROM tool_tags WHERE tool_id = ?", toolID)
		var tagList []string
		if err == nil {
			defer tagRows.Close()
			for tagRows.Next() {
				var tag string
				if err := tagRows.Scan(&tag); err == nil {
					tagList = append(tagList, tag)
				}
			}
		}
		
		// 获取图片
		imageRows, err := r.db.QueryContext(ctx, "SELECT image_url FROM tool_images WHERE tool_id = ? ORDER BY sort_order", toolID)
		var imageList []string
		if err == nil {
			defer imageRows.Close()
			for imageRows.Next() {
				var imgURL string
				if err := imageRows.Scan(&imgURL); err == nil {
					imageList = append(imageList, imgURL)
				}
			}
		}
		
		// 获取贡献者
		contribRows, err := r.db.QueryContext(ctx, 
			"SELECT u.username FROM tool_contributors tc JOIN users u ON tc.user_id = u.id WHERE tc.tool_id = ?", toolID)
		var contribList []string
		if err == nil {
			defer contribRows.Close()
			for contribRows.Next() {
				var username string
				if err := contribRows.Scan(&username); err == nil {
					contribList = append(contribList, username)
				}
			}
		}
		
		tool := map[string]interface{}{
			"resourceId":   toolID,
			"resourceType": resourceType.String,
			"resourceName": resourceName.String,
			"description":  description.String,
			"image":        imageList,
			"catagory":     category.String,
			"tags":         tagList,
			"views":        views,
			"collections":  collections,
			"loves":        loves,
			"contributors": contribList,
			"createdat":    createdAt.Format("2006-01-02"),
		}
		
		tools = append(tools, tool)
	}
	
	return tools, nil
}

func (r *toolRepository) GetByID(ctx context.Context, resourceID string) (map[string]interface{}, error) {
	var toolID int
	var resourceType, resourceName, resourceLink, description, descriptionDetail, category sql.NullString
	var views, collections, loves int
	var createdAt time.Time
	var status string
	
	// 先检查 tool_type 字段是否存在
	var toolTypeExists bool
	var checkQuery string
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) > 0 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = 'tools' 
		AND COLUMN_NAME = 'tool_type'
	`).Scan(&toolTypeExists)
	
	if err != nil {
		// 如果检查失败，假设字段不存在
		toolTypeExists = false
	}
	
	var toolType sql.NullString
	if toolTypeExists {
		checkQuery = `
			SELECT resource_id, resource_type, resource_name, resource_link, 
			       description, description_detail, category, views, collections, loves, created_at, status,
			       COALESCE(tool_type, 'external') as tool_type
			FROM tools
			WHERE resource_id = ?
		`
		fmt.Printf("[DEBUG] GetByID: querying tool with resourceID=%s (with tool_type)\n", resourceID)
		err = r.db.QueryRowContext(ctx, checkQuery, resourceID).Scan(
			&toolID, &resourceType, &resourceName, &resourceLink, &description,
			&descriptionDetail, &category, &views, &collections, &loves, &createdAt, &status, &toolType,
		)
	} else {
		checkQuery = `
			SELECT resource_id, resource_type, resource_name, resource_link, 
			       description, description_detail, category, views, collections, loves, created_at, status
			FROM tools
			WHERE resource_id = ?
		`
		fmt.Printf("[DEBUG] GetByID: querying tool with resourceID=%s (without tool_type)\n", resourceID)
		err = r.db.QueryRowContext(ctx, checkQuery, resourceID).Scan(
			&toolID, &resourceType, &resourceName, &resourceLink, &description,
			&descriptionDetail, &category, &views, &collections, &loves, &createdAt, &status,
		)
		// 如果字段不存在，默认设置为 'external'
		toolType = sql.NullString{String: "external", Valid: true}
	}
	if err != nil {
		fmt.Printf("[DEBUG] GetByID error: resourceID=%s, error=%v\n", resourceID, err)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tool not found")
		}
		return nil, fmt.Errorf("failed to get tool: %w", err)
	}
	fmt.Printf("[DEBUG] GetByID success: toolID=%d, name=%s, status=%s\n", toolID, resourceName.String, status)
	
	// 获取标签
	tagRows, err := r.db.QueryContext(ctx, "SELECT tag FROM tool_tags WHERE tool_id = ?", toolID)
	var tagList []string
	if err == nil {
		defer tagRows.Close()
		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err == nil {
				tagList = append(tagList, tag)
			}
		}
	}
	
	// 获取图片
	imageRows, err := r.db.QueryContext(ctx, "SELECT image_url FROM tool_images WHERE tool_id = ? ORDER BY sort_order", toolID)
	var imageList []string
	if err == nil {
		defer imageRows.Close()
		for imageRows.Next() {
			var imgURL string
			if err := imageRows.Scan(&imgURL); err == nil {
				imageList = append(imageList, imgURL)
			}
		}
	}
	
	// 获取评论数
	var commentCount int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM comments WHERE resource_type = 'tool' AND resource_id = ? AND deleted_at IS NULL", toolID).Scan(&commentCount)
	
	// 从 collections 表实时统计收藏数，而不是使用 tools.collections 字段
	// 这样可以确保数据的准确性，避免冗余字段与实际数据不一致的问题
	actualCollections, err := r.GetCollectionCount(ctx, resourceID)
	if err != nil {
		// 如果获取失败，使用 tools 表中的值作为备用
		actualCollections = collections
	}
	
	// 处理工具类型，如果为空则默认为 'external'
	toolTypeValue := toolType.String
	if toolTypeValue == "" {
		toolTypeValue = "external"
	}
	
	return map[string]interface{}{
		"resourceId":         toolID,
		"resourceType":       resourceType.String,
		"resourceName":       resourceName.String,
		"resourceLink":       resourceLink.String,
		"description":        description.String,
		"description_detail": descriptionDetail.String,
		"catagory":           category.String,
		"image":              imageList,
		"tags":               tagList,
		"type":               toolTypeValue, // 工具类型：internal/external
		"views":              views,
		"collections":        actualCollections, // 使用实时统计的收藏数
		"loves":              loves,
		"iscollected":        false, // 需要根据当前用户判断
		"isliked":            false, // 需要根据当前用户判断
		"comment_count":      commentCount,
		"comments":           []map[string]interface{}{}, // 评论需要单独获取
		"createdDate":        createdAt.Format("2006-01-02"),
	}, nil
}

func (r *toolRepository) Search(ctx context.Context, keyword, cursor string, pageSize int) ([]map[string]interface{}, error) {
	query := `
		SELECT DISTINCT t.resource_id, t.resource_type, t.resource_name, t.description,
		       t.category, t.views, 
		       COALESCE((SELECT COUNT(*) FROM collections WHERE resource_type = 'tool' AND resource_id = t.resource_id), 0) as collections,
		       t.loves, t.created_at, t.status,
		       u.username as submitter_name
		FROM tools t
		LEFT JOIN users u ON t.submitter_id = u.id
		LEFT JOIN tool_tags tt ON t.resource_id = tt.tool_id
		WHERE t.status = 'approved'
		AND (
			t.resource_name LIKE ? OR 
			t.description LIKE ? OR 
			t.description_detail LIKE ? OR
			tt.tag LIKE ?
		)
	`
	
	searchPattern := "%" + keyword + "%"
	args := []interface{}{searchPattern, searchPattern, searchPattern, searchPattern}
	
	query += " ORDER BY t.created_at DESC LIMIT ?"
	args = append(args, pageSize)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search tools: %w", err)
	}
	defer rows.Close()
	
	var tools []map[string]interface{}
	for rows.Next() {
		var toolID int
		var resourceType, resourceName, description, category, submitterName sql.NullString
		var views, collections, loves int
		var createdAt time.Time
		var status string
		
		if err := rows.Scan(&toolID, &resourceType, &resourceName, &description, &category,
			&views, &collections, &loves, &createdAt, &status, &submitterName); err != nil {
			return nil, fmt.Errorf("failed to scan tool: %w", err)
		}
		
		// 获取标签
		tagRows, err := r.db.QueryContext(ctx, "SELECT tag FROM tool_tags WHERE tool_id = ?", toolID)
		var tagList []string
		if err == nil {
			defer tagRows.Close()
			for tagRows.Next() {
				var tag string
				if err := tagRows.Scan(&tag); err == nil {
					tagList = append(tagList, tag)
				}
			}
		}
		
		// 获取图片（第一张）
		var imageURL sql.NullString
		r.db.QueryRowContext(ctx, "SELECT image_url FROM tool_images WHERE tool_id = ? ORDER BY sort_order LIMIT 1", toolID).Scan(&imageURL)
		
		// 获取贡献者
		contribRows, err := r.db.QueryContext(ctx,
			"SELECT u.username FROM tool_contributors tc JOIN users u ON tc.user_id = u.id WHERE tc.tool_id = ?", toolID)
		var contribList []string
		if err == nil {
			defer contribRows.Close()
			for contribRows.Next() {
				var username string
				if err := contribRows.Scan(&username); err == nil {
					contribList = append(contribList, username)
				}
			}
		}
		
		imageList := []string{}
		if imageURL.Valid {
			imageList = append(imageList, imageURL.String)
		}
		
		tool := map[string]interface{}{
			"resourceId":   toolID,
			"resourceType": resourceType.String,
			"resourceName": resourceName.String,
			"description":  description.String,
			"image":        imageList,
			"catagory":     category.String,
			"tags":         tagList,
			"views":        views,
			"collections":  collections,
			"loves":        loves,
			"contributors": contribList,
			"createdat":    createdAt.Format("2006-01-02"),
		}
		
		tools = append(tools, tool)
	}
	
	return tools, nil
}

func (r *toolRepository) Create(ctx context.Context, userID int, data map[string]interface{}) (map[string]interface{}, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// 检查 tool_type 字段是否存在
	var toolTypeExists bool
	r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) > 0 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = 'tools' 
		AND COLUMN_NAME = 'tool_type'
	`).Scan(&toolTypeExists)
	
	var query string
	var result sql.Result
	if toolTypeExists {
		query = `
			INSERT INTO tools (resource_type, resource_name, resource_link, description, 
			                  description_detail, category, tool_type, submitter_id, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'pending', NOW(), NOW())
		`
		toolType := getString(data, "tool_type")
		if toolType == "" {
			toolType = "external" // 默认为外部工具
		}
		result, err = tx.ExecContext(ctx, query,
			"tool",
			getString(data, "name"),
			getString(data, "link"),
			getString(data, "description"),
			getString(data, "description_detail"),
			getString(data, "category"),
			toolType,
			userID,
		)
	} else {
		query = `
			INSERT INTO tools (resource_type, resource_name, resource_link, description, 
			                  description_detail, category, submitter_id, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, 'pending', NOW(), NOW())
		`
		result, err = tx.ExecContext(ctx, query,
			"tool",
			getString(data, "name"),
			getString(data, "link"),
			getString(data, "description"),
			getString(data, "description_detail"),
			getString(data, "category"),
			userID,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to insert tool: %w", err)
	}
	
	toolID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get tool ID: %w", err)
	}
	
	// 插入标签
	if tags, ok := data["tags"].([]string); ok {
		tagStmt, _ := tx.PrepareContext(ctx, "INSERT INTO tool_tags (tool_id, tag) VALUES (?, ?)")
		for _, tag := range tags {
			tagStmt.ExecContext(ctx, toolID, tag)
		}
		tagStmt.Close()
	}
	
	// 插入图片
	if images, ok := data["images"].([]string); ok {
		imgStmt, _ := tx.PrepareContext(ctx, "INSERT INTO tool_images (tool_id, image_url, sort_order) VALUES (?, ?, ?)")
		for i, img := range images {
			imgStmt.ExecContext(ctx, toolID, img, i)
		}
		imgStmt.Close()
	}
	
	// 添加提交者为贡献者
	tx.ExecContext(ctx, "INSERT INTO tool_contributors (tool_id, user_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE tool_id=tool_id", toolID, userID)
	
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	// 获取用户信息
	var username sql.NullString
	r.db.QueryRowContext(ctx, "SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	
	return map[string]interface{}{
		"resourceId":   toolID,
		"resourceType": "tool",
		"resource":     getString(data, "link"),
		"auditStatus":  "pending",
		"submitTime":   time.Now().Format("2006-01-02 15:04:05"),
		"auditTime":    nil,
		"rejectReason": nil,
		"submitor":     username.String,
	}, nil
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (r *toolRepository) AddLike(ctx context.Context, userID int, resourceID string) error {
	// 先检查是否已点赞
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM likes WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?", userID, resourceID).Scan(&exists)
	
	if exists > 0 {
		// 已点赞，直接返回（不做任何操作）
		return nil
	}
	
	// 未点赞，执行插入
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO likes (user_id, resource_type, resource_id, created_at) VALUES (?, 'tool', ?, NOW())",
		userID, resourceID)
	if err != nil {
		return fmt.Errorf("failed to add like: %w", err)
	}
	
	// 增加点赞数
	_, err = r.db.ExecContext(ctx,
		"UPDATE tools SET loves = loves + 1 WHERE resource_id = ?",
		resourceID)
	return err
}

func (r *toolRepository) GetLikes(ctx context.Context, resourceID string) (int, error) {
	var likes int
	err := r.db.QueryRowContext(ctx,
		"SELECT loves FROM tools WHERE resource_id = ?",
		resourceID).Scan(&likes)
	return likes, err
}

func (r *toolRepository) GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT t.resource_id, t.resource_name, t.resource_link, t.description, 
		       t.category, t.created_at, COALESCE(u.nickname, u.username) as nickname, COALESCE(u.avatar, '') as avatar, u.username
		FROM tools t
		LEFT JOIN users u ON t.submitter_id = u.id
		WHERE t.status = 'pending'
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending tools: %w", err)
	}
	defer rows.Close()
	
	var tools []map[string]interface{}
	for rows.Next() {
		var toolID int
		var resourceName, resourceLink, description, category, nickname, avatar, username sql.NullString
		var createdAt time.Time
		
		if err := rows.Scan(&toolID, &resourceName, &resourceLink, &description, &category, &createdAt, &nickname, &avatar, &username); err != nil {
			return nil, fmt.Errorf("failed to scan tool: %w", err)
		}
		
		// 获取标签
		tagRows, err := r.db.QueryContext(ctx, "SELECT tag FROM tool_tags WHERE tool_id = ?", toolID)
		var tagList []string
		if err == nil {
			defer tagRows.Close()
			for tagRows.Next() {
				var tag string
				if err := tagRows.Scan(&tag); err == nil {
					tagList = append(tagList, tag)
				}
			}
		}
		
		tool := map[string]interface{}{
			"id":           toolID,
			"resourceId":   toolID,
			"reourceId":    toolID, // 保持兼容性
			"name":         resourceName.String,
			"resourcename": resourceName.String, // 保持兼容性
			"title":        resourceName.String, // 前端可能使用title
			"url":          resourceLink.String,
			"link":         resourceLink.String, // 保持兼容性
			"desc":         description.String,
			"description":  description.String, // 保持兼容性
			"category":     category.String,
			"catagory":     category.String, // 保持兼容性（拼写错误）
			"tags":         tagList,
			"uploader":     nickname.String, // 使用昵称而不是用户名
			"uploaderNickname": nickname.String, // 明确标识为昵称
			"uploaderAvatar": avatar.String, // 提交者头像
			"submitor":     nickname.String, // 保持兼容性
			"author":       nickname.String, // 前端可能使用author
			"owner":        nickname.String, // 前端可能使用owner
			"username":     username.String, // 保留用户名字段用于备用
			"created_at":   createdAt.Format("2006-01-02T15:04:05Z"),
			"created":      createdAt.Format("2006-01-02T15:04:05Z"), // 前端可能使用created
			"submitDate":   createdAt.Format("2006-01-02 15:04:05"), // 保持兼容性
			"status":       "pending",
		}
		
		tools = append(tools, tool)
	}
	
	return tools, nil
}

// RemoveLike 取消点赞
func (r *toolRepository) RemoveLike(ctx context.Context, userID int, resourceID string) error {
	// 先检查是否已点赞，只有已点赞才执行删除和减少计数
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM likes WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?", userID, resourceID).Scan(&exists)
	
	if exists == 0 {
		// 未点赞，直接返回
		return nil
	}
	
	// 已点赞，执行删除
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM likes WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?",
		userID, resourceID)
	if err != nil {
		return fmt.Errorf("failed to remove like: %w", err)
	}
	
	// 检查删除是否成功（受影响的行数）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	// 只有成功删除才减少计数
	if rowsAffected > 0 {
		_, err = r.db.ExecContext(ctx,
			"UPDATE tools SET loves = GREATEST(loves - 1, 0) WHERE resource_id = ?",
			resourceID)
	}
	return err
}

// CollectTool 收藏工具
func (r *toolRepository) CollectTool(ctx context.Context, userID int, resourceID string) error {
	// 先检查是否已收藏
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM collections WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?", userID, resourceID).Scan(&exists)
	
	if exists > 0 {
		// 已收藏，直接返回（不做任何操作）
		return nil
	}
	
	// 未收藏，执行插入
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO collections (user_id, resource_type, resource_id, created_at) VALUES (?, 'tool', ?, NOW())",
		userID, resourceID)
	if err != nil {
		return fmt.Errorf("failed to collect tool: %w", err)
	}
	
	// 增加收藏数
	_, err = r.db.ExecContext(ctx,
		"UPDATE tools SET collections = collections + 1 WHERE resource_id = ?",
		resourceID)
	return err
}

// UncollectTool 取消收藏
func (r *toolRepository) UncollectTool(ctx context.Context, userID int, resourceID string) error {
	// 先检查是否已收藏，只有已收藏才执行删除和减少计数
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM collections WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?", userID, resourceID).Scan(&exists)
	
	if exists == 0 {
		// 未收藏，直接返回
		return nil
	}
	
	// 已收藏，执行删除
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM collections WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?",
		userID, resourceID)
	if err != nil {
		return fmt.Errorf("failed to uncollect tool: %w", err)
	}
	
	// 检查删除是否成功（受影响的行数）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	// 只有成功删除才减少计数
	if rowsAffected > 0 {
		_, err = r.db.ExecContext(ctx,
			"UPDATE tools SET collections = GREATEST(collections - 1, 0) WHERE resource_id = ?",
			resourceID)
	}
	return err
}

// GetCollectionCount 从 collections 表实时统计收藏数
func (r *toolRepository) GetCollectionCount(ctx context.Context, resourceID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, 
		"SELECT COUNT(*) FROM collections WHERE resource_type = 'tool' AND resource_id = ?", 
		resourceID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get collection count: %w", err)
	}
	return count, nil
}

// GetComments 获取评论列表
func (r *toolRepository) GetComments(ctx context.Context, resourceID string, cursor, limit int) ([]map[string]interface{}, error) {
	// 修复：使用 COALESCE 确保如果昵称为空，使用用户名作为备用，但字段名仍然是 nickname
	query := `SELECT c.comment_id, c.user_id, c.content, c.love_count, c.reply_total, c.created_at, 
	          COALESCE(u.nickname, u.username) as nickname, COALESCE(u.avatar, '') as avatar 
	          FROM comments c 
	          JOIN users u ON c.user_id = u.id 
	          WHERE c.resource_type = 'tool' AND c.resource_id = ? AND c.parent_id IS NULL AND c.deleted_at IS NULL 
	          ORDER BY c.created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, resourceID, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()
	var comments []map[string]interface{}
	for rows.Next() {
		var commentID, userID, loveCount, replyTotal int
		var content, nickname, avatar sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&commentID, &userID, &content, &loveCount, &replyTotal, &createdAt, &nickname, &avatar); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		// 修复：回复评论也使用 COALESCE 获取昵称
		replyRows, _ := r.db.QueryContext(ctx, `SELECT c.comment_id, c.user_id, c.content, c.love_count, c.created_at, 
		                                          COALESCE(u.nickname, u.username) as nickname, COALESCE(u.avatar, '') as avatar, c.parent_id 
		                                          FROM comments c 
		                                          JOIN users u ON c.user_id = u.id 
		                                          WHERE c.parent_id = ? AND c.deleted_at IS NULL 
		                                          ORDER BY c.created_at ASC`, commentID)
		var replies []map[string]interface{}
		if replyRows != nil {
			defer replyRows.Close()
			for replyRows.Next() {
				var replyID, replyUserID, replyLoveCount int
				var replyContent, replyNickname, replyAvatar sql.NullString
				var replyCreatedAt time.Time
				var parentID sql.NullInt64
				if err := replyRows.Scan(&replyID, &replyUserID, &replyContent, &replyLoveCount, &replyCreatedAt, &replyNickname, &replyAvatar, &parentID); err == nil {
					replies = append(replies, map[string]interface{}{
						"comment_Id": replyID, "nickname": replyNickname.String, "avatar": replyAvatar.String, "avater": replyAvatar.String, // 兼容拼写错误
						"comment": replyContent.String, "commentDate": replyCreatedAt.Format("2006-01-02 15:04:05"),
						"love_count": replyLoveCount, "isreply": true, "reply_id": parentID.Int64,
					})
				}
			}
		}
		comments = append(comments, map[string]interface{}{
			"comment_Id": commentID, "nickname": nickname.String, "avatar": avatar.String, "avater": avatar.String, // 兼容拼写错误
			"comment": content.String, "commentDate": createdAt.Format("2006-01-02 15:04:05"),
			"love_count": loveCount, "reply_total": replyTotal, "replies": replies,
		})
	}
	return comments, nil
}

// AddComment 添加评论
func (r *toolRepository) AddComment(ctx context.Context, userID int, resourceID, content string) (map[string]interface{}, error) {
	// 修复：使用 COALESCE 确保如果昵称为空，使用用户名作为备用，但字段名仍然是 nickname
	var nickname, username, avatar sql.NullString
	err := r.db.QueryRowContext(ctx, "SELECT COALESCE(nickname, username) as nickname, username, COALESCE(avatar, '') as avatar FROM users WHERE id = ?", userID).Scan(&nickname, &username, &avatar)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	result, err := r.db.ExecContext(ctx, "INSERT INTO comments (resource_type, resource_id, user_id, content, created_at, updated_at) VALUES ('tool', ?, ?, ?, NOW(), NOW())", resourceID, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get comment ID: %w", err)
	}
	now := time.Now()
	return map[string]interface{}{
		"id":         commentID,
		"commentId":  commentID,
		"userId":     userID,
		"nickname":   nickname.String, // 这里已经是 COALESCE(nickname, username) 的结果，确保有值
		"username":   username.String, // 保留 username 字段用于后端逻辑，但前端应该只使用 nickname
		"avatar":     avatar.String,
		"avater":     avatar.String, // 兼容拼写错误
		"content":    content,
		"comment":    content, // 兼容字段
		"createdAt":  now.Format("2006-01-02 15:04:05"),
		"commentDate": now.Format("2006-01-02 15:04:05"), // 兼容字段
		"likes":      0,
		"love_count": 0, // 兼容字段
		"isLiked":    false,
		"reply_total": 0,
		"replies":    []interface{}{},
	}, nil
}

// DeleteComment 删除评论（软删除）
func (r *toolRepository) DeleteComment(ctx context.Context, userID int, resourceID, commentID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE comments SET deleted_at = NOW() WHERE comment_id = ? AND user_id = ? AND resource_type = 'tool' AND resource_id = ?", commentID, userID, resourceID)
	return err
}

// ReplyComment 回复评论
func (r *toolRepository) ReplyComment(ctx context.Context, userID int, resourceID, commentID, content string) (map[string]interface{}, error) {
	// 修复：使用 COALESCE 确保如果昵称为空，使用用户名作为备用，但字段名仍然是 nickname
	var nickname, username, avatar sql.NullString
	err := r.db.QueryRowContext(ctx, "SELECT COALESCE(nickname, username) as nickname, username, COALESCE(avatar, '') as avatar FROM users WHERE id = ?", userID).Scan(&nickname, &username, &avatar)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	result, err := r.db.ExecContext(ctx, "INSERT INTO comments (resource_type, resource_id, parent_id, user_id, content, created_at, updated_at) VALUES ('tool', ?, ?, ?, ?, NOW(), NOW())", resourceID, commentID, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to add reply: %w", err)
	}
	replyID, _ := result.LastInsertId()
	r.db.ExecContext(ctx, "UPDATE comments SET reply_total = reply_total + 1 WHERE comment_id = ?", commentID)
	return map[string]interface{}{
		"comment_Id": replyID, "nickname": nickname.String, "avatar": avatar.String, "avater": avatar.String, // 兼容拼写错误
		"comment": content, "commentDate": time.Now().Format("2006-01-02 15:04:05"),
		"love_count": 0, "isreply": true, "reply_id": commentID,
	}, nil
}

// DeleteReply 删除回复
func (r *toolRepository) DeleteReply(ctx context.Context, userID int, resourceID, commentID string) error {
	var parentID sql.NullInt64
	r.db.QueryRowContext(ctx, "SELECT parent_id FROM comments WHERE comment_id = ?", commentID).Scan(&parentID)
	_, err := r.db.ExecContext(ctx, "UPDATE comments SET deleted_at = NOW() WHERE comment_id = ? AND user_id = ? AND resource_type = 'tool' AND resource_id = ?", commentID, userID, resourceID)
	if err == nil && parentID.Valid {
		r.db.ExecContext(ctx, "UPDATE comments SET reply_total = GREATEST(reply_total - 1, 0) WHERE comment_id = ?", parentID.Int64)
	}
	return err
}

// LikeComment 点赞评论
func (r *toolRepository) LikeComment(ctx context.Context, userID int, commentID string) error {
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&exists)
	if exists > 0 {
		_, err := r.db.ExecContext(ctx, "DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)
		if err == nil {
			r.db.ExecContext(ctx, "UPDATE comments SET love_count = GREATEST(love_count - 1, 0) WHERE comment_id = ?", commentID)
		}
		return err
	}
	_, err := r.db.ExecContext(ctx, "INSERT INTO comment_likes (comment_id, user_id, created_at) VALUES (?, ?, NOW())", commentID, userID)
	if err == nil {
		r.db.ExecContext(ctx, "UPDATE comments SET love_count = love_count + 1 WHERE comment_id = ?", commentID)
	}
	return err
}

// AddView 增加浏览量
func (r *toolRepository) AddView(ctx context.Context, resourceID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tools SET views = views + 1 WHERE resource_id = ?", resourceID)
	return err
}

// UpdateToolStatus 更新工具审核状态
// 注意：驳回时不再删除工具记录，而是保留记录以便用户查看审核状态
func (r *toolRepository) UpdateToolStatus(ctx context.Context, resourceID, status, rejectReason string) error {
	// 注意：不再删除被驳回的工具，而是保留记录以便用户查看审核状态
	// 只有用户主动撤回时才删除（如果需要的话，可以在用户撤回接口中调用 DeleteTool）
	
	// 否则，只更新状态
	query := "UPDATE tools SET status = ?, audit_time = NOW()"
	args := []interface{}{status}
	
	if rejectReason != "" {
		query += ", reject_reason = ?"
		args = append(args, rejectReason)
	}
	
	query += " WHERE resource_id = ?"
	args = append(args, resourceID)
	
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// CheckUserLike 检查用户是否点赞
func (r *toolRepository) CheckUserLike(ctx context.Context, userID int, resourceID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM likes WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?", userID, resourceID).Scan(&count)
	return count > 0, err
}

// CheckUserCollect 检查用户是否收藏
func (r *toolRepository) CheckUserCollect(ctx context.Context, userID int, resourceID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM collections WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?", userID, resourceID).Scan(&count)
	return count > 0, err
}

// GetCommentLikes 获取评论点赞数
func (r *toolRepository) GetCommentLikes(ctx context.Context, commentID string) (int, error) {
	var likes int
	err := r.db.QueryRowContext(ctx, "SELECT love_count FROM comments WHERE comment_id = ?", commentID).Scan(&likes)
	return likes, err
}

// CheckUserCommentLike 检查用户是否点赞了评论
func (r *toolRepository) CheckUserCommentLike(ctx context.Context, userID int, commentID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&count)
	return count > 0, err
}

// DeleteTool 删除工具（包括数据库记录和图片文件）
func (r *toolRepository) DeleteTool(ctx context.Context, resourceID string) error {
	// 先获取所有图片URL，以便删除文件
	var imageURLs []string
	rows, err := r.db.QueryContext(ctx, "SELECT image_url FROM tool_images WHERE tool_id = ?", resourceID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var imgURL sql.NullString
			if rows.Scan(&imgURL) == nil && imgURL.Valid {
				imageURLs = append(imageURLs, imgURL.String)
			}
		}
	}
	
	// 删除工具记录（CASCADE会自动删除tool_images记录）
	_, err = r.db.ExecContext(ctx, "DELETE FROM tools WHERE resource_id = ?", resourceID)
	if err != nil {
		return fmt.Errorf("failed to delete tool: %w", err)
	}
	
	// 删除图片文件
	for _, imgURL := range imageURLs {
		if err := utils.DeleteImageFile(imgURL); err != nil {
			// 记录错误但不中断流程（文件可能已被删除）
			fmt.Printf("Warning: failed to delete image file %s: %v\n", imgURL, err)
		}
	}
	
	return nil
}

// UpdateTool 更新工具信息
func (r *toolRepository) UpdateTool(ctx context.Context, resourceID string, userID int, data map[string]interface{}) error {
	// 验证用户是否有权限更新该工具
	var submitterID int
	err := r.db.QueryRowContext(ctx, "SELECT submitter_id FROM tools WHERE resource_id = ?", resourceID).Scan(&submitterID)
	if err != nil {
		return fmt.Errorf("tool not found: %w", err)
	}
	if submitterID != userID {
		return fmt.Errorf("unauthorized: you can only update your own tools")
	}
	
	// 验证工具状态是否为pending（只有pending状态的工具可以编辑）
	var status string
	err = r.db.QueryRowContext(ctx, "SELECT status FROM tools WHERE resource_id = ?", resourceID).Scan(&status)
	if err != nil {
		return fmt.Errorf("failed to get tool status: %w", err)
	}
	if status != "pending" {
		return fmt.Errorf("only pending tools can be updated")
	}
	
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// 检查 tool_type 字段是否存在
	var toolTypeExists bool
	r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) > 0 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = 'tools' 
		AND COLUMN_NAME = 'tool_type'
	`).Scan(&toolTypeExists)
	
	// 更新工具基本信息
	var query string
	if toolTypeExists {
		query = `
			UPDATE tools 
			SET resource_name = ?, resource_link = ?, description = ?, 
			    description_detail = ?, category = ?, tool_type = ?, updated_at = NOW()
			WHERE resource_id = ?
		`
		toolType := getString(data, "tool_type")
		if toolType == "" {
			toolType = "external" // 默认为外部工具
		}
		_, err = tx.ExecContext(ctx, query,
			getString(data, "name"),
			getString(data, "link"),
			getString(data, "description"),
			getString(data, "description_detail"),
			getString(data, "category"),
			toolType,
			resourceID,
		)
	} else {
		query = `
			UPDATE tools 
			SET resource_name = ?, resource_link = ?, description = ?, 
			    description_detail = ?, category = ?, updated_at = NOW()
			WHERE resource_id = ?
		`
		_, err = tx.ExecContext(ctx, query,
			getString(data, "name"),
			getString(data, "link"),
			getString(data, "description"),
			getString(data, "description_detail"),
			getString(data, "category"),
			resourceID,
		)
	}
	if err != nil {
		return fmt.Errorf("failed to update tool: %w", err)
	}
	
	// 删除旧标签
	_, err = tx.ExecContext(ctx, "DELETE FROM tool_tags WHERE tool_id = ?", resourceID)
	if err != nil {
		return fmt.Errorf("failed to delete old tags: %w", err)
	}
	
	// 插入新标签
	if tags, ok := data["tags"].([]string); ok && len(tags) > 0 {
		tagStmt, _ := tx.PrepareContext(ctx, "INSERT INTO tool_tags (tool_id, tag) VALUES (?, ?)")
		for _, tag := range tags {
			tagStmt.ExecContext(ctx, resourceID, tag)
		}
		tagStmt.Close()
	}
	
	// 删除旧图片
	_, err = tx.ExecContext(ctx, "DELETE FROM tool_images WHERE tool_id = ?", resourceID)
	if err != nil {
		return fmt.Errorf("failed to delete old images: %w", err)
	}
	
	// 插入新图片
	if images, ok := data["images"].([]string); ok && len(images) > 0 {
		imgStmt, _ := tx.PrepareContext(ctx, "INSERT INTO tool_images (tool_id, image_url, sort_order) VALUES (?, ?, ?)")
		for i, img := range images {
			if img != "" {
				imgStmt.ExecContext(ctx, resourceID, img, i)
			}
		}
		imgStmt.Close()
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}
