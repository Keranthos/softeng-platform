package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CourseRepository interface {
	GetCourses(ctx context.Context, semester string, category []string, sort string, limit, cursor int) ([]map[string]interface{}, error)
	GetByID(ctx context.Context, courseID string, userID int) (map[string]interface{}, error)
	Search(ctx context.Context, keyword string, category []string, limit, cursor int) ([]map[string]interface{}, error)
	UploadResource(ctx context.Context, userID int, courseID string, data map[string]interface{}) (map[string]interface{}, error)
	DownloadTextbook(ctx context.Context, courseID, textbookID string) (string, error)
	AddComment(ctx context.Context, userID int, courseID, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, courseID string, commentID string) (map[string]interface{}, error) // commentID 为空则删除该用户最新一条
	ReplyComment(ctx context.Context, userID int, courseID, commentID, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error)
	AddView(ctx context.Context, courseID string) (int, error)
	CollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	UncollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	LikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	UnlikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error)
	GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error) // 新增方法
}

type courseRepository struct {
	db *Database
}

func NewCourseRepository(db *Database) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) GetCourses(ctx context.Context, semester string, category []string, sort string, limit, cursor int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	var (
		whereParts []string
		args       []interface{}
	)

	if semester != "" {
		whereParts = append(whereParts, "c.semester = ?")
		args = append(args, semester)
	}

	if len(category) > 0 {
		whereParts = append(whereParts, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM course_categories cc2 WHERE cc2.course_id = c.course_id AND cc2.category IN (%s))",
			placeholders(len(category)),
		))
		for _, cat := range category {
			args = append(args, cat)
		}
	}

	if cursor > 0 {
		whereParts = append(whereParts, "c.course_id < ?")
		args = append(args, cursor)
	}

	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = "WHERE " + strings.Join(whereParts, " AND ")
	}

	orderBy := "c.created_at DESC, c.course_id DESC"
	switch strings.ToLower(sort) {
	case "views":
		orderBy = "c.views DESC, c.course_id DESC"
	case "collections":
		orderBy = "c.collections DESC, c.course_id DESC"
	case "loves", "likes":
		orderBy = "c.loves DESC, c.course_id DESC"
	}

	query := fmt.Sprintf(`
		SELECT
			c.course_id,
			c.resource_type,
			c.name,
			COALESCE(GROUP_CONCAT(DISTINCT ct.teacher_name SEPARATOR ','), '') AS teachers_csv,
			COALESCE(GROUP_CONCAT(DISTINCT cc.category SEPARATOR ','), '') AS categories_csv,
			c.semester,
			COALESCE(c.credit, 0) AS credit,
			COALESCE(c.cover, '') AS cover,
			c.views,
			c.loves,
			c.collections,
			c.created_at
		FROM courses c
		LEFT JOIN course_teachers ct ON ct.course_id = c.course_id
		LEFT JOIN course_categories cc ON cc.course_id = c.course_id
		%s
		GROUP BY c.course_id, c.resource_type, c.name, c.semester, c.credit, c.cover, c.views, c.loves, c.collections, c.created_at
		ORDER BY %s
		LIMIT ?
	`, whereSQL, orderBy)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query courses: %v", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var (
			courseID     int
			resourceType sql.NullString
			name         sql.NullString
			teachersCSV  sql.NullString
			categoriesCSV sql.NullString
			semesterNS   sql.NullString
			credit       int
			cover        sql.NullString
			views        int
			loves        int
			collections  int
			createdAt    time.Time
		)

		if err := rows.Scan(
			&courseID,
			&resourceType,
			&name,
			&teachersCSV,
			&categoriesCSV,
			&semesterNS,
			&credit,
			&cover,
			&views,
			&loves,
			&collections,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan course row: %v", err)
		}

		result = append(result, map[string]interface{}{
			"courseId":     courseID,
			"resourceType": nullString(resourceType),
			"name":         nullString(name),
			"teacher":      splitCSV(nullString(teachersCSV)),
			"category":     splitCSV(nullString(categoriesCSV)),
			"semester":     nullString(semesterNS),
			"credit":       credit,
			"cover":        nullString(cover),
			"views":        views,
			"loves":        loves,
			"collections":  collections,
			"createdat":    createdAt.Format("2006-01-02"),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate course rows: %v", err)
	}

	return result, nil
}

