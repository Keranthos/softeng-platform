package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"softeng-platform/internal/model"
	"strings"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByNickname(ctx context.Context, nickname string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, userID int, hashedPassword string) error
	GetCollections(ctx context.Context, userID int, resourceType string) ([]map[string]interface{}, error)
	DeleteCollection(ctx context.Context, userID int, resourceType string, resourceID int) error
	GetUserPendingItems(ctx context.Context, userID int, resourceType string) ([]map[string]interface{}, error)
	GetUserSubmissions(ctx context.Context, userID int, resourceType string) ([]map[string]interface{}, error)
	GetResourceStatus(ctx context.Context, resourceType, resourceID string, userID int) error
	GetResourceOldStatus(ctx context.Context, resourceType, resourceID string) (string, error)
	LogStatusChange(ctx context.Context, resourceType, resourceID, oldStatus, newStatus string, operatorID int) error
	UpdateEmail(ctx context.Context, userID int, newEmail string) error
}

type userRepository struct {
	db *Database
}

func NewUserRepository(db *Database) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, nickname, email, password, avatar, description, face_photo, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Nickname,
		user.Email,
		user.Password,
		user.Avatar,
		user.Description,
		user.FacePhoto,
		user.Role,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT id, username, nickname, email, 
		       COALESCE(avatar, '') as avatar, 
		       COALESCE(description, '') as description, 
		       COALESCE(face_photo, '') as face_photo, 
		       role, created_at, updated_at
		FROM users WHERE id = ?
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.Avatar,
		&user.Description,
		&user.FacePhoto,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, nickname, email, password, 
		       COALESCE(avatar, '') as avatar, 
		       COALESCE(description, '') as description, 
		       COALESCE(face_photo, '') as face_photo, 
		       role, created_at, updated_at
		FROM users WHERE username = ?
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Avatar,
		&user.Description,
		&user.FacePhoto,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, nickname, email, password, 
		       COALESCE(avatar, '') as avatar, 
		       COALESCE(description, '') as description, 
		       COALESCE(face_photo, '') as face_photo, 
		       role, created_at, updated_at
		FROM users WHERE email = ?
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Avatar,
		&user.Description,
		&user.FacePhoto,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByNickname(ctx context.Context, nickname string) (*model.User, error) {
	query := `
		SELECT id, username, nickname, email, password, 
		       COALESCE(avatar, '') as avatar, 
		       COALESCE(description, '') as description, 
		       COALESCE(face_photo, '') as face_photo, 
		       role, created_at, updated_at
		FROM users WHERE nickname = ?
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, nickname).Scan(
		&user.ID,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.Password,
		&user.Avatar,
		&user.Description,
		&user.FacePhoto,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users 
		SET nickname = ?, avatar = ?, description = ?, face_photo = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Nickname,
		user.Avatar,
		user.Description,
		user.FacePhoto,
		time.Now(),
		user.ID,
	)

	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	query := `UPDATE users SET password = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, hashedPassword, time.Now(), userID)
	return err
}

