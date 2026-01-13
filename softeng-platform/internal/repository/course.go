package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type CourseRepository interface {
	GetCourses(ctx context.Context, semester string, category []string, sort string, limit, cursor int) ([]map[string]interface{}, error)
	GetByID(ctx context.Context, courseID string) (map[string]interface{}, error)
	Search(ctx context.Context, keyword string, category []string, limit, cursor int) ([]map[string]interface{}, error)
	Create(ctx context.Context, userID int, data map[string]interface{}) (map[string]interface{}, error)
	UploadResource(ctx context.Context, userID int, courseID string, data map[string]interface{}) (map[string]interface{}, error)
	DownloadTextbook(ctx context.Context, courseID, textbookID string) (string, error)
	GetComments(ctx context.Context, courseID string, cursor, limit int) ([]map[string]interface{}, error)
	AddComment(ctx context.Context, userID int, courseID, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error)
	ReplyComment(ctx context.Context, userID int, courseID, commentID, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error)
	LikeComment(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error)
	GetResources(ctx context.Context, courseID string) (map[string]interface{}, error)
	AddView(ctx context.Context, courseID string) (int, error)
	CollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	UncollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	LikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	UnlikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error)
	UpdateCourseStatus(ctx context.Context, courseID, status, rejectReason string) error
	CheckUserLike(ctx context.Context, userID int, courseID string) (bool, error)
	CheckUserCollect(ctx context.Context, userID int, courseID string) (bool, error)
}

type courseRepository struct {
	db *Database
}

func NewCourseRepository(db *Database) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) GetCourses(ctx context.Context, semester string, category []string, sort string, limit, cursor int) ([]map[string]interface{}, error) {
	// 使用子查询从 collections 表实时统计收藏数，而不是使用 courses.collections 字段
	query := `SELECT c.course_id, c.resource_type, c.name, c.semester, c.credit, c.cover, c.views, c.loves, 
	          COALESCE((SELECT COUNT(*) FROM collections WHERE resource_type = 'course' AND resource_id = c.course_id), 0) as collections, 
	          c.created_at FROM courses c WHERE 1=1`
	args := []interface{}{}
	
	if semester != "" {
		query += " AND c.semester = ?"
		args = append(args, semester)
	}
	
	if len(category) > 0 {
		query += ` AND EXISTS (SELECT 1 FROM course_categories cc WHERE cc.course_id = c.course_id AND cc.category IN (` + strings.Repeat("?,", len(category))[:len(strings.Repeat("?,", len(category)))-1] + `))`
		for _, cat := range category {
			args = append(args, cat)
		}
	}
	
	switch sort {
	case "最新", "newest":
		query += " ORDER BY c.created_at DESC"
	case "最多浏览", "views":
		query += " ORDER BY c.views DESC"
	case "最多点赞", "loves":
		query += " ORDER BY c.loves DESC"
	case "最多资料":
		// 使用子查询统计的资源数量进行排序（不是收藏数）
		query += " ORDER BY (SELECT COUNT(*) FROM course_resources_web WHERE course_id = c.course_id) + (SELECT COUNT(*) FROM course_resources_upload WHERE course_id = c.course_id) DESC"
	case "学分最高":
		query += " ORDER BY c.credit DESC"
	default:
		// 默认排序：按学期和课程ID
		query += " ORDER BY c.semester ASC, c.course_id ASC"
	}
	
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, cursor)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query courses: %w", err)
	}
	defer rows.Close()
	
	var courses []map[string]interface{}
	for rows.Next() {
		var courseID, credit, views, loves, collections int
		var resourceType, name, semester, cover sql.NullString
		var createdAt time.Time
		
		if err := rows.Scan(&courseID, &resourceType, &name, &semester, &credit, &cover, &views, &loves, &collections, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}
		
		// 获取教师（取第一个）
		teacherRows, _ := r.db.QueryContext(ctx, "SELECT teacher_name FROM course_teachers WHERE course_id = ? LIMIT 1", courseID)
		var teacher string
		if teacherRows != nil {
			if teacherRows.Next() {
				teacherRows.Scan(&teacher)
			}
			teacherRows.Close()
		}
		
		// 获取分类（取第一个）
		catRows, _ := r.db.QueryContext(ctx, "SELECT category FROM course_categories WHERE course_id = ? LIMIT 1", courseID)
		var courseType string
		if catRows != nil {
			if catRows.Next() {
				catRows.Scan(&courseType)
			}
			catRows.Close()
		}
		
		// 统计课程资源数量（URL资源 + 上传资源）
		var resourceCount int
		var webCount, uploadCount int
		r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM course_resources_web WHERE course_id = ?", courseID).Scan(&webCount)
		r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM course_resources_upload WHERE course_id = ?", courseID).Scan(&uploadCount)
		resourceCount = webCount + uploadCount
		
		// collections 现在是从子查询实时统计的收藏数
		courseMap := map[string]interface{}{
			"id":         courseID,
			"courseId":   courseID, // 保留兼容
			"name":       name.String,
			"code":       "", // code字段暂不提供，如需要可以从其他表获取
			"semester":   semester.String,
			"type":       courseType,
			"teacher":    teacher,
			"credit":     float64(credit),
			"resources":  resourceCount, // 使用实际资源数量，而不是收藏数
			"likes":      loves,
			// 保留原有字段以兼容其他可能的使用
			"resourceType": resourceType.String,
			"cover":        cover.String,
			"views":        views,
			"loves":        loves,
			"collections":  collections, // 从 collections 表实时统计的收藏数
		}
		
		// 调试日志：输出收藏数
		if courseID == 1015 || courseID == 1016 { // 只输出特定课程的日志，避免日志过多
			fmt.Printf("[DEBUG] Course %d: collections=%d (from real-time query)\n", courseID, collections)
		}
		
		courses = append(courses, courseMap)
	}
	
	// 如果排序方式是"最多资料"，需要按资源数量重新排序
	if sort == "最多资料" {
		// 使用稳定排序，按资源数量降序
		for i := 0; i < len(courses)-1; i++ {
			for j := i + 1; j < len(courses); j++ {
				resI, _ := courses[i]["resources"].(int)
				resJ, _ := courses[j]["resources"].(int)
				if resI < resJ {
					courses[i], courses[j] = courses[j], courses[i]
				}
			}
		}
	}
	
	return courses, nil
}

