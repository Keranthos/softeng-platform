package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ToolRepository interface {
	GetTools(ctx context.Context, category, tags []string, sort, cursor string, pageSize int) ([]map[string]interface{}, error)
	GetByID(ctx context.Context, resourceID string, userID int) (map[string]interface{}, error)
	Search(ctx context.Context, keyword, cursor string, pageSize int) ([]map[string]interface{}, error)
	Create(ctx context.Context, userID int, data map[string]interface{}) (map[string]interface{}, error)

	// 点赞
	AddLike(ctx context.Context, userID int, resourceID string) error
	RemoveLike(ctx context.Context, userID int, resourceID string) error
	GetLikes(ctx context.Context, resourceID string) (int, error)

	// 收藏
	AddCollection(ctx context.Context, userID int, resourceID string) error
	RemoveCollection(ctx context.Context, userID int, resourceID string) error
	GetCollections(ctx context.Context, resourceID string) (int, error)

	// 评论/回复（comments 表）
	AddComment(ctx context.Context, userID int, resourceID string, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, resourceID string, commentID string) (map[string]interface{}, error) // commentID 为空则删除该用户最新一条
	ReplyComment(ctx context.Context, userID int, resourceID string, parentCommentID string, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, resourceID string, replyID string) (map[string]interface{}, error)

	// 浏览量
	AddView(ctx context.Context, resourceID string) (int, error)

	GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error) // 新增方法
}

type toolRepository struct {
	db *Database
}

func NewToolRepository(db *Database) ToolRepository {
	return &toolRepository{db: db}
}

func (r *toolRepository) GetTools(ctx context.Context, category, tags []string, sort, cursor string, pageSize int) ([]map[string]interface{}, error) {
	if pageSize <= 0 {
		pageSize = 10
	}

	var (
		whereParts []string
		args       []interface{}
	)
	whereParts = append(whereParts, "t.resource_type = 'tool'")

	// category 过滤（支持多个）
	if len(category) > 0 {
		whereParts = append(whereParts, fmt.Sprintf("t.category IN (%s)", placeholders(len(category))))
		for _, c := range category {
			args = append(args, c)
		}
	}

	// tags 过滤（任意一个tag命中即可）
	if len(tags) > 0 {
		whereParts = append(whereParts, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM tool_tags tt2 WHERE tt2.tool_id = t.resource_id AND tt2.tag IN (%s))",
			placeholders(len(tags)),
		))
		for _, t := range tags {
			args = append(args, t)
		}
	}

	// cursor 分页（按 resource_id 递减翻页）
	if cursor != "" {
		whereParts = append(whereParts, "t.resource_id < ?")
		args = append(args, cursor)
	}

	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = "WHERE " + strings.Join(whereParts, " AND ")
	}

	orderBy := "t.created_at DESC, t.resource_id DESC"
	switch strings.ToLower(sort) {
	case "views":
		orderBy = "t.views DESC, t.resource_id DESC"
	case "collections":
		orderBy = "t.collections DESC, t.resource_id DESC"
	case "loves", "likes":
		orderBy = "t.loves DESC, t.resource_id DESC"
	}

	// 列表一次性聚合 tags / contributors，image 取最小 sort_order 的那张
	query := fmt.Sprintf(`
		SELECT
			t.resource_id,
			t.resource_type,
			t.resource_name,
			t.description,
			COALESCE(MIN(ti.image_url), '') AS image,
			COALESCE(t.category, '') AS catagory,
			t.views,
			t.collections,
			t.loves,
			COALESCE(GROUP_CONCAT(DISTINCT tt.tag SEPARATOR ','), '') AS tags_csv,
			COALESCE(GROUP_CONCAT(DISTINCT u.username SEPARATOR ','), '') AS contributors_csv,
			t.created_at
		FROM tools t
		LEFT JOIN tool_images ti ON ti.tool_id = t.resource_id
		LEFT JOIN tool_tags tt ON tt.tool_id = t.resource_id
		LEFT JOIN tool_contributors tc ON tc.tool_id = t.resource_id
		LEFT JOIN users u ON u.id = tc.user_id
		%s
		GROUP BY t.resource_id, t.resource_type, t.resource_name, t.description, t.category, t.views, t.collections, t.loves, t.created_at
		ORDER BY %s
		LIMIT ?
	`, whereSQL, orderBy)

	args = append(args, pageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tools: %v", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var (
			id              int
			resourceType    string
			resourceName    string
			description     sql.NullString
			image           sql.NullString
			catagory        sql.NullString
			views           int
			collections     int
			loves           int
			tagsCSV         sql.NullString
			contributorsCSV sql.NullString
			createdAt       time.Time
		)
		if err := rows.Scan(
			&id,
			&resourceType,
			&resourceName,
			&description,
			&image,
			&catagory,
			&views,
			&collections,
			&loves,
			&tagsCSV,
			&contributorsCSV,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan tool row: %v", err)
		}

		result = append(result, map[string]interface{}{
			"resourceId":   id,
			"resourceType": resourceType,
			"resourceName": resourceName,
			"description":  nullString(description),
			"image":        nullString(image),
			"catagory":     nullString(catagory),
			"tags":         splitCSV(nullString(tagsCSV)),
			"views":        views,
			"collections":  collections,
			"loves":        loves,
			"contributors": splitCSV(nullString(contributorsCSV)),
			"createdat":    createdAt.Format("2006-01-02"),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tool rows: %v", err)
	}

	return result, nil
}

