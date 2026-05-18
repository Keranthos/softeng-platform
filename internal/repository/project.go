package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type ProjectRepository interface {
	GetProjects(ctx context.Context, category string, techStack []string, sort string, limit int, cursor string) ([]map[string]interface{}, error)
	GetByID(ctx context.Context, projectID string) (map[string]interface{}, error)
	Search(ctx context.Context, keyword string, category []string, cursor string, limit int) ([]map[string]interface{}, error)
	Create(ctx context.Context, userID int, data map[string]interface{}) (map[string]interface{}, error)
	Update(ctx context.Context, userID int, projectID string, data map[string]interface{}) (map[string]interface{}, error)
	LikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
	UnlikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
	GetComments(ctx context.Context, projectID string, cursor, limit int) ([]map[string]interface{}, error)
	AddComment(ctx context.Context, userID int, projectID, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, projectID, commentID string) error
	LikeComment(ctx context.Context, userID int, commentID string) (bool, int, error)
	ReplyComment(ctx context.Context, userID int, projectID, commentID, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, projectID, commentID string) (map[string]interface{}, error)
	AddView(ctx context.Context, projectID string) (int, error)
	CollectProject(ctx context.Context, userID int, projectID string) error
	UncollectProject(ctx context.Context, userID int, projectID string) error
	GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error)
	UpdateProjectStatus(ctx context.Context, resourceID, status, rejectReason string) error
	GetCollectionCount(ctx context.Context, projectID string) (int, error) // 从 collections 表实时统计收藏数
}

type projectRepository struct {
	db *Database
}

func NewProjectRepository(db *Database) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) GetProjects(ctx context.Context, category string, techStack []string, sort string, limit int, cursor string) ([]map[string]interface{}, error) {
	// 使用子查询从 collections 表实时统计收藏数，而不是使用 projects.collections 字段
	query := `
		SELECT p.project_id, p.resource_type, p.name, p.description, p.category, p.cover, p.views, p.loves, 
		       COALESCE((SELECT COUNT(*) FROM collections WHERE resource_type = 'project' AND resource_id = p.project_id), 0) as collections,
		       p.created_at 
		FROM projects p 
		WHERE p.status = 'approved'
	`
	args := []interface{}{}
	
	if category != "" {
		query += " AND p.category = ?"
		args = append(args, category)
	}
	
	if len(techStack) > 0 {
		query += ` AND EXISTS (SELECT 1 FROM project_tech_stack pts WHERE pts.project_id = p.project_id AND pts.tech IN (` + strings.Repeat("?,", len(techStack))[:len(strings.Repeat("?,", len(techStack)))-1] + `))`
		for _, tech := range techStack {
			args = append(args, tech)
		}
	}
	
	switch sort {
	case "最新", "newest":
		query += " ORDER BY p.created_at DESC"
	case "最多浏览", "views":
		query += " ORDER BY p.views DESC"
	case "最多收藏", "collections":
		// 使用子查询统计的收藏数进行排序，而不是 projects.collections 字段
		query += " ORDER BY (SELECT COUNT(*) FROM collections WHERE resource_type = 'project' AND resource_id = p.project_id) DESC"
	case "最多点赞", "loves":
		query += " ORDER BY p.loves DESC"
	default:
		query += " ORDER BY p.created_at DESC"
	}
	
	cursorInt := 0
	if cursor != "" {
		fmt.Sscanf(cursor, "%d", &cursorInt)
	}
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, cursorInt)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()
	
	var projects []map[string]interface{}
	for rows.Next() {
		var projectID, views, loves, collections int
		var resourceType, name, description, category, cover sql.NullString
		var createdAt time.Time
		
		// collections 现在是从子查询中获取的实时统计值
		if err := rows.Scan(&projectID, &resourceType, &name, &description, &category, &cover, &views, &loves, &collections, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		
		techRows, _ := r.db.QueryContext(ctx, "SELECT tech FROM project_tech_stack WHERE project_id = ?", projectID)
		var techList []string
		if techRows != nil {
			defer techRows.Close()
			for techRows.Next() {
				var tech string
				if err := techRows.Scan(&tech); err == nil {
					techList = append(techList, tech)
				}
			}
		}
		
		authorRows, _ := r.db.QueryContext(ctx, "SELECT u.username FROM project_authors pa JOIN users u ON pa.user_id = u.id WHERE pa.project_id = ?", projectID)
		var authorList []string
		if authorRows != nil {
			defer authorRows.Close()
			for authorRows.Next() {
				var username string
				if err := authorRows.Scan(&username); err == nil {
					authorList = append(authorList, username)
				}
			}
		}
		
		// 统一使用 collections 字段，不再返回 loves 和 stars 字段
		projects = append(projects, map[string]interface{}{
			"projectId": projectID, "id": projectID, // 兼容两种字段名
			"resourceType": resourceType.String, "name": name.String,
			"description": description.String, "category": category.String,
			"techStack": techList, "technologies": techList, // 兼容两种字段名
			"authername": authorList, "contributors": authorList, // 兼容两种字段名
			"cover": cover.String, "coverImage": cover.String, "logo": cover.String, // 兼容多种字段名
			"createdat": createdAt.Format("2006-01-02"), "createdAt": createdAt.Format("2006-01-02"), // 兼容两种字段名
			"collections": collections, "views": views, // collections 是从子查询实时统计的收藏数
		})
	}
	
	return projects, nil
}