func (r *courseRepository) GetByID(ctx context.Context, courseID string) (map[string]interface{}, error) {
	var courseIDInt int
	var resourceType, name, semester, cover sql.NullString
	var credit, views, loves, collections int
	var createdAt time.Time
	
	var description sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT course_id, resource_type, name, semester, credit, cover, views, loves, collections, COALESCE(description, '') as description, created_at FROM courses WHERE course_id = ?",
		courseID).Scan(&courseIDInt, &resourceType, &name, &semester, &credit, &cover, &views, &loves, &collections, &description, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("course not found")
		}
		return nil, fmt.Errorf("failed to get course: %w", err)
	}
	
	// 获取课程分类
	catRows, _ := r.db.QueryContext(ctx, "SELECT category FROM course_categories WHERE course_id = ?", courseIDInt)
	var categoryList []string
	if catRows != nil {
		defer catRows.Close()
		for catRows.Next() {
			var cat string
			if err := catRows.Scan(&cat); err == nil {
				categoryList = append(categoryList, cat)
			}
		}
	}
	
	// 获取课程教师
	teacherRows, _ := r.db.QueryContext(ctx, "SELECT teacher_name FROM course_teachers WHERE course_id = ?", courseIDInt)
	var teacherList []string
	if teacherRows != nil {
		defer teacherRows.Close()
		for teacherRows.Next() {
			var teacher string
			if err := teacherRows.Scan(&teacher); err == nil {
				teacherList = append(teacherList, teacher)
			}
		}
	}
	
	// 获取贡献者
	contributorRows, _ := r.db.QueryContext(ctx, "SELECT u.username FROM course_contributors cc JOIN users u ON cc.user_id = u.id WHERE cc.course_id = ?", courseIDInt)
	var contributorList []string
	if contributorRows != nil {
		defer contributorRows.Close()
		for contributorRows.Next() {
			var username string
			if err := contributorRows.Scan(&username); err == nil {
				contributorList = append(contributorList, username)
			}
		}
	}
	
	// 获取评论总数
	var commentTotal int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM comments WHERE resource_type = 'course' AND resource_id = ? AND deleted_at IS NULL", courseIDInt).Scan(&commentTotal)
	
	// 获取资源（URL和上传的资源）
	resources, _ := r.GetResources(ctx, courseID)
	urlForm, _ := resources["url_form"].([]map[string]interface{})
	uploadForm, _ := resources["upload_form"].([]map[string]interface{})
	
	return map[string]interface{}{
		"courseId":      courseIDInt,
		"resourceType":  resourceType.String,
		"name":          name.String,
		"description":   description.String, // 添加description字段
		"catagory":      func() string {
			if len(categoryList) > 0 {
				return categoryList[0]
			}
			return ""
		}(),
		"category":      categoryList, // 也提供数组格式
		"teacher":       teacherList,
		"semester":      semester.String,
		"credit":        credit,
		"cover":         cover.String,
		"url_form":      urlForm,
		"upload_form":   uploadForm,
		"contributor":   contributorList,
		"collections":   collections,
		"views":         views,
		"likes":         loves,
		"loves":         loves, // 兼容字段
		"isliked":       false, // service层会根据用户设置
		"iscollected":   false, // service层会根据用户设置
		"comment_total": commentTotal,
		"comments":      []map[string]interface{}{}, // 评论列表通过单独的接口获取
		"createdAt":     createdAt.Format("2006-01-02"),
	}, nil
}