func (r *toolRepository) GetByID(ctx context.Context, resourceID string, userID int) (map[string]interface{}, error) {
	query := `
		SELECT
			resource_id,
			resource_type,
			resource_name,
			resource_link,
			description,
			description_detail,
			category,
			views,
			collections,
			loves,
			created_at
		FROM tools
		WHERE resource_id = ?
		LIMIT 1
	`

	var (
		id             int
		resourceType   string
		resourceName   string
		resourceLink   sql.NullString
		description    sql.NullString
		descriptionDtl sql.NullString
		category       sql.NullString
		views          int
		collections    int
		loves          int
		createdAt      time.Time
	)

	err := r.db.QueryRowContext(ctx, query, resourceID).Scan(
		&id,
		&resourceType,
		&resourceName,
		&resourceLink,
		&description,
		&descriptionDtl,
		&category,
		&views,
		&collections,
		&loves,
		&createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get tool by id: %v", err)
	}

	images, err := r.fetchToolImages(ctx, id)
	if err != nil {
		return nil, err
	}
	tags, err := r.fetchToolTags(ctx, id)
	if err != nil {
		return nil, err
	}
	contributors, err := r.fetchToolContributors(ctx, id)
	if err != nil {
		return nil, err
	}

	isLiked := false
	isCollected := false
	if userID > 0 {
		if v, err := r.isToolLiked(ctx, userID, id); err == nil {
			isLiked = v
		} else {
			return nil, err
		}
		if v, err := r.isToolCollected(ctx, userID, id); err == nil {
			isCollected = v
		} else {
			return nil, err
		}
	}

	comments, commentCount, err := r.fetchToolComments(ctx, id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"resourceId":         id,
		"resourceType":       resourceType,
		"resourceName":       resourceName,
		"resourceLink":       nullString(resourceLink),
		"description":        nullString(description),
		"description_detail": nullString(descriptionDtl),
		"catagory":           nullString(category),
		"image":              images,
		"tags":               tags,
		"contributors":       contributors,
		"views":              views,
		"collections":        collections,
		"loves":              loves,
		"iscollected":        isCollected,
		"isliked":            isLiked,
		"comment_count":      commentCount,
		"comments":           comments,
		"createdDate":        createdAt.Format("2006-01-02"),
	}, nil
}