func (r *projectRepository) GetByID(ctx context.Context, projectID string) (map[string]interface{}, error) {
	var pID, views, loves, collections int
	var resourceType, name, description, githubURL, category, cover sql.NullString
	var detail sql.NullString
	var createdAt time.Time
	
	// 只允许查看已审核通过的项目（pending和rejected的项目不应该公开显示）
	err := r.db.QueryRowContext(ctx, `SELECT project_id, resource_type, name, description, detail, github_url, category, cover, views, loves, collections, created_at FROM projects WHERE project_id = ? AND status = 'approved'`, projectID).Scan(&pID, &resourceType, &name, &description, &detail, &githubURL, &category, &cover, &views, &loves, &collections, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	
	techRows, _ := r.db.QueryContext(ctx, "SELECT tech FROM project_tech_stack WHERE project_id = ?", pID)
	var techList []string
	if techRows != nil {
		defer techRows.Close()
		for techRows.Next() {
			var tech string
			if err := techRows.Scan(&tech); err == nil {
				techList = append(techList, tech)
			}
		}
	}
	
	imageRows, _ := r.db.QueryContext(ctx, "SELECT image_url FROM project_images WHERE project_id = ? ORDER BY sort_order", pID)
	var imageList []string
	if imageRows != nil {
		defer imageRows.Close()
		for imageRows.Next() {
			var imgURL string
			if err := imageRows.Scan(&imgURL); err == nil {
				imageList = append(imageList, imgURL)
			}
		}
	}
	
	// 查询项目作者/贡献者，返回完整的用户信息（对象数组）
	authorRows, _ := r.db.QueryContext(ctx, "SELECT u.id, COALESCE(u.nickname, u.username) as name, COALESCE(u.avatar, '') as avatar, u.username FROM project_authors pa JOIN users u ON pa.user_id = u.id WHERE pa.project_id = ?", pID)
	var authorList []string  // 保持向后兼容（字符串数组）
	var contributorsList []map[string]interface{}  // 新的对象数组格式
	if authorRows != nil {
		defer authorRows.Close()
		for authorRows.Next() {
			var userID int
			var name, avatar, username sql.NullString
			if err := authorRows.Scan(&userID, &name, &avatar, &username); err == nil {
				authorList = append(authorList, username.String)
				// 构建贡献者对象，前端期望格式：{ id, name, avatar, role }
				contributorsList = append(contributorsList, map[string]interface{}{
					"id":     userID,
					"name":   name.String,
					"avatar": avatar.String,
					"role":   "贡献者", // 默认角色，如果以后数据库中有角色字段可以修改
				})
			}
		}
	}
	
	var commentCount int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM comments WHERE resource_type = 'project' AND resource_id = ? AND deleted_at IS NULL", pID).Scan(&commentCount)
	
	// 从 collections 表实时统计收藏数，而不是使用 projects.collections 字段
	// 这样可以确保数据的准确性，避免冗余字段与实际数据不一致的问题
	actualCollections, err := r.GetCollectionCount(ctx, projectID)
	if err != nil {
		// 如果获取失败，使用 projects 表中的值作为备用
		actualCollections = collections
	}
	
	// 使用 cover 或 images 的第一张作为封面图
	coverImage := cover.String
	if coverImage == "" && len(imageList) > 0 {
		coverImage = imageList[0]
	}
	
	return map[string]interface{}{
		"projectId": pID, "resourceType": resourceType.String, "name": name.String,
		"description": description.String, 
		"detail": detail.String, "details": detail.String, // 兼容两种字段名
		"githubURL": githubURL.String, "githubUrl": githubURL.String, // 兼容两种字段名
		"techStack": techList, "technologies": techList, // 兼容两种字段名
		"catagory": category.String, "category": category.String, // 兼容两种字段名
		"cover": cover.String, "coverImage": coverImage, "logo": coverImage, // 兼容多种字段名
		"images": imageList,
		"views": views, "collections": actualCollections, // 使用实时统计的收藏数
		"isliked": false, "iscollected": false,
		"author": authorList, "contributors": contributorsList, // contributors 使用对象数组格式
		"comment_count": commentCount, "comments": []map[string]interface{}{},
		"createdAt": createdAt.Format("2006-01-02"),
	}, nil
}

func (r *projectRepository) Search(ctx context.Context, keyword string, category []string, cursor string, limit int) ([]map[string]interface{}, error) {
	// 使用子查询从 collections 表实时统计收藏数，而不是使用 projects.collections 字段
	query := `
		SELECT DISTINCT p.project_id, p.resource_type, p.name, p.description, p.category, p.cover, p.views, p.loves, 
		       COALESCE((SELECT COUNT(*) FROM collections WHERE resource_type = 'project' AND resource_id = p.project_id), 0) as collections,
		       p.created_at 
		FROM projects p 
		LEFT JOIN project_tech_stack pts ON p.project_id = pts.project_id 
		WHERE p.status = 'approved'
	`
	args := []interface{}{}
	
	if keyword != "" {
		query += " AND (p.name LIKE ? OR p.description LIKE ? OR pts.tech LIKE ?)"
		pattern := "%" + keyword + "%"
		args = append(args, pattern, pattern, pattern)
	}
	
	if len(category) > 0 {
		query += " AND p.category IN (" + strings.Repeat("?,", len(category))[:len(strings.Repeat("?,", len(category)))-1] + ")"
		for _, cat := range category {
			args = append(args, cat)
		}
	}
	
	query += " ORDER BY p.created_at DESC"
	cursorInt := 0
	if cursor != "" {
		fmt.Sscanf(cursor, "%d", &cursorInt)
	}
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, cursorInt)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search projects: %w", err)
	}
	defer rows.Close()
	
	var projects []map[string]interface{}
	for rows.Next() {
		var projectID, views, loves, collections int
		var resourceType, name, description, category, cover sql.NullString
		var createdAt time.Time
		
		if err := rows.Scan(&projectID, &resourceType, &name, &description, &category, &cover, &views, &loves, &collections, &createdAt); err != nil {
			continue
		}
		
		techRows, _ := r.db.QueryContext(ctx, "SELECT tech FROM project_tech_stack WHERE project_id = ?", projectID)
		var techList []string
		if techRows != nil {
			defer techRows.Close()
			for techRows.Next() {
				var tech string
				if err := techRows.Scan(&tech); err == nil {
					techList = append(techList, tech)
				}
			}
		}
		
		authorRows, _ := r.db.QueryContext(ctx, "SELECT u.username FROM project_authors pa JOIN users u ON pa.user_id = u.id WHERE pa.project_id = ?", projectID)
		var authorList []string
		if authorRows != nil {
			defer authorRows.Close()
			for authorRows.Next() {
				var username string
				if err := authorRows.Scan(&username); err == nil {
					authorList = append(authorList, username)
				}
			}
		}
		
		// 统一使用 collections 字段，不再返回 loves 和 stars 字段
		projects = append(projects, map[string]interface{}{
			"projectId": projectID, "id": projectID, // 兼容两种字段名
			"resourceType": resourceType.String, "name": name.String,
			"description": description.String, "category": category.String,
			"techStack": techList, "technologies": techList, // 兼容两种字段名
			"authername": authorList, "contributors": authorList, // 兼容两种字段名
			"cover": cover.String, "coverImage": cover.String, "logo": cover.String, // 兼容多种字段名
			"createdat": createdAt.Format("2006-01-02"), "createdAt": createdAt.Format("2006-01-02"), // 兼容两种字段名
			"collections": collections, "views": views, // collections 是从子查询实时统计的收藏数
		})
	}
	
	return projects, nil
}