func (r *courseRepository) Search(ctx context.Context, keyword string, category []string, limit, cursor int) ([]map[string]interface{}, error) {
	// 使用子查询从 collections 表实时统计收藏数，而不是使用 courses.collections 字段
	query := `SELECT DISTINCT c.course_id, c.resource_type, c.name, c.semester, c.credit, c.cover, c.views, c.loves, 
	          COALESCE((SELECT COUNT(*) FROM collections WHERE resource_type = 'course' AND resource_id = c.course_id), 0) as collections, 
	          c.created_at FROM courses c LEFT JOIN course_categories cc ON c.course_id = cc.course_id WHERE 1=1`
	args := []interface{}{}
	
	if keyword != "" {
		query += " AND (c.name LIKE ?)"
		args = append(args, "%"+keyword+"%")
	}
	
	if len(category) > 0 {
		query += ` AND EXISTS (SELECT 1 FROM course_categories cc2 WHERE cc2.course_id = c.course_id AND cc2.category IN (` + strings.Repeat("?,", len(category))[:len(strings.Repeat("?,", len(category)))-1] + `))`
		for _, cat := range category {
			args = append(args, cat)
		}
	}
	
	query += " ORDER BY c.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, cursor)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search courses: %w", err)
	}
	defer rows.Close()
	
	var courses []map[string]interface{}
	for rows.Next() {
		var courseID, credit, views, loves, collections int
		var resourceType, name, semester, cover sql.NullString
		var createdAt time.Time
		
		if err := rows.Scan(&courseID, &resourceType, &name, &semester, &credit, &cover, &views, &loves, &collections, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}
		
		teacherRows, _ := r.db.QueryContext(ctx, "SELECT teacher_name FROM course_teachers WHERE course_id = ?", courseID)
		var teacherList []string
		if teacherRows != nil {
			defer teacherRows.Close()
			for teacherRows.Next() {
				var teacher string
				if err := teacherRows.Scan(&teacher); err == nil {
					teacherList = append(teacherList, teacher)
				}
			}
		}
		
		catRows, _ := r.db.QueryContext(ctx, "SELECT category FROM course_categories WHERE course_id = ?", courseID)
		var catList []string
		if catRows != nil {
			defer catRows.Close()
			for catRows.Next() {
				var cat string
				if err := catRows.Scan(&cat); err == nil {
					catList = append(catList, cat)
				}
			}
		}
		
		courses = append(courses, map[string]interface{}{
			"courseId": courseID, "resourceType": resourceType.String, "name": name.String,
			"teacher": teacherList, "category": catList, "semester": semester.String,
			"credit": credit, "cover": cover.String, "views": views, "loves": loves, "collections": collections,
		})
	}
	
	return courses, nil
}