func (r *courseRepository) GetByID(ctx context.Context, courseID string, userID int) (map[string]interface{}, error) {
	query := `
		SELECT
			course_id,
			resource_type,
			name,
			semester,
			COALESCE(credit, 0) AS credit,
			cover,
			views,
			loves,
			collections,
			created_at
		FROM courses
		WHERE course_id = ?
		LIMIT 1
	`

	var (
		id           int
		resourceType sql.NullString
		name         sql.NullString
		semesterNS   sql.NullString
		credit       int
		cover        sql.NullString
		views        int
		loves        int
		collections  int
		createdAt    time.Time
	)

	err := r.db.QueryRowContext(ctx, query, courseID).Scan(
		&id,
		&resourceType,
		&name,
		&semesterNS,
		&credit,
		&cover,
		&views,
		&loves,
		&collections,
		&createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get course by id: %v", err)
	}

	categories, err := r.fetchCourseCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	teachers, err := r.fetchCourseTeachers(ctx, id)
	if err != nil {
		return nil, err
	}
	urlForm, err := r.fetchCourseWebResources(ctx, id)
	if err != nil {
		return nil, err
	}
	uploadForm, err := r.fetchCourseUploadResources(ctx, id)
	if err != nil {
		return nil, err
	}
	contributors, err := r.fetchCourseContributors(ctx, id)
	if err != nil {
		return nil, err
	}

	catagory := ""
	if len(categories) > 0 {
		catagory = categories[0]
	}

	// 登录态：返回 isliked/iscollected
	isLiked := false
	isCollected := false
	if userID > 0 {
		if v, err := r.isCourseLiked(ctx, userID, id); err == nil {
			isLiked = v
		} else {
			return nil, err
		}
		if v, err := r.isCourseCollected(ctx, userID, id); err == nil {
			isCollected = v
		} else {
			return nil, err
		}
	}

	comments, commentTotal, err := r.fetchCourseComments(ctx, id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"courseId":     id,
		"resourceType": nullString(resourceType),
		"name":         nullString(name),
		"catagory":     catagory,
		"url_form":     urlForm,
		"upload_form":  uploadForm,
		"teacher":      teachers,
		"semester":     nullString(semesterNS),
		"credit":       credit,
		"cover":        nullString(cover),
		"contributor":  contributors,
		"collections":  collections,
		"views":        views,
		"likes":        loves, // schema 里 courses.loves，这里按接口习惯返回 likes
		"isliked":      isLiked,
		"iscollected":  isCollected,
		"comment_total": commentTotal,
		"comments":     comments,
		"createdAt":    createdAt.Format("2006-01-02"),
	}, nil
}