func (r *userRepository) UpdateEmail(ctx context.Context, userID int, newEmail string) error {
	query := `UPDATE users SET email = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, newEmail, time.Now(), userID)
	return err
}

func (r *userRepository) GetResourceStatus(ctx context.Context, resourceType, resourceID string, userID int) error {
	// 验证用户是否有权限操作该资源
	switch resourceType {
	case "tool":
		var submitterID int
		err := r.db.QueryRowContext(ctx, "SELECT submitter_id FROM tools WHERE resource_id = ?", resourceID).Scan(&submitterID)
		if err != nil {
			return fmt.Errorf("resource not found")
		}
		if submitterID != userID {
			return fmt.Errorf("unauthorized: user is not the submitter")
		}
	case "project":
		var authorCount int
		err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM project_authors WHERE project_id = ? AND user_id = ?", resourceID, userID).Scan(&authorCount)
		if err != nil {
			return fmt.Errorf("resource not found")
		}
		if authorCount == 0 {
			return fmt.Errorf("unauthorized: user is not the author")
		}
	}
	return nil
}

func (r *userRepository) GetResourceOldStatus(ctx context.Context, resourceType, resourceID string) (string, error) {
	switch resourceType {
	case "tool":
		var status string
		err := r.db.QueryRowContext(ctx, "SELECT status FROM tools WHERE resource_id = ?", resourceID).Scan(&status)
		if err != nil {
			return "", err
		}
		return status, nil
	case "project":
		var status string
		err := r.db.QueryRowContext(ctx, "SELECT status FROM projects WHERE project_id = ?", resourceID).Scan(&status)
		if err != nil {
			return "", err
		}
		return status, nil
	}
	return "", nil
}


func (r *userRepository) LogStatusChange(ctx context.Context, resourceType, resourceID, oldStatus, newStatus string, operatorID int) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO resource_status_logs (resource_type, resource_id, old_status, new_status, operator_id, operate_time) VALUES (?, ?, ?, ?, ?, NOW())`,
		resourceType, resourceID, oldStatus, newStatus, operatorID)
	return err
}