func (r *projectRepository) Create(ctx context.Context, userID int, data map[string]interface{}) (map[string]interface{}, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	result, err := tx.ExecContext(ctx, `INSERT INTO projects (resource_type, name, description, detail, github_url, category, cover, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, 'pending', NOW(), NOW())`,
		"project", getString(data, "name"), getString(data, "description"), getString(data, "detail"), getString(data, "github"), getString(data, "category"), "")
	if err != nil {
		return nil, fmt.Errorf("failed to insert project: %w", err)
	}
	
	projectID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get project ID: %w", err)
	}
	
	if techStack, ok := data["techStack"].([]string); ok {
		stmt, _ := tx.PrepareContext(ctx, "INSERT INTO project_tech_stack (project_id, tech) VALUES (?, ?)")
		for _, tech := range techStack {
			stmt.ExecContext(ctx, projectID, tech)
		}
		stmt.Close()
	}
	
	if images, ok := data["images"].([]string); ok {
		stmt, _ := tx.PrepareContext(ctx, "INSERT INTO project_images (project_id, image_url, sort_order) VALUES (?, ?, ?)")
		for i, img := range images {
			stmt.ExecContext(ctx, projectID, img, i)
		}
		stmt.Close()
	}
	
	tx.ExecContext(ctx, "INSERT INTO project_authors (project_id, user_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE project_id=project_id", projectID, userID)
	
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	var username sql.NullString
	r.db.QueryRowContext(ctx, "SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	
	return map[string]interface{}{
		"resourceId": projectID, "resourceType": "resource", "resource": getString(data, "github"),
		"auditStatus": "pending", "submitTime": time.Now().Format("2006-01-02 15:04:05"),
		"auditTime": nil, "rejectReason": nil, "submitor": username.String,
	}, nil
}