func (r *courseRepository) Search(ctx context.Context, keyword string, category []string, limit, cursor int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	var (
		whereParts []string
		args       []interface{}
	)

	if keyword != "" {
		whereParts = append(whereParts, "c.name LIKE ?")
		args = append(args, "%"+keyword+"%")
	}

	if len(category) > 0 {
		whereParts = append(whereParts, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM course_categories cc2 WHERE cc2.course_id = c.course_id AND cc2.category IN (%s))",
			placeholders(len(category)),
		))
		for _, cat := range category {
			args = append(args, cat)
		}
	}

	if cursor > 0 {
		whereParts = append(whereParts, "c.course_id < ?")
		args = append(args, cursor)
	}

	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = "WHERE " + strings.Join(whereParts, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT
			c.course_id,
			c.resource_type,
			c.name,
			COALESCE(GROUP_CONCAT(DISTINCT ct.teacher_name SEPARATOR ','), '') AS teachers_csv,
			COALESCE(GROUP_CONCAT(DISTINCT cc.category SEPARATOR ','), '') AS categories_csv,
			c.semester,
			COALESCE(c.credit, 0) AS credit,
			COALESCE(c.cover, '') AS cover,
			c.views,
			c.loves,
			c.collections,
			c.created_at
		FROM courses c
		LEFT JOIN course_teachers ct ON ct.course_id = c.course_id
		LEFT JOIN course_categories cc ON cc.course_id = c.course_id
		%s
		GROUP BY c.course_id, c.resource_type, c.name, c.semester, c.credit, c.cover, c.views, c.loves, c.collections, c.created_at
		ORDER BY c.created_at DESC, c.course_id DESC
		LIMIT ?
	`, whereSQL)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search courses: %v", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var (
			courseID     int
			resourceType sql.NullString
			name         sql.NullString
			teachersCSV  sql.NullString
			categoriesCSV sql.NullString
			semesterNS   sql.NullString
			credit       int
			cover        sql.NullString
			views        int
			loves        int
			collections  int
			createdAt    time.Time
		)

		if err := rows.Scan(
			&courseID,
			&resourceType,
			&name,
			&teachersCSV,
			&categoriesCSV,
			&semesterNS,
			&credit,
			&cover,
			&views,
			&loves,
			&collections,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan course search row: %v", err)
		}

		result = append(result, map[string]interface{}{
			"courseId":     courseID,
			"resourceType": nullString(resourceType),
			"name":         nullString(name),
			"teacher":      splitCSV(nullString(teachersCSV)),
			"category":     splitCSV(nullString(categoriesCSV)),
			"semester":     nullString(semesterNS),
			"credit":       credit,
			"cover":        nullString(cover),
			"views":        views,
			"loves":        loves,
			"collections":  collections,
			"createdat":    createdAt.Format("2006-01-02"),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate course search rows: %v", err)
	}

	return result, nil
}

func (r *courseRepository) UploadResource(ctx context.Context, userID int, courseID string, data map[string]interface{}) (map[string]interface{}, error) {
	// 实现上传课程资源的逻辑
	return map[string]interface{}{
		"resourceId":   1,
		"resourceType": "teach",
		"resource1": map[string]interface{}{
			"resource_intro": data["description"],
			"resource_url":   data["resource"],
			"resource_id":    1,
		},
		"resource2": map[string]interface{}{
			"resource_intro":  data["description"],
			"resource_upload": data["file"],
			"resource_id":     2,
		},
		"auditStatus":  "pending",
		"submitTime":   "2023-12-01 10:00:00",
		"auditTime":    nil,
		"rejectReason": nil,
	}, nil
}

func (r *courseRepository) DownloadTextbook(ctx context.Context, courseID, textbookID string) (string, error) {
	// 实现下载课本的逻辑
	return "Textbook content for course " + courseID + " textbook " + textbookID, nil
}

func (r *courseRepository) AddComment(ctx context.Context, userID int, courseID, content string) (map[string]interface{}, error) {
	cid, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course id")
	}

	res, err := r.db.ExecContext(ctx, `
		INSERT INTO comments (resource_type, resource_id, parent_id, user_id, content)
		VALUES ('course', ?, NULL, ?, ?)
	`, cid, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to insert comment: %v", err)
	}
	commentID64, _ := res.LastInsertId()
	commentID := int(commentID64)

	nickname, avatar, err := fetchUserDisplayCourse(ctx, r.db, userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"comment_Id":  commentID,
		"nickname":    nickname,
		"avater":      avatar,
		"comment":     content,
		"commentDate": time.Now().Format("2006-01-02 15:04:05"),
		"love_count":  0,
		"reply_total": 0,
		"replies":     []map[string]interface{}{},
	}, nil
}

func (r *courseRepository) DeleteComment(ctx context.Context, userID int, courseID string, commentID string) (map[string]interface{}, error) {
	cid, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course id")
	}

	var (
		targetID int
		content  sql.NullString
	)

	if strings.TrimSpace(commentID) != "" {
		targetID, err = strconv.Atoi(commentID)
		if err != nil {
			return nil, fmt.Errorf("invalid comment id")
		}
		err = r.db.QueryRowContext(ctx, `
			SELECT comment_id, content
			FROM comments
			WHERE comment_id = ? AND resource_type = 'course' AND resource_id = ? AND user_id = ? AND parent_id IS NULL AND deleted_at IS NULL
			LIMIT 1
		`, targetID, cid, userID).Scan(&targetID, &content)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("comment not found")
			}
			return nil, fmt.Errorf("failed to read comment: %v", err)
		}
	} else {
		err = r.db.QueryRowContext(ctx, `
			SELECT comment_id, content
			FROM comments
			WHERE resource_type = 'course' AND resource_id = ? AND user_id = ? AND parent_id IS NULL AND deleted_at IS NULL
			ORDER BY created_at DESC, comment_id DESC
			LIMIT 1
		`, cid, userID).Scan(&targetID, &content)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("comment not found")
			}
			return nil, fmt.Errorf("failed to read latest comment: %v", err)
		}
	}

	if _, err := r.db.ExecContext(ctx, `
		UPDATE comments
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE comment_id = ? AND user_id = ? AND resource_type = 'course' AND resource_id = ? AND deleted_at IS NULL
	`, targetID, userID, cid); err != nil {
		return nil, fmt.Errorf("failed to delete comment: %v", err)
	}

	nickname, avatar, err := fetchUserDisplayCourse(ctx, r.db, userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"comment_Id":  targetID,
		"nickname":    nickname,
		"avater":      avatar,
		"comment":     "已删除的评论",
		"delete_Date": time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *courseRepository) ReplyComment(ctx context.Context, userID int, courseID, commentID, content string) (map[string]interface{}, error) {
	cid, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course id")
	}
	parentID, err := strconv.Atoi(commentID)
	if err != nil {
		return nil, fmt.Errorf("invalid comment id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	var exists int
	if err := tx.QueryRowContext(ctx, `
		SELECT 1
		FROM comments
		WHERE comment_id = ? AND resource_type = 'course' AND resource_id = ? AND parent_id IS NULL AND deleted_at IS NULL
		LIMIT 1
	`, parentID, cid).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("parent comment not found")
		}
		return nil, fmt.Errorf("failed to check parent comment: %v", err)
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO comments (resource_type, resource_id, parent_id, user_id, content)
		VALUES ('course', ?, ?, ?, ?)
	`, cid, parentID, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to insert reply: %v", err)
	}
	replyID64, _ := res.LastInsertId()
	replyID := int(replyID64)

	if _, err := tx.ExecContext(ctx, `UPDATE comments SET reply_total = reply_total + 1 WHERE comment_id = ?`, parentID); err != nil {
		return nil, fmt.Errorf("failed to update reply_total: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	nickname, avatar, err := fetchUserDisplayCourse(ctx, r.db, userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"comment_Id":  replyID,
		"nickname":    nickname,
		"avater":      avatar,
		"comment":     content,
		"commentDate": time.Now().Format("2006-01-02 15:04:05"),
		"isreply":     true,
		"reply_id":    parentID,
		"replies":     []map[string]interface{}{},
	}, nil
}