func (r *toolRepository) Search(ctx context.Context, keyword, cursor string, pageSize int) ([]map[string]interface{}, error) {
	if pageSize <= 0 {
		pageSize = 10
	}

	var (
		whereParts []string
		args       []interface{}
	)

	whereParts = append(whereParts, "t.resource_type = 'tool'")
	if keyword != "" {
		whereParts = append(whereParts, "(t.resource_name LIKE ? OR t.description LIKE ?)")
		like := "%" + keyword + "%"
		args = append(args, like, like)
	}
	if cursor != "" {
		whereParts = append(whereParts, "t.resource_id < ?")
		args = append(args, cursor)
	}

	whereSQL := "WHERE " + strings.Join(whereParts, " AND ")

	query := fmt.Sprintf(`
		SELECT
			t.resource_id,
			t.resource_type,
			t.resource_name,
			t.description,
			COALESCE(MIN(ti.image_url), '') AS image,
			COALESCE(t.category, '') AS catagory,
			t.views,
			t.collections,
			t.loves,
			COALESCE(GROUP_CONCAT(DISTINCT tt.tag SEPARATOR ','), '') AS tags_csv,
			COALESCE(GROUP_CONCAT(DISTINCT u.username SEPARATOR ','), '') AS contributors_csv,
			t.created_at
		FROM tools t
		LEFT JOIN tool_images ti ON ti.tool_id = t.resource_id
		LEFT JOIN tool_tags tt ON tt.tool_id = t.resource_id
		LEFT JOIN tool_contributors tc ON tc.tool_id = t.resource_id
		LEFT JOIN users u ON u.id = tc.user_id
		%s
		GROUP BY t.resource_id, t.resource_type, t.resource_name, t.description, t.category, t.views, t.collections, t.loves, t.created_at
		ORDER BY t.created_at DESC, t.resource_id DESC
		LIMIT ?
	`, whereSQL)
	args = append(args, pageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search tools: %v", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var (
			id              int
			resourceType    string
			resourceName    string
			description     sql.NullString
			image           sql.NullString
			catagory        sql.NullString
			views           int
			collections     int
			loves           int
			tagsCSV         sql.NullString
			contributorsCSV sql.NullString
			createdAt       time.Time
		)
		if err := rows.Scan(
			&id,
			&resourceType,
			&resourceName,
			&description,
			&image,
			&catagory,
			&views,
			&collections,
			&loves,
			&tagsCSV,
			&contributorsCSV,
			&createdAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan tool search row: %v", err)
		}

		result = append(result, map[string]interface{}{
			"resourceId":   id,
			"resourceType": resourceType,
			"resourceName": resourceName,
			"description":  nullString(description),
			"image":        nullString(image),
			"catagory":     nullString(catagory),
			"tags":         splitCSV(nullString(tagsCSV)),
			"views":        views,
			"collections":  collections,
			"loves":        loves,
			"contributors": splitCSV(nullString(contributorsCSV)),
			"createdat":    createdAt.Format("2006-01-02"),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tool search rows: %v", err)
	}

	return result, nil
}

func (r *toolRepository) Create(ctx context.Context, userID int, data map[string]interface{}) (map[string]interface{}, error) {
	// 最小可用实现：写 tools + tool_tags（图片/贡献者后续可扩展）
	name, _ := data["name"].(string)
	link, _ := data["link"].(string)
	description, _ := data["description"].(string)
	descriptionDetail, _ := data["description_detail"].(string)
	category, _ := data["category"].(string)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT INTO tools (resource_type, resource_name, resource_link, description, description_detail, category, status, submitter_id)
		VALUES ('tool', ?, ?, ?, ?, ?, 'pending', ?)
	`, name, link, description, descriptionDetail, category, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert tool: %v", err)
	}

	newID64, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inserted id: %v", err)
	}
	newID := int(newID64)

	// tags
	if tagsAny, ok := data["tags"]; ok {
		if tags, ok := tagsAny.([]string); ok {
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if tag == "" {
					continue
				}
				_, _ = tx.ExecContext(ctx, `INSERT IGNORE INTO tool_tags (tool_id, tag) VALUES (?, ?)`, newID, tag)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	return map[string]interface{}{
		"resourceId":   newID,
		"resourceType": "tool",
		"resource":     link,
		"auditStatus":  "pending",
		"submitTime":   time.Now().Format("2006-01-02 15:04:05"),
		"auditTime":    nil,
		"rejectReason": nil,
	}, nil
}

func (r *toolRepository) AddLike(ctx context.Context, userID int, resourceID string) error {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return fmt.Errorf("invalid resource id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT IGNORE INTO likes (user_id, resource_type, resource_id)
		VALUES (?, 'tool', ?)
	`, userID, toolID)
	if err != nil {
		return fmt.Errorf("failed to insert like: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE tools SET loves = loves + 1 WHERE resource_id = ?`, toolID); err != nil {
			return fmt.Errorf("failed to update tool loves: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %v", err)
	}
	return nil
}

func (r *toolRepository) GetLikes(ctx context.Context, resourceID string) (int, error) {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return 0, fmt.Errorf("invalid resource id")
	}

	var loves int
	err = r.db.QueryRowContext(ctx, `SELECT loves FROM tools WHERE resource_id = ?`, toolID).Scan(&loves)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get likes: %v", err)
	}
	return loves, nil
}