func (r *courseRepository) UploadResource(ctx context.Context, userID int, courseID string, data map[string]interface{}) (map[string]interface{}, error) {
	var resource1, resource2 map[string]interface{}
	
	// 如果有URL资源，插入到course_resources_web
	if resource := getString(data, "resource"); resource != "" {
		result, err := r.db.ExecContext(ctx, "INSERT INTO course_resources_web (course_id, resource_intro, resource_url, created_at) VALUES (?, ?, ?, NOW())", courseID, getString(data, "description"), resource)
		if err != nil {
			return nil, fmt.Errorf("failed to upload web resource: %w", err)
		}
		resourceID, _ := result.LastInsertId()
		resource1 = map[string]interface{}{"resource_id": resourceID, "resource_intro": getString(data, "description"), "resource_url": resource}
	}
	
	// 如果有上传文件，插入到course_resources_upload
	if file := getString(data, "file"); file != "" {
		result, err := r.db.ExecContext(ctx, "INSERT INTO course_resources_upload (course_id, resource_intro, resource_upload, created_at) VALUES (?, ?, ?, NOW())", courseID, getString(data, "description"), file)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file resource: %w", err)
		}
		resourceID, _ := result.LastInsertId()
		resource2 = map[string]interface{}{"resource_id": resourceID, "resource_intro": getString(data, "description"), "resource_upload": file}
	}
	
	// 添加用户为贡献者
	r.db.ExecContext(ctx, "INSERT INTO course_contributors (course_id, user_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE course_id=course_id", courseID, userID)
	
	// 确定返回的resourceId（优先使用resource1的ID，如果没有则使用resource2的ID）
	var resourceID int64
	if resource1 != nil {
		if id, ok := resource1["resource_id"].(int64); ok {
			resourceID = id
		}
	} else if resource2 != nil {
		if id, ok := resource2["resource_id"].(int64); ok {
			resourceID = id
		}
	}
	
	return map[string]interface{}{
		"resourceId": resourceID, "resourceType": "teach", "resource1": resource1, "resource2": resource2,
		"auditStatus": "pending", "submitTime": time.Now().Format("2006-01-02 15:04:05"), "auditTime": nil, "rejectReason": nil,
	}, nil
}

func (r *courseRepository) DownloadTextbook(ctx context.Context, courseID, textbookID string) (string, error) {
	var uploadURL sql.NullString
	err := r.db.QueryRowContext(ctx, "SELECT resource_upload FROM course_resources_upload WHERE resource_id = ? AND course_id = ?", textbookID, courseID).Scan(&uploadURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("textbook not found")
		}
		return "", fmt.Errorf("failed to get textbook: %w", err)
	}
	return uploadURL.String, nil
}

func (r *courseRepository) AddComment(ctx context.Context, userID int, courseID, content string) (map[string]interface{}, error) {
	var username, avatar sql.NullString
	err := r.db.QueryRowContext(ctx, "SELECT username, avatar FROM users WHERE id = ?", userID).Scan(&username, &avatar)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	result, err := r.db.ExecContext(ctx, "INSERT INTO comments (resource_type, resource_id, user_id, content, created_at, updated_at) VALUES ('course', ?, ?, ?, NOW(), NOW())", courseID, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get comment ID: %w", err)
	}
	return map[string]interface{}{
		"comment_Id": commentID, "nickname": username.String, "avater": avatar.String,
		"comment": content, "commentDate": time.Now().Format("2006-01-02 15:04:05"),
		"love_count": 0, "reply_total": 0, "replies": []interface{}{},
	}, nil
}

func (r *courseRepository) DeleteComment(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error) {
	_, err := r.db.ExecContext(ctx, "UPDATE comments SET deleted_at = NOW() WHERE comment_id = ? AND user_id = ? AND resource_type = 'course' AND resource_id = ?", commentID, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete comment: %w", err)
	}
	return map[string]interface{}{
		"message": "Comment deleted successfully",
	}, nil
}