func (r *courseRepository) DeleteReply(ctx context.Context, userID int, courseID, commentID string) (map[string]interface{}, error) {
	cid, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course id")
	}
	rid, err := strconv.Atoi(commentID)
	if err != nil {
		return nil, fmt.Errorf("invalid reply id")
	}

	var parentID sql.NullInt64
	err = r.db.QueryRowContext(ctx, `
		SELECT parent_id
		FROM comments
		WHERE comment_id = ? AND resource_type = 'course' AND resource_id = ? AND user_id = ? AND deleted_at IS NULL
		LIMIT 1
	`, rid, cid, userID).Scan(&parentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reply not found")
		}
		return nil, fmt.Errorf("failed to read reply: %v", err)
	}
	if !parentID.Valid {
		return nil, fmt.Errorf("not a reply")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `
		UPDATE comments
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE comment_id = ? AND user_id = ? AND resource_type = 'course' AND resource_id = ? AND deleted_at IS NULL
	`, rid, userID, cid); err != nil {
		return nil, fmt.Errorf("failed to delete reply: %v", err)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE comments
		SET reply_total = GREATEST(reply_total - 1, 0), updated_at = NOW()
		WHERE comment_id = ?
	`, parentID.Int64); err != nil {
		return nil, fmt.Errorf("failed to decrement parent reply_total: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	nickname, avatar, err := fetchUserDisplayCourse(ctx, r.db, userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"comment_Id":  rid,
		"nickname":    nickname,
		"avater":      avatar,
		"comment":     "已删除的回复",
		"delete_Date": time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *courseRepository) AddView(ctx context.Context, courseID string) (int, error) {
	cid, err := strconv.Atoi(courseID)
	if err != nil {
		return 0, fmt.Errorf("invalid course id")
	}

	if _, err := r.db.ExecContext(ctx, `UPDATE courses SET views = views + 1 WHERE course_id = ?`, cid); err != nil {
		return 0, fmt.Errorf("failed to update views: %v", err)
	}

	var views int
	if err := r.db.QueryRowContext(ctx, `SELECT views FROM courses WHERE course_id = ?`, cid).Scan(&views); err != nil {
		return 0, fmt.Errorf("failed to read views: %v", err)
	}
	return views, nil
}

func (r *courseRepository) CollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	cid, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT IGNORE INTO collections (user_id, resource_type, resource_id)
		VALUES (?, 'course', ?)
	`, userID, cid)
	if err != nil {
		return nil, fmt.Errorf("failed to insert collection: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE courses SET collections = collections + 1 WHERE course_id = ?`, cid); err != nil {
			return nil, fmt.Errorf("failed to update collections: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	var collections int
	if err := r.db.QueryRowContext(ctx, `SELECT collections FROM courses WHERE course_id = ?`, cid).Scan(&collections); err != nil {
		return nil, fmt.Errorf("failed to read collections: %v", err)
	}

	return map[string]interface{}{
		"iscollected": true,
		"collections": collections,
	}, nil
}

func (r *courseRepository) UncollectCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	cid, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM collections
		WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?
	`, userID, cid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete collection: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE courses SET collections = GREATEST(collections - 1, 0) WHERE course_id = ?`, cid); err != nil {
			return nil, fmt.Errorf("failed to decrement collections: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	var collections int
	if err := r.db.QueryRowContext(ctx, `SELECT collections FROM courses WHERE course_id = ?`, cid).Scan(&collections); err != nil {
		return nil, fmt.Errorf("failed to read collections: %v", err)
	}

	return map[string]interface{}{
		"iscollected": false,
		"collections": collections,
	}, nil
}

func (r *courseRepository) LikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	cid, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT IGNORE INTO likes (user_id, resource_type, resource_id)
		VALUES (?, 'course', ?)
	`, userID, cid)
	if err != nil {
		return nil, fmt.Errorf("failed to insert like: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE courses SET loves = loves + 1 WHERE course_id = ?`, cid); err != nil {
			return nil, fmt.Errorf("failed to update loves: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	var loves int
	if err := r.db.QueryRowContext(ctx, `SELECT loves FROM courses WHERE course_id = ?`, cid).Scan(&loves); err != nil {
		return nil, fmt.Errorf("failed to read loves: %v", err)
	}

	return map[string]interface{}{
		"isliked": true,
		"likes":   loves,
	}, nil
}

func (r *courseRepository) UnlikeCourse(ctx context.Context, userID int, courseID string) (map[string]interface{}, error) {
	cid, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM likes
		WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?
	`, userID, cid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete like: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE courses SET loves = GREATEST(loves - 1, 0) WHERE course_id = ?`, cid); err != nil {
			return nil, fmt.Errorf("failed to decrement loves: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	var loves int
	if err := r.db.QueryRowContext(ctx, `SELECT loves FROM courses WHERE course_id = ?`, cid).Scan(&loves); err != nil {
		return nil, fmt.Errorf("failed to read loves: %v", err)
	}

	return map[string]interface{}{
		"isliked": false,
		"likes":   loves,
	}, nil
}

func fetchUserDisplayCourse(ctx context.Context, db *Database, userID int) (string, string, error) {
	var nickname sql.NullString
	var username sql.NullString
	var avatar sql.NullString
	err := db.QueryRowContext(ctx, `SELECT nickname, username, avatar FROM users WHERE id = ? LIMIT 1`, userID).Scan(&nickname, &username, &avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("user not found")
		}
		return "", "", fmt.Errorf("failed to query user: %v", err)
	}
	display := nullString(nickname)
	if strings.TrimSpace(display) == "" {
		display = nullString(username)
	}
	return display, nullString(avatar), nil
}

func (r *courseRepository) isCourseLiked(ctx context.Context, userID int, courseID int) (bool, error) {
	var one int
	err := r.db.QueryRowContext(ctx, `
		SELECT 1
		FROM likes
		WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?
		LIMIT 1
	`, userID, courseID).Scan(&one)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check isliked: %v", err)
	}
	return true, nil
}

func (r *courseRepository) isCourseCollected(ctx context.Context, userID int, courseID int) (bool, error) {
	var one int
	err := r.db.QueryRowContext(ctx, `
		SELECT 1
		FROM collections
		WHERE user_id = ? AND resource_type = 'course' AND resource_id = ?
		LIMIT 1
	`, userID, courseID).Scan(&one)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check iscollected: %v", err)
	}
	return true, nil
}