func (r *projectRepository) Update(ctx context.Context, userID int, projectID string, data map[string]interface{}) (map[string]interface{}, error) {
	// 检查用户是否是项目作者
	var authorCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM project_authors WHERE project_id = ? AND user_id = ?", projectID, userID).Scan(&authorCount)
	if err != nil || authorCount == 0 {
		return nil, fmt.Errorf("unauthorized: user is not the author of this project")
	}
	
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	_, err = tx.ExecContext(ctx, `UPDATE projects SET name = ?, description = ?, detail = ?, github_url = ?, category = ?, updated_at = NOW() WHERE project_id = ?`,
		getString(data, "name"), getString(data, "description"), getString(data, "detail"), getString(data, "github"), getString(data, "category"), projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	
	// 更新技术栈
	if techStack, ok := data["techStack"].([]string); ok {
		tx.ExecContext(ctx, "DELETE FROM project_tech_stack WHERE project_id = ?", projectID)
		stmt, _ := tx.PrepareContext(ctx, "INSERT INTO project_tech_stack (project_id, tech) VALUES (?, ?)")
		for _, tech := range techStack {
			stmt.ExecContext(ctx, projectID, tech)
		}
		stmt.Close()
	}
	
	// 更新图片
	if images, ok := data["images"].([]string); ok {
		tx.ExecContext(ctx, "DELETE FROM project_images WHERE project_id = ?", projectID)
		stmt, _ := tx.PrepareContext(ctx, "INSERT INTO project_images (project_id, image_url, sort_order) VALUES (?, ?, ?)")
		for i, img := range images {
			stmt.ExecContext(ctx, projectID, img, i)
		}
		stmt.Close()
	}
	
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	project, _ := r.GetByID(ctx, projectID)
	return project, nil
}

func (r *projectRepository) LikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	// 先检查是否已点赞
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM likes WHERE user_id = ? AND resource_type = 'project' AND resource_id = ?", userID, projectID).Scan(&exists)
	
	if exists > 0 {
		// 已点赞，直接返回
		var likes int
		r.db.QueryRowContext(ctx, "SELECT loves FROM projects WHERE project_id = ?", projectID).Scan(&likes)
		return map[string]interface{}{"likecounts": likes, "isliked": true}, nil
	}
	
	// 未点赞，执行插入
	_, err := r.db.ExecContext(ctx, "INSERT INTO likes (user_id, resource_type, resource_id, created_at) VALUES (?, 'project', ?, NOW())", userID, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to like project: %w", err)
	}
	
	// 增加点赞数
	_, err = r.db.ExecContext(ctx, "UPDATE projects SET loves = loves + 1 WHERE project_id = ?", projectID)
	var likes int
	r.db.QueryRowContext(ctx, "SELECT loves FROM projects WHERE project_id = ?", projectID).Scan(&likes)
	return map[string]interface{}{"likecounts": likes, "isliked": true}, err
}