func (r *courseRepository) GetComments(ctx context.Context, courseID string, cursor, limit int) ([]map[string]interface{}, error) {
	query := `SELECT c.comment_id, c.user_id, c.content, c.love_count, c.reply_total, c.created_at, 
	          COALESCE(u.nickname, u.username) as nickname, 
	          COALESCE(u.avatar, '') as avater 
	          FROM comments c 
	          JOIN users u ON c.user_id = u.id 
	          WHERE c.resource_type = 'course' AND c.resource_id = ? AND c.parent_id IS NULL AND c.deleted_at IS NULL 
	          ORDER BY c.created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, courseID, limit, cursor)
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
		replyRows, _ := r.db.QueryContext(ctx, `SELECT c.comment_id, c.user_id, c.content, c.love_count, c.created_at, 
		                                          COALESCE(u.nickname, u.username) as nickname, 
		                                          COALESCE(u.avatar, '') as avater, c.parent_id 
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
						"comment_Id": replyID, "nickname": replyNickname.String, "avater": replyAvatar.String,
						"comment": replyContent.String, "commentDate": replyCreatedAt.Format("2006-01-02 15:04:05"),
						"love_count": replyLoveCount, "isreply": true, "reply_id": parentID.Int64,
					})
				}
			}
		}
		comments = append(comments, map[string]interface{}{
			"comment_Id": commentID, "nickname": nickname.String, "avater": avatar.String,
			"comment": content.String, "commentDate": createdAt.Format("2006-01-02 15:04:05"),
			"love_count": loveCount, "reply_total": replyTotal, "replies": replies,
			"userId": userID, // 前端可能需要userId来判断是否可以删除
		})
	}
	return comments, nil
}

func (r *courseRepository) LikeComment(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error) {
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&exists)
	isLiked := exists > 0
	
	if isLiked {
		// 取消点赞
		_, err := r.db.ExecContext(ctx, "DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to unlike comment: %w", err)
		}
		r.db.ExecContext(ctx, "UPDATE comments SET love_count = GREATEST(love_count - 1, 0) WHERE comment_id = ?", commentID)
		isLiked = false
	} else {
		// 点赞
		_, err := r.db.ExecContext(ctx, "INSERT INTO comment_likes (comment_id, user_id, created_at) VALUES (?, ?, NOW())", commentID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to like comment: %w", err)
		}
		r.db.ExecContext(ctx, "UPDATE comments SET love_count = love_count + 1 WHERE comment_id = ?", commentID)
		isLiked = true
	}
	
	// 获取更新后的点赞数
	var likes int
	r.db.QueryRowContext(ctx, "SELECT love_count FROM comments WHERE comment_id = ?", commentID).Scan(&likes)
	
	return map[string]interface{}{
		"isliked": isLiked,
		"likes":   likes,
	}, nil
}

func (r *courseRepository) GetResources(ctx context.Context, courseID string) (map[string]interface{}, error) {
	// 查询URL资源（course_resources_web）
	urlRows, err := r.db.QueryContext(ctx, 
		"SELECT resource_id, resource_intro, resource_url FROM course_resources_web WHERE course_id = ? ORDER BY sort_order ASC, created_at ASC",
		courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query web resources: %w", err)
	}
	defer urlRows.Close()
	
	var urlForm []map[string]interface{}
	for urlRows.Next() {
		var resourceID int
		var resourceIntro, resourceURL sql.NullString
		if err := urlRows.Scan(&resourceID, &resourceIntro, &resourceURL); err == nil {
			urlForm = append(urlForm, map[string]interface{}{
				"resource_id":   resourceID,
				"resource_intro": resourceIntro.String,
				"resource_url":   resourceURL.String,
			})
		}
	}
	
	// 查询上传资源（course_resources_upload）
	uploadRows, err := r.db.QueryContext(ctx,
		"SELECT resource_id, resource_intro, resource_upload FROM course_resources_upload WHERE course_id = ? ORDER BY sort_order ASC, created_at ASC",
		courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query upload resources: %w", err)
	}
	defer uploadRows.Close()
	
	var uploadForm []map[string]interface{}
	for uploadRows.Next() {
		var resourceID int
		var resourceIntro, resourceUpload sql.NullString
		if err := uploadRows.Scan(&resourceID, &resourceIntro, &resourceUpload); err == nil {
			uploadForm = append(uploadForm, map[string]interface{}{
				"resource_id":     resourceID,
				"resource_intro":  resourceIntro.String,
				"resource_upload": resourceUpload.String,
			})
		}
	}
	
	return map[string]interface{}{
		"url_form":    urlForm,
		"upload_form": uploadForm,
	}, nil
}