type courseCommentRow struct {
	ID         int
	ParentID   sql.NullInt64
	UserID     int
	Nickname   sql.NullString
	Username   sql.NullString
	Avatar     sql.NullString
	Content    sql.NullString
	LoveCount  int
	ReplyTotal int
	CreatedAt  time.Time
}

func (r *courseRepository) fetchCourseComments(ctx context.Context, courseID int) ([]map[string]interface{}, int, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			c.comment_id,
			c.parent_id,
			c.user_id,
			u.nickname,
			u.username,
			u.avatar,
			c.content,
			c.love_count,
			c.reply_total,
			c.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.resource_type = 'course' AND c.resource_id = ? AND c.deleted_at IS NULL
		ORDER BY c.created_at ASC, c.comment_id ASC
	`, courseID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query comments: %v", err)
	}
	defer rows.Close()

	var all []courseCommentRow
	for rows.Next() {
		var row courseCommentRow
		if err := rows.Scan(
			&row.ID,
			&row.ParentID,
			&row.UserID,
			&row.Nickname,
			&row.Username,
			&row.Avatar,
			&row.Content,
			&row.LoveCount,
			&row.ReplyTotal,
			&row.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan comment row: %v", err)
		}
		all = append(all, row)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate comment rows: %v", err)
	}

	top := make([]courseCommentRow, 0)
	repliesByParent := make(map[int][]courseCommentRow)
	for _, c := range all {
		if !c.ParentID.Valid {
			top = append(top, c)
			continue
		}
		pid := int(c.ParentID.Int64)
		repliesByParent[pid] = append(repliesByParent[pid], c)
	}

	commentTotal := len(top)
	out := make([]map[string]interface{}, 0, len(top))
	for _, c := range top {
		nickname := nullString(c.Nickname)
		if strings.TrimSpace(nickname) == "" {
			nickname = nullString(c.Username)
		}

		item := map[string]interface{}{
			"comment_Id":  c.ID,
			"nickname":    nickname,
			"avater":      nullString(c.Avatar),
			"comment":     nullString(c.Content),
			"commentDate": c.CreatedAt.Format("2006-01-02 15:04:05"),
			"love_count":  c.LoveCount,
			"reply_total": c.ReplyTotal,
			"replies":     []map[string]interface{}{},
		}

		replies := repliesByParent[c.ID]
		replyOut := make([]map[string]interface{}, 0, len(replies))
		for _, rpl := range replies {
			rnick := nullString(rpl.Nickname)
			if strings.TrimSpace(rnick) == "" {
				rnick = nullString(rpl.Username)
			}
			replyOut = append(replyOut, map[string]interface{}{
				"comment_Id":  rpl.ID,
				"nickname":    rnick,
				"avater":      nullString(rpl.Avatar),
				"comment":     nullString(rpl.Content),
				"commentDate": rpl.CreatedAt.Format("2006-01-02 15:04:05"),
				"love_count":  rpl.LoveCount,
				"isreply":     true,
				"reply_id":    c.ID,
			})
		}
		item["replies"] = replyOut
		out = append(out, item)
	}

	return out, commentTotal, nil
}

func (r *courseRepository) GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{
			"submitor":     "user1",
			"submitDate":   "2023-12-01 10:00:00",
			"reourceId":    1,
			"resourceType": "course",
			"resourcename": "新课程资源",
			"catagory":     "教学资料",
			"link":         "https://example.com/resource",
			"description":  "课程相关资源",
			"tags":         []string{"讲义", "视频"},
			"file":         "resource.pdf",
		},
	}, nil
}

func (r *courseRepository) fetchCourseTeachers(ctx context.Context, courseID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT teacher_name FROM course_teachers WHERE course_id = ? ORDER BY teacher_name ASC`, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query course teachers: %v", err)
	}
	defer rows.Close()

	var teachers []string
	for rows.Next() {
		var t sql.NullString
		if err := rows.Scan(&t); err != nil {
			return nil, fmt.Errorf("failed to scan teacher: %v", err)
		}
		if s := nullString(t); s != "" {
			teachers = append(teachers, s)
		}
	}
	return teachers, rows.Err()
}