func (r *projectRepository) UnlikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	// 先检查是否已点赞，只有已点赞才执行删除和减少计数
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM likes WHERE user_id = ? AND resource_type = 'project' AND resource_id = ?", userID, projectID).Scan(&exists)
	
	if exists == 0 {
		// 未点赞，直接返回当前状态
		var likes int
		r.db.QueryRowContext(ctx, "SELECT loves FROM projects WHERE project_id = ?", projectID).Scan(&likes)
		return map[string]interface{}{"likecounts": likes, "isliked": false}, nil
	}
	
	// 已点赞，执行删除
	result, err := r.db.ExecContext(ctx, "DELETE FROM likes WHERE user_id = ? AND resource_type = 'project' AND resource_id = ?", userID, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to unlike project: %w", err)
	}
	
	// 检查删除是否成功（受影响的行数）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	
	// 只有成功删除才减少计数
	if rowsAffected > 0 {
		_, err = r.db.ExecContext(ctx, "UPDATE projects SET loves = GREATEST(loves - 1, 0) WHERE project_id = ?", projectID)
	}
	
	var likes int
	r.db.QueryRowContext(ctx, "SELECT loves FROM projects WHERE project_id = ?", projectID).Scan(&likes)
	return map[string]interface{}{"likecounts": likes, "isliked": false}, err
}