func (r *toolRepository) RemoveLike(ctx context.Context, userID int, resourceID string) error {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return fmt.Errorf("invalid resource id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM likes
		WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?
	`, userID, toolID)
	if err != nil {
		return fmt.Errorf("failed to delete like: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE tools SET loves = GREATEST(loves - 1, 0) WHERE resource_id = ?`, toolID); err != nil {
			return fmt.Errorf("failed to decrement tool loves: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %v", err)
	}
	return nil
}

func (r *toolRepository) AddCollection(ctx context.Context, userID int, resourceID string) error {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return fmt.Errorf("invalid resource id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT IGNORE INTO collections (user_id, resource_type, resource_id)
		VALUES (?, 'tool', ?)
	`, userID, toolID)
	if err != nil {
		return fmt.Errorf("failed to insert collection: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE tools SET collections = collections + 1 WHERE resource_id = ?`, toolID); err != nil {
			return fmt.Errorf("failed to update tool collections: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %v", err)
	}
	return nil
}

func (r *toolRepository) RemoveCollection(ctx context.Context, userID int, resourceID string) error {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return fmt.Errorf("invalid resource id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM collections
		WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?
	`, userID, toolID)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE tools SET collections = GREATEST(collections - 1, 0) WHERE resource_id = ?`, toolID); err != nil {
			return fmt.Errorf("failed to decrement tool collections: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %v", err)
	}
	return nil
}

func (r *toolRepository) GetCollections(ctx context.Context, resourceID string) (int, error) {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return 0, fmt.Errorf("invalid resource id")
	}

	var c int
	err = r.db.QueryRowContext(ctx, `SELECT collections FROM tools WHERE resource_id = ?`, toolID).Scan(&c)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get collections: %v", err)
	}
	return c, nil
}

func (r *toolRepository) AddView(ctx context.Context, resourceID string) (int, error) {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return 0, fmt.Errorf("invalid resource id")
	}

	if _, err := r.db.ExecContext(ctx, `UPDATE tools SET views = views + 1 WHERE resource_id = ?`, toolID); err != nil {
		return 0, fmt.Errorf("failed to update views: %v", err)
	}

	var views int
	if err := r.db.QueryRowContext(ctx, `SELECT views FROM tools WHERE resource_id = ?`, toolID).Scan(&views); err != nil {
		return 0, fmt.Errorf("failed to read views: %v", err)
	}
	return views, nil
}