func (r *courseRepository) Create(ctx context.Context, userID int, data map[string]interface{}) (map[string]interface{}, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	result, err := tx.ExecContext(ctx, `INSERT INTO courses (resource_type, name, semester, credit, cover, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`,
		"course", getString(data, "name"), getString(data, "semester"), getInt(data, "credit"), getString(data, "cover"))
	if err != nil {
		return nil, fmt.Errorf("failed to insert course: %w", err)
	}
	
	courseID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get course ID: %w", err)
	}
	
	// 插入教师
	if teachers, ok := data["teacher"].([]string); ok {
		stmt, _ := tx.PrepareContext(ctx, "INSERT INTO course_teachers (course_id, teacher_name) VALUES (?, ?)")
		for _, teacher := range teachers {
			stmt.ExecContext(ctx, courseID, teacher)
		}
		stmt.Close()
	}
	
	// 插入分类
	if categories, ok := data["category"].([]string); ok {
		stmt, _ := tx.PrepareContext(ctx, "INSERT INTO course_categories (course_id, category) VALUES (?, ?)")
		for _, cat := range categories {
			stmt.ExecContext(ctx, courseID, cat)
		}
		stmt.Close()
	}
	
	// 添加提交者为贡献者
	tx.ExecContext(ctx, "INSERT INTO course_contributors (course_id, user_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE course_id=course_id", courseID, userID)
	
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	var username sql.NullString
	r.db.QueryRowContext(ctx, "SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	
	return map[string]interface{}{
		"resourceId": courseID, "resourceType": "course", "name": getString(data, "name"),
		"auditStatus": "pending", "submitTime": time.Now().Format("2006-01-02 15:04:05"),
		"auditTime": nil, "rejectReason": nil, "submitor": username.String,
	}, nil
}

func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
	}
	return 0
}

func (r *courseRepository) ReplyComment(ctx context.Context, userID int, courseID, commentID, content string) (map[string]interface{}, error) {
	var username, avatar sql.NullString
	err := r.db.QueryRowContext(ctx, "SELECT username, avatar FROM users WHERE id = ?", userID).Scan(&username, &avatar)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	result, err := r.db.ExecContext(ctx, "INSERT INTO comments (resource_type, resource_id, parent_id, user_id, content, created_at, updated_at) VALUES ('course', ?, ?, ?, ?, NOW(), NOW())", courseID, commentID, userID, content)
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

func (r *courseRepository) DeleteReply(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error) {
	var parentID sql.NullInt64
	r.db.QueryRowContext(ctx, "SELECT parent_id FROM comments WHERE comment_id = ?", commentID).Scan(&parentID)
	_, err := r.db.ExecContext(ctx, "UPDATE comments SET deleted_at = NOW() WHERE comment_id = ? AND user_id = ? AND resource_type = 'course' AND resource_id = ?", commentID, userID, courseID)
	if err == nil && parentID.Valid {
		r.db.ExecContext(ctx, "UPDATE comments SET reply_total = GREATEST(reply_total - 1, 0) WHERE comment_id = ?", parentID.Int64)
	}
	return map[string]interface{}{"commentId": commentID}, err
}

func (r *courseRepository) AddView(ctx context.Context, courseID string) (int, error) {
	_, err := r.db.ExecContext(ctx, "UPDATE courses SET views = views + 1 WHERE course_id = ?", courseID)
	if err != nil {
		return 0, err
	}
	var views int
	r.db.QueryRowContext(ctx, "SELECT views FROM courses WHERE course_id = ?", courseID).Scan(&views)
	return views, nil
}

func (r *courseRepository) CollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	// 先检查是否已收藏
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM collections WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?", userID, courseID).Scan(&exists)
	
	if exists > 0 {
		// 已收藏，直接返回
		var collections int
		r.db.QueryRowContext(ctx, "SELECT collections FROM courses WHERE course_id = ?", courseID).Scan(&collections)
		return map[string]interface{}{"iscollected": true, "collections": collections}, nil
	}
	
	// 未收藏，执行插入
	_, err := r.db.ExecContext(ctx, "INSERT INTO collections (user_id, resource_type, resource_id, created_at) VALUES (?, 'course', ?, NOW())", userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to collect course: %w", err)
	}
	
	// 增加收藏数
	_, err = r.db.ExecContext(ctx, "UPDATE courses SET collections = collections + 1 WHERE course_id = ?", courseID)
	var collections int
	r.db.QueryRowContext(ctx, "SELECT collections FROM courses WHERE course_id = ?", courseID).Scan(&collections)
	return map[string]interface{}{"iscollected": true, "collections": collections}, err
}