func (r *projectRepository) AddComment(ctx context.Context, userID int, projectID, content string) (map[string]interface{}, error) {
	// 修复：使用 COALESCE 确保如果昵称为空，使用用户名作为备用，但字段名仍然是 nickname
	var nickname, username, avatar sql.NullString
	err := r.db.QueryRowContext(ctx, "SELECT COALESCE(nickname, username) as nickname, username, COALESCE(avatar, '') as avatar FROM users WHERE id = ?", userID).Scan(&nickname, &username, &avatar)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	result, err := r.db.ExecContext(ctx, "INSERT INTO comments (resource_type, resource_id, user_id, content, created_at, updated_at) VALUES ('project', ?, ?, ?, NOW(), NOW())", projectID, userID, content)
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
		"comment_Id": commentID, // 兼容字段
		"userId":     userID,
		"user_id":    userID, // 兼容字段
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

// GetComments 获取评论列表
func (r *projectRepository) GetComments(ctx context.Context, projectID string, cursor, limit int) ([]map[string]interface{}, error) {
	query := `SELECT c.comment_id, c.user_id, c.content, c.love_count, c.reply_total, c.created_at, COALESCE(u.nickname, u.username) as nickname, COALESCE(u.avatar, '') as avatar FROM comments c JOIN users u ON c.user_id = u.id WHERE c.resource_type = 'project' AND c.resource_id = ? AND c.parent_id IS NULL AND c.deleted_at IS NULL ORDER BY c.created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, projectID, limit, cursor)
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
		replyRows, _ := r.db.QueryContext(ctx, `SELECT c.comment_id, c.user_id, c.content, c.love_count, c.created_at, COALESCE(u.nickname, u.username) as nickname, COALESCE(u.avatar, '') as avatar, c.parent_id FROM comments c JOIN users u ON c.user_id = u.id WHERE c.parent_id = ? AND c.deleted_at IS NULL ORDER BY c.created_at ASC`, commentID)
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
						"comment_Id": replyID, "nickname": replyNickname.String, "avater": replyAvatar.String,
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
			"userId": userID, // 前端可能需要userId来判断是否可以删除
		})
	}
	return comments, nil
}

func (r *projectRepository) DeleteComment(ctx context.Context, userID int, projectID, commentID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE comments SET deleted_at = NOW() WHERE comment_id = ? AND user_id = ? AND resource_type = 'project' AND resource_id = ?", commentID, userID, projectID)
	return err
}

// LikeComment 点赞评论，返回点赞后的状态和点赞数
func (r *projectRepository) LikeComment(ctx context.Context, userID int, commentID string) (bool, int, error) {
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&exists)
	isLiked := exists > 0
	
	if isLiked {
		// 取消点赞
		_, err := r.db.ExecContext(ctx, "DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)
		if err != nil {
			return false, 0, err
		}
		r.db.ExecContext(ctx, "UPDATE comments SET love_count = GREATEST(love_count - 1, 0) WHERE comment_id = ?", commentID)
		isLiked = false
	} else {
		// 点赞
		_, err := r.db.ExecContext(ctx, "INSERT INTO comment_likes (comment_id, user_id, created_at) VALUES (?, ?, NOW())", commentID, userID)
		if err != nil {
			return false, 0, err
		}
		r.db.ExecContext(ctx, "UPDATE comments SET love_count = love_count + 1 WHERE comment_id = ?", commentID)
		isLiked = true
	}
	
	// 获取更新后的点赞数
	var likes int
	r.db.QueryRowContext(ctx, "SELECT love_count FROM comments WHERE comment_id = ?", commentID).Scan(&likes)
	
	return isLiked, likes, nil
}

func (r *projectRepository) ReplyComment(ctx context.Context, userID int, projectID, commentID, content string) (map[string]interface{}, error) {
	var username, avatar sql.NullString
	err := r.db.QueryRowContext(ctx, "SELECT username, avatar FROM users WHERE id = ?", userID).Scan(&username, &avatar)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	result, err := r.db.ExecContext(ctx, "INSERT INTO comments (resource_type, resource_id, parent_id, user_id, content, created_at, updated_at) VALUES ('project', ?, ?, ?, ?, NOW(), NOW())", projectID, commentID, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to add reply: %w", err)
	}
	replyID, _ := result.LastInsertId()
	r.db.ExecContext(ctx, "UPDATE comments SET reply_total = reply_total + 1 WHERE comment_id = ?", commentID)
	return map[string]interface{}{
		"comment_Id": replyID, "nickname": username.String, "avater": avatar.String,
		"comment": content, "commentDate": time.Now().Format("2006-01-02 15:04:05"),
		"love_count": 0, "isreply": true, "reply_id": commentID,
	}, nil
}

func (r *projectRepository) DeleteReply(ctx context.Context, userID int, projectID, commentID string) (map[string]interface{}, error) {
	var parentID sql.NullInt64
	r.db.QueryRowContext(ctx, "SELECT parent_id FROM comments WHERE comment_id = ?", commentID).Scan(&parentID)
	_, err := r.db.ExecContext(ctx, "UPDATE comments SET deleted_at = NOW() WHERE comment_id = ? AND user_id = ? AND resource_type = 'project' AND resource_id = ?", commentID, userID, projectID)
	if err == nil && parentID.Valid {
		r.db.ExecContext(ctx, "UPDATE comments SET reply_total = GREATEST(reply_total - 1, 0) WHERE comment_id = ?", parentID.Int64)
	}
	return map[string]interface{}{}, err
}

func (r *projectRepository) AddView(ctx context.Context, projectID string) (int, error) {
	_, err := r.db.ExecContext(ctx, "UPDATE projects SET views = views + 1 WHERE project_id = ?", projectID)
	if err != nil {
		return 0, err
	}
	var views int
	r.db.QueryRowContext(ctx, "SELECT views FROM projects WHERE project_id = ?", projectID).Scan(&views)
	return views, nil
}

func (r *projectRepository) CollectProject(ctx context.Context, userID int, projectID string) error {
	// 先检查是否已收藏
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM collections WHERE user_id = ? AND resource_type = 'project' AND resource_id = ?", userID, projectID).Scan(&exists)
	
	if exists > 0 {
		// 已收藏，直接返回（不做任何操作）
		return nil
	}
	
	// 未收藏，执行插入
	_, err := r.db.ExecContext(ctx, "INSERT INTO collections (user_id, resource_type, resource_id, created_at) VALUES (?, 'project', ?, NOW())", userID, projectID)
	if err != nil {
		return fmt.Errorf("failed to collect project: %w", err)
	}
	
	// 增加收藏数
	_, err = r.db.ExecContext(ctx, "UPDATE projects SET collections = collections + 1 WHERE project_id = ?", projectID)
	return err
}

func (r *projectRepository) UncollectProject(ctx context.Context, userID int, projectID string) error {
	// 先检查是否已收藏，只有已收藏才执行删除和减少计数
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM collections WHERE user_id = ? AND resource_type = 'project' AND resource_id = ?", userID, projectID).Scan(&exists)
	
	if exists == 0 {
		// 未收藏，直接返回
		return nil
	}
	
	// 已收藏，执行删除
	result, err := r.db.ExecContext(ctx, "DELETE FROM collections WHERE user_id = ? AND resource_type = 'project' AND resource_id = ?", userID, projectID)
	if err != nil {
		return fmt.Errorf("failed to uncollect project: %w", err)
	}
	
	// 检查删除是否成功（受影响的行数）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	// 只有成功删除才减少计数
	if rowsAffected > 0 {
		_, err = r.db.ExecContext(ctx, "UPDATE projects SET collections = GREATEST(collections - 1, 0) WHERE project_id = ?", projectID)
	}
	return err
}

// GetCollectionCount 从 collections 表实时统计收藏数
func (r *projectRepository) GetCollectionCount(ctx context.Context, projectID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, 
		"SELECT COUNT(*) FROM collections WHERE resource_type = 'project' AND resource_id = ?", 
		projectID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get collection count: %w", err)
	}
	return count, nil
}

func (r *projectRepository) GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error) {
	// 查询status为pending的项目
	// 使用子查询获取第一个作者作为提交者
	query := `SELECT p.project_id, p.name, p.description, p.category, p.github_url, p.status, p.created_at, p.audit_time, p.reject_reason, 
	          (SELECT u.username FROM project_authors pa2 JOIN users u ON pa2.user_id = u.id WHERE pa2.project_id = p.project_id LIMIT 1) as username
	          FROM projects p 
	          WHERE p.status = 'pending' 
	          ORDER BY p.created_at DESC 
	          LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending projects: %w", err)
	}
	defer rows.Close()
	var projects []map[string]interface{}
	for rows.Next() {
		var projectID int
		var name, description, category, githubURL, status sql.NullString
		var createdAt time.Time
		var auditTime sql.NullTime
		var rejectReason, username sql.NullString
		// SQL查询顺序：project_id, name, description, category, github_url, status, created_at, audit_time, reject_reason, username
		if err := rows.Scan(&projectID, &name, &description, &category, &githubURL, &status, &createdAt, &auditTime, &rejectReason, &username); err != nil {
			continue
		}
		techRows, _ := r.db.QueryContext(ctx, "SELECT tech FROM project_tech_stack WHERE project_id = ?", projectID)
		var techList []string
		if techRows != nil {
			defer techRows.Close()
			for techRows.Next() {
				var tech string
				if err := techRows.Scan(&tech); err == nil {
					techList = append(techList, tech)
				}
			}
		}
		projects = append(projects, map[string]interface{}{
			"id":           projectID,
			"resourceId":   projectID,
			"_id":          projectID, // 前端可能使用_id
			"reourceId":    projectID, // 保持兼容性
			"name":         name.String,
			"title":        name.String, // 前端可能使用title
			"resourcename": name.String, // 保持兼容性
			"githubUrl":    githubURL.String,
			"demoUrl":      "", // 项目表没有demoUrl字段
			"description":  description.String,
			"category":     category.String,
			"catagory":     category.String, // 保持兼容性
			"technologies": techList,
			"tags":         techList, // 前端可能使用tags
			"uploader":     username.String,
			"submitor":     username.String, // 保持兼容性
			"author":       username.String, // 前端可能使用author
			"owner":        username.String, // 前端可能使用owner
			"status":       status.String,
			"created_at":   createdAt.Format("2006-01-02T15:04:05Z"),
			"created":      createdAt.Format("2006-01-02T15:04:05Z"), // 前端可能使用created
			"submitDate":   createdAt.Format("2006-01-02 15:04:05"), // 保持兼容性
			"auditTime": func() interface{} {
				if auditTime.Valid {
					return auditTime.Time.Format("2006-01-02 15:04:05")
				}
				return nil
			}(),
			"rejectReason": rejectReason.String,
		})
	}
	return projects, nil
}

// UpdateProjectStatus 更新项目审核状态
func (r *projectRepository) UpdateProjectStatus(ctx context.Context, resourceID, status, rejectReason string) error {
	query := "UPDATE projects SET status = ?, audit_time = NOW()"
	args := []interface{}{status}
	
	if rejectReason != "" {
		query += ", reject_reason = ?"
		args = append(args, rejectReason)
	}
	
	query += " WHERE project_id = ?"
	args = append(args, resourceID)
	
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