func (r *userRepository) GetCollections(ctx context.Context, userID int, resourceType string) ([]map[string]interface{}, error) {
	query := `SELECT c.resource_id, c.created_at FROM collections c WHERE c.user_id = ? AND c.resource_type = ? ORDER BY c.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID, resourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to query collections: %w", err)
	}
	defer rows.Close()
	
	var items []map[string]interface{}
	for rows.Next() {
		var resourceID int
		var createdAt time.Time
		if err := rows.Scan(&resourceID, &createdAt); err != nil {
			fmt.Printf("[DEBUG] Failed to scan collection row: %v\n", err)
			continue
		}
		fmt.Printf("[DEBUG] Processing %s with ID=%d, createdAt=%v\n", resourceType, resourceID, createdAt)
		
		// 根据资源类型获取详细信息
		var item map[string]interface{}
		switch resourceType {
		case "tool":
			var name, category, description sql.NullString
			var views, collections, loves int
			r.db.QueryRowContext(ctx, "SELECT resource_name, category, description, views, collections, loves FROM tools WHERE resource_id = ?", resourceID).Scan(&name, &category, &description, &views, &collections, &loves)
			
			// 获取工具标签
			var tags []string
			tagRows, tagErr := r.db.QueryContext(ctx, "SELECT tag FROM tool_tags WHERE tool_id = ? LIMIT 1", resourceID)
			if tagErr != nil {
				fmt.Printf("[DEBUG] Failed to query tool tags for tool_id=%d: %v\n", resourceID, tagErr)
			}
			if tagRows != nil {
				defer tagRows.Close()
				for tagRows.Next() {
					var tag sql.NullString
					if tagRows.Scan(&tag) == nil && tag.Valid {
						tags = append(tags, tag.String)
					}
				}
			}
			fmt.Printf("[DEBUG] Tool ID=%d, tags=%v\n", resourceID, tags)
			
			// 获取工具图片（从 tool_images 表）
			var images []string
			imageRows, _ := r.db.QueryContext(ctx, "SELECT image_url FROM tool_images WHERE tool_id = ? ORDER BY sort_order LIMIT 1", resourceID)
			if imageRows != nil {
				defer imageRows.Close()
				for imageRows.Next() {
					var imgURL sql.NullString
					if imageRows.Scan(&imgURL) == nil && imgURL.Valid {
						images = append(images, imgURL.String)
						break // 只需要第一张图片
					}
				}
			}
			
			// 确保即使数组为空也返回空数组，而不是 nil
			if tags == nil {
				tags = []string{}
			}
			if images == nil {
				images = []string{}
			}
			
			item = map[string]interface{}{
				"resourceId": resourceID, "resourceType": resourceType, "resourceName": name.String,
				"name": name.String, "introduce": description.String, "description": description.String,
				"catagory": category.String, "category": category.String,
				"tags": tags, "image": images,
				"views": views, "collections": collections, "loves": loves,
				"created_at": createdAt.Format("2006-01-02 15:04:05"), // 收藏时间
			}
		case "course":
			var name, semester, description, cover sql.NullString
			var views, loves, collections int
			r.db.QueryRowContext(ctx, "SELECT name, semester, description, cover, views, loves, collections FROM courses WHERE course_id = ?", resourceID).Scan(&name, &semester, &description, &cover, &views, &loves, &collections)
			
			// 获取课程教师（字段名是 teacher_name，不是 teacher）
			var teachers []string
			teacherRows, teacherErr := r.db.QueryContext(ctx, "SELECT teacher_name FROM course_teachers WHERE course_id = ? LIMIT 1", resourceID)
			if teacherErr != nil {
				fmt.Printf("[DEBUG] Failed to query course teachers for course_id=%d: %v\n", resourceID, teacherErr)
			}
			if teacherRows != nil {
				defer teacherRows.Close()
				for teacherRows.Next() {
					var teacher sql.NullString
					if teacherRows.Scan(&teacher) == nil && teacher.Valid {
						teachers = append(teachers, teacher.String)
					}
				}
			}
			fmt.Printf("[DEBUG] Course ID=%d, teachers=%v\n", resourceID, teachers)
			
			// 确保即使数组为空也返回空数组，而不是 nil
			if teachers == nil {
				teachers = []string{}
			}
			
			item = map[string]interface{}{
				"resourceId": resourceID, "resourceType": resourceType, "resourceName": name.String,
				"name": name.String, "courseName": name.String, "introduce": description.String, "description": description.String,
				"semester": semester.String, "cover": cover.String,
				"teacher": teachers, // 教师列表
				"views": views, "loves": loves, "collections": collections,
				"created_at": createdAt.Format("2006-01-02 15:04:05"), // 收藏时间
			}
		case "project":
			var name, category, description, cover sql.NullString
			var views, loves, collections int
			r.db.QueryRowContext(ctx, "SELECT name, category, description, cover, views, loves, collections FROM projects WHERE project_id = ?", resourceID).Scan(&name, &category, &description, &cover, &views, &loves, &collections)
			
			// 获取项目技术栈
			var techStack []string
			techRows, techErr := r.db.QueryContext(ctx, "SELECT tech FROM project_tech_stack WHERE project_id = ? LIMIT 1", resourceID)
			if techErr != nil {
				fmt.Printf("[DEBUG] Failed to query project tech stack for project_id=%d: %v\n", resourceID, techErr)
			}
			if techRows != nil {
				defer techRows.Close()
				for techRows.Next() {
					var tech sql.NullString
					if techRows.Scan(&tech) == nil && tech.Valid {
						techStack = append(techStack, tech.String)
					}
				}
			}
			fmt.Printf("[DEBUG] Project ID=%d, techStack=%v\n", resourceID, techStack)
			
			// 获取项目图片
			var imageList []string
			if cover.Valid && cover.String != "" {
				imageList = []string{cover.String}
			}
			
			// 确保即使数组为空也返回空数组，而不是 nil
			if techStack == nil {
				techStack = []string{}
			}
			if imageList == nil {
				imageList = []string{}
			}
			
			item = map[string]interface{}{
				"resourceId": resourceID, "resourceType": resourceType, "resourceName": name.String,
				"name": name.String, "projectName": name.String, "introduce": description.String, "description": description.String,
				"catagory": category.String, "category": category.String,
				"techStack": techStack, "technologies": techStack, // 技术栈
				"coverImage": cover.String, "cover": cover.String, "image": imageList,
				"views": views, "loves": loves, "collections": collections,
				"created_at": createdAt.Format("2006-01-02 15:04:05"), // 收藏时间
			}
		}
		if item != nil {
			// 调试：打印生成的 item 数据
			if resourceType == "tool" {
				fmt.Printf("[DEBUG] Tool item: %+v\n", item)
			} else if resourceType == "project" {
				fmt.Printf("[DEBUG] Project item: %+v\n", item)
			} else if resourceType == "course" {
				fmt.Printf("[DEBUG] Course item: %+v\n", item)
			}
			items = append(items, item)
		}
	}
	return items, nil
}

func (r *userRepository) DeleteCollection(ctx context.Context, userID int, resourceType string, resourceID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM collections WHERE user_id = ? AND resource_type = ? AND resource_id = ?", userID, resourceType, resourceID)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}
	
	// 更新资源的收藏数
	switch resourceType {
	case "tool":
		r.db.ExecContext(ctx, "UPDATE tools SET collections = GREATEST(collections - 1, 0) WHERE resource_id = ?", resourceID)
	case "course":
		r.db.ExecContext(ctx, "UPDATE courses SET collections = GREATEST(collections - 1, 0) WHERE course_id = ?", resourceID)
	case "project":
		r.db.ExecContext(ctx, "UPDATE projects SET collections = GREATEST(collections - 1, 0) WHERE project_id = ?", resourceID)
	}
	
	return nil
}

func (r *userRepository) GetUserPendingItems(ctx context.Context, userID int, resourceType string) ([]map[string]interface{}, error) {
	switch resourceType {
	case "tool":
		// 返回所有状态的记录（pending, approved, rejected），而不仅仅是pending
		query := `SELECT t.resource_id, t.resource_name, t.category, t.status, t.created_at, t.audit_time, t.reject_reason FROM tools t WHERE t.submitter_id = ? ORDER BY t.created_at DESC`
		rows, err := r.db.QueryContext(ctx, query, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []map[string]interface{}
		for rows.Next() {
			var resourceID int
			var name, category, status sql.NullString
			var createdAt time.Time
			var auditTime sql.NullTime
			var rejectReason sql.NullString
			if err := rows.Scan(&resourceID, &name, &category, &status, &createdAt, &auditTime, &rejectReason); err != nil {
				continue
			}
			
			// 获取工具图片（从 tool_images 表）
			var images []string
			imageRows, _ := r.db.QueryContext(ctx, "SELECT image_url FROM tool_images WHERE tool_id = ? ORDER BY sort_order LIMIT 1", resourceID)
			if imageRows != nil {
				defer imageRows.Close()
				for imageRows.Next() {
					var imgURL sql.NullString
					if imageRows.Scan(&imgURL) == nil && imgURL.Valid {
						images = append(images, imgURL.String)
						break // 只需要第一张图片
					}
				}
			}
			if images == nil {
				images = []string{}
			}
			
			// 将 status 映射为 auditStatus，并添加 introduce 字段
			auditStatus := "pending"
			if status.Valid {
				auditStatus = status.String
			}
			items = append(items, map[string]interface{}{
				"resourceId": resourceID, "resourceType": resourceType, "resourceName": name.String,
				"introduce": name.String, // 添加 introduce 字段供前端使用
				"image": images, "images": images, // 添加图片字段
				"catagory": category.String, "status": auditStatus, "auditStatus": auditStatus, // 同时返回 status 和 auditStatus
				"submitTime": createdAt.Format("2006-01-02 15:04:05"),
				"auditTime": func() interface{} {
					if auditTime.Valid {
						return auditTime.Time.Format("2006-01-02 15:04:05")
					}
					return nil
				}(), "rejectReason": rejectReason.String,
			})
		}
		return items, nil
	case "course":
		// 课程没有status字段，暂时返回空
		return []map[string]interface{}{}, nil
	case "project":
		// 返回所有状态的记录（pending, approved, rejected），而不仅仅是pending
		query := `SELECT p.project_id, p.name, p.category, p.status, p.created_at, p.audit_time, p.reject_reason 
		          FROM projects p 
		          WHERE EXISTS (SELECT 1 FROM project_authors pa WHERE pa.project_id = p.project_id AND pa.user_id = ?) 
		          ORDER BY p.created_at DESC`
		rows, err := r.db.QueryContext(ctx, query, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []map[string]interface{}
		for rows.Next() {
			var projectID int
			var name, category, status sql.NullString
			var createdAt time.Time
			var auditTime sql.NullTime
			var rejectReason sql.NullString
			if err := rows.Scan(&projectID, &name, &category, &status, &createdAt, &auditTime, &rejectReason); err != nil {
				continue
			}
			// 将 status 映射为 auditStatus，并添加 introduce 字段
			auditStatus := "pending"
			if status.Valid {
				auditStatus = status.String
			}
			items = append(items, map[string]interface{}{
				"resourceId": projectID, "resourceType": resourceType, "resourceName": name.String,
				"introduce": name.String, // 添加 introduce 字段供前端使用
				"catagory": category.String, "status": auditStatus, "auditStatus": auditStatus, // 同时返回 status 和 auditStatus
				"submitTime": createdAt.Format("2006-01-02 15:04:05"),
				"auditTime": func() interface{} {
					if auditTime.Valid {
						return auditTime.Time.Format("2006-01-02 15:04:05")
					}
					return nil
				}(), "rejectReason": rejectReason.String,
			})
		}
		return items, nil
	}
	return []map[string]interface{}{}, nil
}

func (r *userRepository) GetUserSubmissions(ctx context.Context, userID int, resourceType string) ([]map[string]interface{}, error) {
	switch resourceType {
	case "tool":
		// 先查询所有工具，避免在查询时处理图片导致性能问题
		query := `SELECT t.resource_id, t.resource_name, t.category, t.status, t.created_at, t.audit_time, t.reject_reason
		          FROM tools t 
		          WHERE t.submitter_id = ? 
		          ORDER BY t.created_at DESC`
		rows, err := r.db.QueryContext(ctx, query, userID)
		if err != nil {
			fmt.Printf("[GetUserSubmissions] 查询工具失败: %v\n", err)
			return nil, err
		}
		defer rows.Close()
		
		var toolIDs []int
		var items []map[string]interface{}
		toolMap := make(map[int]*map[string]interface{})
		
		// 先收集所有工具ID和基本信息
		for rows.Next() {
			var resourceID int
			var name, category, status sql.NullString
			var createdAt time.Time
			var auditTime sql.NullTime
			var rejectReason sql.NullString
			if err := rows.Scan(&resourceID, &name, &category, &status, &createdAt, &auditTime, &rejectReason); err != nil {
				continue
			}
			
			toolIDs = append(toolIDs, resourceID)
			
			// 将 status 映射为 auditStatus
			auditStatus := "pending"
			if status.Valid {
				auditStatus = status.String
			}
			
			item := map[string]interface{}{
				"resourceId": resourceID, "resourceType": resourceType, "resourceName": name.String,
				"introduce": name.String,
				"image": []string{}, "images": []string{}, // 先设置为空，后面再填充
				"catagory": category.String, "status": auditStatus, "auditStatus": auditStatus,
				"submitTime": createdAt.Format("2006-01-02 15:04:05"),
				"auditTime": func() interface{} {
					if auditTime.Valid {
						return auditTime.Time.Format("2006-01-02 15:04:05")
					}
					return nil
				}(), "rejectReason": rejectReason.String,
			}
			items = append(items, item)
			toolMap[resourceID] = &item
		}
		
		// 如果有工具，批量查询所有图片（简化查询，避免复杂JOIN导致性能问题）
		if len(toolIDs) > 0 {
			// 构建IN子句的占位符
			placeholders := make([]string, len(toolIDs))
			args := make([]interface{}, len(toolIDs))
			for i, id := range toolIDs {
				placeholders[i] = "?"
				args[i] = id
			}
			
			// 简化查询：直接查询所有图片，然后在代码中处理
			// 这样避免复杂的子查询和JOIN，性能更好
			imageQuery := `SELECT tool_id, image_url, sort_order, id
			               FROM tool_images 
			               WHERE tool_id IN (` + strings.Join(placeholders, ",") + `)
			               ORDER BY tool_id, sort_order ASC, id ASC`
			
			imageRows, err := r.db.QueryContext(ctx, imageQuery, args...)
			if err != nil {
				// 如果查询失败，记录错误但不影响主流程（工具列表仍然返回，只是没有图片）
				fmt.Printf("[GetUserSubmissions] 查询图片失败 (toolIDs=%v): %v\n", toolIDs, err)
			} else {
				defer imageRows.Close()
				// 使用map记录每个工具的第一张图片
				firstImageMap := make(map[int]string)
				for imageRows.Next() {
					var toolID int
					var imageURL sql.NullString
					var sortOrder sql.NullInt64
					var imageID sql.NullInt64
					if err := imageRows.Scan(&toolID, &imageURL, &sortOrder, &imageID); err != nil {
						fmt.Printf("[GetUserSubmissions] 扫描图片数据失败: %v\n", err)
						continue
					}
					// 只记录每个工具的第一张图片
					if _, exists := firstImageMap[toolID]; !exists && imageURL.Valid && imageURL.String != "" {
						firstImageMap[toolID] = imageURL.String
					}
				}
				// 将图片URL填充到对应的工具项中
				for toolID, imageURL := range firstImageMap {
					if item, ok := toolMap[toolID]; ok {
						images := []string{imageURL}
						(*item)["image"] = images
						(*item)["images"] = images
					}
				}
				fmt.Printf("[GetUserSubmissions] 成功加载 %d 个工具的图片\n", len(firstImageMap))
			}
		}
		
		return items, nil
	case "course":
		// 使用 INNER JOIN 替代 EXISTS，性能更好
		query := `SELECT DISTINCT c.course_id, c.name, c.semester, c.created_at 
		          FROM courses c 
		          INNER JOIN course_contributors cc ON cc.course_id = c.course_id 
		          WHERE cc.user_id = ? 
		          ORDER BY c.created_at DESC`
		rows, err := r.db.QueryContext(ctx, query, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []map[string]interface{}
		for rows.Next() {
			var courseID int
			var name, semester sql.NullString
			var createdAt time.Time
			if err := rows.Scan(&courseID, &name, &semester, &createdAt); err != nil {
				continue
			}
			// 课程没有 status 字段，默认设置为 pending
			items = append(items, map[string]interface{}{
				"resourceId": courseID, "resourceType": resourceType, "resourceName": name.String,
				"introduce": name.String, // 添加 introduce 字段供前端使用
				"semester": semester.String, "submitTime": createdAt.Format("2006-01-02 15:04:05"),
				"status": "pending", "auditStatus": "pending", // 课程默认 pending
			})
		}
		return items, nil
	case "project":
		// 使用 INNER JOIN 替代 EXISTS，性能更好
		query := `SELECT DISTINCT p.project_id, p.name, p.category, p.status, p.created_at, p.audit_time, p.reject_reason 
		          FROM projects p 
		          INNER JOIN project_authors pa ON pa.project_id = p.project_id 
		          WHERE pa.user_id = ? 
		          ORDER BY p.created_at DESC`
		rows, err := r.db.QueryContext(ctx, query, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var items []map[string]interface{}
		for rows.Next() {
			var projectID int
			var name, category, status sql.NullString
			var createdAt time.Time
			var auditTime sql.NullTime
			var rejectReason sql.NullString
			if err := rows.Scan(&projectID, &name, &category, &status, &createdAt, &auditTime, &rejectReason); err != nil {
				continue
			}
			// 将 status 映射为 auditStatus，并添加 introduce 字段
			auditStatus := "pending"
			if status.Valid {
				auditStatus = status.String
			}
			items = append(items, map[string]interface{}{
				"resourceId": projectID, "resourceType": resourceType, "resourceName": name.String,
				"introduce": name.String, // 添加 introduce 字段供前端使用
				"catagory": category.String, "status": auditStatus, "auditStatus": auditStatus, // 同时返回 status 和 auditStatus
				"submitTime": createdAt.Format("2006-01-02 15:04:05"),
				"auditTime": func() interface{} {
					if auditTime.Valid {
						return auditTime.Time.Format("2006-01-02 15:04:05")
					}
					return nil
				}(), "rejectReason": rejectReason.String,
			})
		}
		return items, nil
	}
	return []map[string]interface{}{}, nil
}