func (r *toolRepository) AddComment(ctx context.Context, userID int, resourceID string, content string) (map[string]interface{}, error) {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return nil, fmt.Errorf("invalid resource id")
	}

	res, err := r.db.ExecContext(ctx, `
		INSERT INTO comments (resource_type, resource_id, parent_id, user_id, content)
		VALUES ('tool', ?, NULL, ?, ?)
	`, toolID, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to insert comment: %v", err)
	}
	commentID64, _ := res.LastInsertId()
	commentID := int(commentID64)

	nickname, avatar, err := r.fetchUserDisplay(ctx, userID)
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

func (r *toolRepository) DeleteComment(ctx context.Context, userID int, resourceID string, commentID string) (map[string]interface{}, error) {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return nil, fmt.Errorf("invalid resource id")
	}

	// 允许 commentID 为空：删除该用户最新一条一级评论
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
			WHERE comment_id = ? AND resource_type = 'tool' AND resource_id = ? AND user_id = ? AND parent_id IS NULL AND deleted_at IS NULL
			LIMIT 1
		`, targetID, toolID, userID).Scan(&targetID, &content)
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
			WHERE resource_type = 'tool' AND resource_id = ? AND user_id = ? AND parent_id IS NULL AND deleted_at IS NULL
			ORDER BY created_at DESC, comment_id DESC
			LIMIT 1
		`, toolID, userID).Scan(&targetID, &content)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("comment not found")
			}
			return nil, fmt.Errorf("failed to read latest comment: %v", err)
		}
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE comments
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE comment_id = ? AND user_id = ? AND resource_type = 'tool' AND resource_id = ? AND deleted_at IS NULL
	`, targetID, userID, toolID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete comment: %v", err)
	}

	nickname, avatar, err := r.fetchUserDisplay(ctx, userID)
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

func (r *toolRepository) ReplyComment(ctx context.Context, userID int, resourceID string, parentCommentID string, content string) (map[string]interface{}, error) {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return nil, fmt.Errorf("invalid resource id")
	}
	parentID, err := strconv.Atoi(parentCommentID)
	if err != nil {
		return nil, fmt.Errorf("invalid comment id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	// 确认父评论存在且是一级评论
	var exists int
	if err := tx.QueryRowContext(ctx, `
		SELECT 1
		FROM comments
		WHERE comment_id = ? AND resource_type = 'tool' AND resource_id = ? AND parent_id IS NULL AND deleted_at IS NULL
		LIMIT 1
	`, parentID, toolID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("parent comment not found")
		}
		return nil, fmt.Errorf("failed to check parent comment: %v", err)
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO comments (resource_type, resource_id, parent_id, user_id, content)
		VALUES ('tool', ?, ?, ?, ?)
	`, toolID, parentID, userID, content)
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

	nickname, avatar, err := r.fetchUserDisplay(ctx, userID)
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

func (r *toolRepository) DeleteReply(ctx context.Context, userID int, resourceID string, replyID string) (map[string]interface{}, error) {
	toolID, err := strconv.Atoi(resourceID)
	if err != nil {
		return nil, fmt.Errorf("invalid resource id")
	}
	rid, err := strconv.Atoi(replyID)
	if err != nil {
		return nil, fmt.Errorf("invalid reply id")
	}

	var parentID sql.NullInt64
	err = r.db.QueryRowContext(ctx, `
		SELECT parent_id
		FROM comments
		WHERE comment_id = ? AND resource_type = 'tool' AND resource_id = ? AND user_id = ? AND deleted_at IS NULL
		LIMIT 1
	`, rid, toolID, userID).Scan(&parentID)
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
		WHERE comment_id = ? AND user_id = ? AND resource_type = 'tool' AND resource_id = ? AND deleted_at IS NULL
	`, rid, userID, toolID); err != nil {
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

	nickname, avatar, err := r.fetchUserDisplay(ctx, userID)
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

func (r *toolRepository) GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error) {
	// 实现获取待审核工具的逻辑
	return []map[string]interface{}{
		{
			"submitor":     "user1",
			"submitDate":   "2023-12-01 10:00:00",
			"reourceId":    1,
			"resourceType": "tool",
			"resourcename": "新工具",
			"catagory":     "软件开发",
			"link":         "https://example.com/tool",
			"description":  "这是一个新工具",
			"tags":         []string{"AI", "免费"},
			"file":         "tool.zip",
		},
	}, nil
}

func (r *toolRepository) fetchToolImages(ctx context.Context, toolID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT image_url FROM tool_images WHERE tool_id = ? ORDER BY sort_order ASC, id ASC`, toolID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tool images: %v", err)
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var url sql.NullString
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("failed to scan tool image: %v", err)
		}
		if s := nullString(url); s != "" {
			images = append(images, s)
		}
	}
	return images, rows.Err()
}