func (r *courseRepository) UncollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	// 先检查是否已收藏，只有已收藏才执行删除和减少计数
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM collections WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?", userID, courseID).Scan(&exists)
	
	if exists == 0 {
		// 未收藏，直接返回当前状态
		var collections int
		r.db.QueryRowContext(ctx, "SELECT collections FROM courses WHERE course_id = ?", courseID).Scan(&collections)
		return map[string]interface{}{"iscollected": false, "collections": collections}, nil
	}
	
	// 已收藏，执行删除
	result, err := r.db.ExecContext(ctx, "DELETE FROM collections WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?", userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to uncollect course: %w", err)
	}
	
	// 检查删除是否成功（受影响的行数）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	
	// 只有成功删除才减少计数
	if rowsAffected > 0 {
		_, err = r.db.ExecContext(ctx, "UPDATE courses SET collections = GREATEST(collections - 1, 0) WHERE course_id = ?", courseID)
	}
	
	var collections int
	r.db.QueryRowContext(ctx, "SELECT collections FROM courses WHERE course_id = ?", courseID).Scan(&collections)
	return map[string]interface{}{"iscollected": false, "collections": collections}, err
}

func (r *courseRepository) LikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	// 先检查是否已点赞
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM likes WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?", userID, courseID).Scan(&exists)
	
	if exists > 0 {
		// 已点赞，直接返回（或者可以考虑返回错误，取决于业务逻辑）
		var likes int
		r.db.QueryRowContext(ctx, "SELECT loves FROM courses WHERE course_id = ?", courseID).Scan(&likes)
		return map[string]interface{}{"isliked": true, "likes": likes}, nil
	}
	
	// 未点赞，执行插入
	_, err := r.db.ExecContext(ctx, "INSERT INTO likes (user_id, resource_type, resource_id, created_at) VALUES (?, 'course', ?, NOW())", userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to like course: %w", err)
	}
	
	// 增加点赞数
	_, err = r.db.ExecContext(ctx, "UPDATE courses SET loves = loves + 1 WHERE course_id = ?", courseID)
	var likes int
	r.db.QueryRowContext(ctx, "SELECT loves FROM courses WHERE course_id = ?", courseID).Scan(&likes)
	return map[string]interface{}{"isliked": true, "likes": likes}, err
}

func (r *courseRepository) UnlikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	// 先检查是否已点赞，只有已点赞才执行删除和减少计数
	var exists int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM likes WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?", userID, courseID).Scan(&exists)
	
	if exists == 0 {
		// 未点赞，直接返回当前状态
		var likes int
		r.db.QueryRowContext(ctx, "SELECT loves FROM courses WHERE course_id = ?", courseID).Scan(&likes)
		return map[string]interface{}{"isliked": false, "likes": likes}, nil
	}
	
	// 已点赞，执行删除
	result, err := r.db.ExecContext(ctx, "DELETE FROM likes WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?", userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to unlike course: %w", err)
	}
	
	// 检查删除是否成功（受影响的行数）
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	
	// 只有成功删除才减少计数
	if rowsAffected > 0 {
		_, err = r.db.ExecContext(ctx, "UPDATE courses SET loves = GREATEST(loves - 1, 0) WHERE course_id = ?", courseID)
	}
	
	var likes int
	r.db.QueryRowContext(ctx, "SELECT loves FROM courses WHERE course_id = ?", courseID).Scan(&likes)
	return map[string]interface{}{"isliked": false, "likes": likes}, err
}