func (r *courseRepository) fetchCourseCategories(ctx context.Context, courseID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT category FROM course_categories WHERE course_id = ? ORDER BY category ASC`, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query course categories: %v", err)
	}
	defer rows.Close()

	var cats []string
	for rows.Next() {
		var c sql.NullString
		if err := rows.Scan(&c); err != nil {
			return nil, fmt.Errorf("failed to scan category: %v", err)
		}
		if s := nullString(c); s != "" {
			cats = append(cats, s)
		}
	}
	return cats, rows.Err()
}

func (r *courseRepository) fetchCourseWebResources(ctx context.Context, courseID int) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT resource_id, resource_intro, resource_url
		FROM course_resources_web
		WHERE course_id = ?
		ORDER BY sort_order ASC, resource_id ASC
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query course web resources: %v", err)
	}
	defer rows.Close()

	var res []map[string]interface{}
	for rows.Next() {
		var (
			id    int
			intro sql.NullString
			url   sql.NullString
		)
		if err := rows.Scan(&id, &intro, &url); err != nil {
			return nil, fmt.Errorf("failed to scan course web resource: %v", err)
		}
		res = append(res, map[string]interface{}{
			"resource_intro": nullString(intro),
			"resource_url":   nullString(url),
			"resource_id":    id,
		})
	}
	return res, rows.Err()
}

func (r *courseRepository) fetchCourseUploadResources(ctx context.Context, courseID int) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT resource_id, resource_intro, resource_upload
		FROM course_resources_upload
		WHERE course_id = ?
		ORDER BY sort_order ASC, resource_id ASC
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query course upload resources: %v", err)
	}
	defer rows.Close()

	var res []map[string]interface{}
	for rows.Next() {
		var (
			id     int
			intro  sql.NullString
			upload sql.NullString
		)
		if err := rows.Scan(&id, &intro, &upload); err != nil {
			return nil, fmt.Errorf("failed to scan course upload resource: %v", err)
		}
		res = append(res, map[string]interface{}{
			"resource_intro":  nullString(intro),
			"resource_upload": nullString(upload),
			"resource_id":     id,
		})
	}
	return res, rows.Err()
}

func (r *courseRepository) fetchCourseContributors(ctx context.Context, courseID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.username
		FROM course_contributors cc
		JOIN users u ON u.id = cc.user_id
		WHERE cc.course_id = ?
		ORDER BY u.username ASC
	`, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query course contributors: %v", err)
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username sql.NullString
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("failed to scan contributor: %v", err)
		}
		if s := nullString(username); s != "" {
			users = append(users, s)
		}
	}
	return users, rows.Err()
}