func (r *toolRepository) fetchToolTags(ctx context.Context, toolID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT tag FROM tool_tags WHERE tool_id = ? ORDER BY tag ASC`, toolID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tool tags: %v", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag sql.NullString
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("failed to scan tool tag: %v", err)
		}
		if s := nullString(tag); s != "" {
			tags = append(tags, s)
		}
	}
	return tags, rows.Err()
}

func (r *toolRepository) fetchToolContributors(ctx context.Context, toolID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.username
		FROM tool_contributors tc
		JOIN users u ON u.id = tc.user_id
		WHERE tc.tool_id = ?
		ORDER BY u.username ASC
	`, toolID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tool contributors: %v", err)
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

func (r *toolRepository) fetchUserDisplay(ctx context.Context, userID int) (string, string, error) {
	var nickname sql.NullString
	var username sql.NullString
	var avatar sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT nickname, username, avatar FROM users WHERE id = ? LIMIT 1`, userID).Scan(&nickname, &username, &avatar)
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

func (r *toolRepository) isToolLiked(ctx context.Context, userID int, toolID int) (bool, error) {
	var one int
	err := r.db.QueryRowContext(ctx, `
		SELECT 1
		FROM likes
		WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?
		LIMIT 1
	`, userID, toolID).Scan(&one)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check isliked: %v", err)
	}
	return true, nil
}

func (r *toolRepository) isToolCollected(ctx context.Context, userID int, toolID int) (bool, error) {
	var one int
	err := r.db.QueryRowContext(ctx, `
		SELECT 1
		FROM collections
		WHERE user_id = ? AND resource_type = 'tool' AND resource_id = ?
		LIMIT 1
	`, userID, toolID).Scan(&one)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check iscollected: %v", err)
	}
	return true, nil
}

type toolCommentRow struct {
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

func (r *toolRepository) fetchToolComments(ctx context.Context, toolID int) ([]map[string]interface{}, int, error) {
	// 拉取该工具下所有“未删除”的评论（包含一级评论和回复），一次查询后在内存里组装树形结构
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
		WHERE c.resource_type = 'tool' AND c.resource_id = ? AND c.deleted_at IS NULL
		ORDER BY c.created_at ASC, c.comment_id ASC
	`, toolID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query comments: %v", err)
	}
	defer rows.Close()

	var all []toolCommentRow
	for rows.Next() {
		var rrow toolCommentRow
		if err := rows.Scan(
			&rrow.ID,
			&rrow.ParentID,
			&rrow.UserID,
			&rrow.Nickname,
			&rrow.Username,
			&rrow.Avatar,
			&rrow.Content,
			&rrow.LoveCount,
			&rrow.ReplyTotal,
			&rrow.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan comment row: %v", err)
		}
		all = append(all, rrow)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate comment rows: %v", err)
	}

	// 顶层评论（parent_id NULL）与 replies 分组
	top := make([]toolCommentRow, 0)
	repliesByParent := make(map[int][]toolCommentRow)
	for _, c := range all {
		if !c.ParentID.Valid {
			top = append(top, c)
			continue
		}
		pid := int(c.ParentID.Int64)
		repliesByParent[pid] = append(repliesByParent[pid], c)
	}

	commentCount := len(top)
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

	return out, commentCount, nil
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	// "?, ?, ?"
	return strings.TrimRight(strings.Repeat("?,", n), ",")
}

func splitCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