func (r *courseRepository) GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error) {
	query := `SELECT c.course_id, c.name, c.created_at, u.username FROM courses c LEFT JOIN users u ON EXISTS (SELECT 1 FROM course_contributors cc WHERE cc.course_id = c.course_id AND cc.user_id = u.id LIMIT 1) WHERE c.course_id IN (SELECT DISTINCT cc2.course_id FROM course_contributors cc2) ORDER BY c.created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending courses: %w", err)
	}
	defer rows.Close()
	var courses []map[string]interface{}
	for rows.Next() {
		var courseID int
		var name, username sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&courseID, &name, &createdAt, &username); err != nil {
			continue
		}
		var catagory sql.NullString
		r.db.QueryRowContext(ctx, "SELECT category FROM course_categories WHERE course_id = ? LIMIT 1", courseID).Scan(&catagory)
		courses = append(courses, map[string]interface{}{
			"id":           courseID,
			"resourceId":   courseID,
			"reourceId":    courseID, // 保持兼容性
			"courseId":     courseID,
			"courseName":   name.String,
			"name":         name.String, // 前端可能使用name
			"title":        name.String, // 前端可能使用title
			"resourcename": name.String, // 保持兼容性
			"type":         "doc", // 默认类型，可以从资源表获取
			"link":         "", // 可以从资源表获取
			"description":  "",
			"category":     catagory.String,
			"catagory":     catagory.String, // 保持兼容性
			"uploader":     username.String,
			"submitor":     username.String, // 保持兼容性
			"author":       username.String, // 前端可能使用author
			"owner":         username.String, // 前端可能使用owner
			"created_at":   createdAt.Format("2006-01-02T15:04:05Z"),
			"created":      createdAt.Format("2006-01-02T15:04:05Z"), // 前端可能使用created
			"submitDate":   createdAt.Format("2006-01-02 15:04:05"), // 保持兼容性
		})
	}
	return courses, nil
}

// CheckUserLike 检查用户是否点赞了课程
func (r *courseRepository) CheckUserLike(ctx context.Context, userID int, courseID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM likes WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?", userID, courseID).Scan(&count)
	return count > 0, err
}

// CheckUserCollect 检查用户是否收藏了课程
func (r *courseRepository) CheckUserCollect(ctx context.Context, userID int, courseID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM collections WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?", userID, courseID).Scan(&count)
	return count > 0, err
}

// UpdateCourseStatus 更新课程审核状态
// 注意：如果 courses 表没有 status 字段，此方法会尝试添加字段或使用 resource_status_logs 表
func (r *courseRepository) UpdateCourseStatus(ctx context.Context, courseID, status, rejectReason string) error {
	// 首先检查 courses 表是否有 status 字段
	var hasStatusField bool
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) > 0 
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = 'courses' 
		AND COLUMN_NAME = 'status'
	`).Scan(&hasStatusField)
	
	if err != nil {
		// 如果查询失败，假设字段不存在，使用 resource_status_logs 表
		hasStatusField = false
	}
	
	if hasStatusField {
		// 如果 courses 表有 status 字段，直接更新
		query := "UPDATE courses SET status = ?, audit_time = NOW()"
		args := []interface{}{status}
		
		if rejectReason != "" {
			// 检查是否有 reject_reason 字段
			var hasRejectReasonField bool
			r.db.QueryRowContext(ctx, `
				SELECT COUNT(*) > 0 
				FROM information_schema.COLUMNS 
				WHERE TABLE_SCHEMA = DATABASE() 
				AND TABLE_NAME = 'courses' 
				AND COLUMN_NAME = 'reject_reason'
			`).Scan(&hasRejectReasonField)
			
			if hasRejectReasonField {
				query += ", reject_reason = ?"
				args = append(args, rejectReason)
			}
		}
		
		query += " WHERE course_id = ?"
		args = append(args, courseID)
		
		_, err := r.db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to update course status: %w", err)
		}
	}
	
	// 无论是否有 status 字段，都记录到 resource_status_logs 表
	// 获取旧状态
	var oldStatus string
	r.db.QueryRowContext(ctx, `
		SELECT new_status 
		FROM resource_status_logs 
		WHERE resource_type = 'course' AND resource_id = ? 
		ORDER BY operate_time DESC 
		LIMIT 1
	`, courseID).Scan(&oldStatus)
	
	if oldStatus == "" {
		oldStatus = "pending" // 默认旧状态
	}
	
	// 记录状态变更日志
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO resource_status_logs (resource_type, resource_id, old_status, new_status, operator_id, operate_time, reject_reason)
		VALUES ('course', ?, ?, ?, 0, NOW(), ?)
	`, courseID, oldStatus, status, rejectReason)
	
	if err != nil {
		return fmt.Errorf("failed to log status change: %w", err)
	}
	
	return nil
}
