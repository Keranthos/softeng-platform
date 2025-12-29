package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ProjectRepository interface {
	GetProjects(ctx context.Context, category string, techStack []string, sort string, limit int, cursor string) ([]map[string]interface{}, error)
	GetByID(ctx context.Context, projectID string, userID int) (map[string]interface{}, error)
	Search(ctx context.Context, keyword string, category []string, cursor string, limit int) ([]map[string]interface{}, error)
	Create(ctx context.Context, userID int, data map[string]interface{}) (map[string]interface{}, error)
	Update(ctx context.Context, userID int, projectID string, data map[string]interface{}) (map[string]interface{}, error)
	LikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
	UnlikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
	AddComment(ctx context.Context, userID int, projectID, content string) (map[string]interface{}, error)
	DeleteComment(ctx context.Context, userID int, projectID string, commentID string) (map[string]interface{}, error) // commentID 为空则删除该用户最新一条
	ReplyComment(ctx context.Context, userID int, projectID, commentID, content string) (map[string]interface{}, error)
	DeleteReply(ctx context.Context, userID int, projectID, commentID string) (map[string]interface{}, error)
	AddView(ctx context.Context, projectID string) (int, error)
	CollectProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
	UncollectProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error)
	GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error) // 新增方法
}

type projectRepository struct {
	db *Database
}

func NewProjectRepository(db *Database) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) GetProjects(ctx context.Context, category string, techStack []string, sort string, limit int, cursor string) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	var (
		whereParts []string
		args       []interface{}
	)

	whereParts = append(whereParts, "p.resource_type = 'project'")

	if category != "" {
		whereParts = append(whereParts, "p.category = ?")
		args = append(args, category)
	}

	if len(techStack) > 0 {
		whereParts = append(whereParts, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM project_tech_stack pts2 WHERE pts2.project_id = p.project_id AND pts2.tech IN (%s))",
			placeholders(len(techStack)),
		))
		for _, tech := range techStack {
			args = append(args, tech)
		}
	}

	if cursor != "" {
		whereParts = append(whereParts, "p.project_id < ?")
		args = append(args, cursor)
	}

	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = "WHERE " + strings.Join(whereParts, " AND ")
	}

	orderBy := "p.created_at DESC, p.project_id DESC"
	switch strings.ToLower(sort) {
	case "views":
		orderBy = "p.views DESC, p.project_id DESC"
	case "collections":
		orderBy = "p.collections DESC, p.project_id DESC"
	case "loves", "likes":
		orderBy = "p.loves DESC, p.project_id DESC"
	}

	query := fmt.Sprintf(`
		SELECT
			p.project_id,
			p.resource_type,
			p.name,
			p.description,
			COALESCE(p.category, '') AS category,
			COALESCE(GROUP_CONCAT(DISTINCT pts.tech SEPARATOR ','), '') AS tech_csv,
			COALESCE(GROUP_CONCAT(DISTINCT u.username SEPARATOR ','), '') AS author_csv,
			COALESCE(p.cover, '') AS cover,
			p.created_at,
			p.loves,
			p.collections,
			p.views
		FROM projects p
		LEFT JOIN project_tech_stack pts ON pts.project_id = p.project_id
		LEFT JOIN project_authors pa ON pa.project_id = p.project_id
		LEFT JOIN users u ON u.id = pa.user_id
		%s
		GROUP BY p.project_id, p.resource_type, p.name, p.description, p.category, p.cover, p.created_at, p.loves, p.collections, p.views
		ORDER BY %s
		LIMIT ?
	`, whereSQL, orderBy)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %v", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var (
			projectID    int
			resourceType sql.NullString
			name         sql.NullString
			description  sql.NullString
			cat          sql.NullString
			techCSV      sql.NullString
			authorCSV    sql.NullString
			cover        sql.NullString
			createdAt    time.Time
			loves        int
			collections  int
			views        int
		)

		if err := rows.Scan(
			&projectID,
			&resourceType,
			&name,
			&description,
			&cat,
			&techCSV,
			&authorCSV,
			&cover,
			&createdAt,
			&loves,
			&collections,
			&views,
		); err != nil {
			return nil, fmt.Errorf("failed to scan project row: %v", err)
		}

		result = append(result, map[string]interface{}{
			"projectId":    projectID,
			"resourceType": nullString(resourceType),
			"name":         nullString(name),
			"description":  nullString(description),
			"category":     nullString(cat),
			"techStack":    splitCSV(nullString(techCSV)),
			"likecount":    loves,
			"authername":   splitCSV(nullString(authorCSV)),
			"cover":        nullString(cover),
			"createdat":    createdAt.Format("2006-01-02"),
			"loves":        loves,
			"collections":  collections,
			"views":        views,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate project rows: %v", err)
	}

	return result, nil
}

func (r *projectRepository) GetByID(ctx context.Context, projectID string, userID int) (map[string]interface{}, error) {
	query := `
		SELECT
			project_id,
			resource_type,
			name,
			description,
			detail,
			github_url,
			category,
			cover,
			views,
			loves,
			collections,
			created_at
		FROM projects
		WHERE project_id = ?
		LIMIT 1
	`

	var (
		id           int
		resourceType sql.NullString
		name         sql.NullString
		description  sql.NullString
		detail       sql.NullString
		githubURL    sql.NullString
		cat          sql.NullString
		cover        sql.NullString
		views        int
		loves        int
		collections  int
		createdAt    time.Time
	)

	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&id,
		&resourceType,
		&name,
		&description,
		&detail,
		&githubURL,
		&cat,
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
		return nil, fmt.Errorf("failed to get project by id: %v", err)
	}

	techStack, err := r.fetchProjectTechStack(ctx, id)
	if err != nil {
		return nil, err
	}
	images, err := r.fetchProjectImages(ctx, id)
	if err != nil {
		return nil, err
	}
	authors, err := r.fetchProjectAuthors(ctx, id)
	if err != nil {
		return nil, err
	}

	isLiked := false
	isCollected := false
	if userID > 0 {
		if v, err := r.isProjectLiked(ctx, userID, id); err == nil {
			isLiked = v
		} else {
			return nil, err
		}
		if v, err := r.isProjectCollected(ctx, userID, id); err == nil {
			isCollected = v
		} else {
			return nil, err
		}
	}

	comments, commentCount, err := r.fetchProjectComments(ctx, id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"projectId":     id,
		"resourceType":  nullString(resourceType),
		"name":          nullString(name),
		"description":   nullString(description),
		"detail":        nullString(detail),
		"githubURL":     nullString(githubURL),
		"techStack":     techStack,
		"catagory":      nullString(cat),
		"cover":         nullString(cover),
		"images":        images,
		"likes":         loves,
		"views":         views,
		"collections":   collections,
		"isliked":       isLiked,
		"iscollected":   isCollected,
		"author":        authors,
		"comment_count": commentCount,
		"comments":      comments,
		"createdAt":     createdAt.Format("2006-01-02"),
	}, nil
}

func (r *projectRepository) Search(ctx context.Context, keyword string, category []string, cursor string, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	var (
		whereParts []string
		args       []interface{}
	)

	whereParts = append(whereParts, "p.resource_type = 'project'")

	if keyword != "" {
		whereParts = append(whereParts, "(p.name LIKE ? OR p.description LIKE ?)")
		like := "%" + keyword + "%"
		args = append(args, like, like)
	}

	if len(category) > 0 {
		whereParts = append(whereParts, fmt.Sprintf("p.category IN (%s)", placeholders(len(category))))
		for _, c := range category {
			args = append(args, c)
		}
	}

	if cursor != "" {
		whereParts = append(whereParts, "p.project_id < ?")
		args = append(args, cursor)
	}

	whereSQL := "WHERE " + strings.Join(whereParts, " AND ")

	query := fmt.Sprintf(`
		SELECT
			p.project_id,
			p.resource_type,
			p.name,
			p.description,
			COALESCE(p.category, '') AS category,
			COALESCE(GROUP_CONCAT(DISTINCT pts.tech SEPARATOR ','), '') AS tech_csv,
			COALESCE(GROUP_CONCAT(DISTINCT u.username SEPARATOR ','), '') AS author_csv,
			COALESCE(p.cover, '') AS cover,
			p.created_at,
			p.loves,
			p.collections,
			p.views
		FROM projects p
		LEFT JOIN project_tech_stack pts ON pts.project_id = p.project_id
		LEFT JOIN project_authors pa ON pa.project_id = p.project_id
		LEFT JOIN users u ON u.id = pa.user_id
		%s
		GROUP BY p.project_id, p.resource_type, p.name, p.description, p.category, p.cover, p.created_at, p.loves, p.collections, p.views
		ORDER BY p.created_at DESC, p.project_id DESC
		LIMIT ?
	`, whereSQL)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search projects: %v", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var (
			id           int
			resourceType sql.NullString
			name         sql.NullString
			description  sql.NullString
			cat          sql.NullString
			techCSV      sql.NullString
			authorCSV    sql.NullString
			cover        sql.NullString
			createdAt    time.Time
			loves        int
			collections  int
			views        int
		)

		if err := rows.Scan(
			&id,
			&resourceType,
			&name,
			&description,
			&cat,
			&techCSV,
			&authorCSV,
			&cover,
			&createdAt,
			&loves,
			&collections,
			&views,
		); err != nil {
			return nil, fmt.Errorf("failed to scan project search row: %v", err)
		}

		result = append(result, map[string]interface{}{
			"projectId":    id,
			"resourceType": nullString(resourceType),
			"name":         nullString(name),
			"description":  nullString(description),
			"category":     nullString(cat),
			"techStack":    splitCSV(nullString(techCSV)),
			"likecount":    loves,
			"authername":   splitCSV(nullString(authorCSV)),
			"cover":        nullString(cover),
			"createdat":    createdAt.Format("2006-01-02"),
			"loves":        loves,
			"collections":  collections,
			"views":        views,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate project search rows: %v", err)
	}

	return result, nil
}

func (r *projectRepository) Create(ctx context.Context, userID int, data map[string]interface{}) (map[string]interface{}, error) {
	name, _ := data["name"].(string)
	description, _ := data["description"].(string)
	detail, _ := data["detail"].(string)
	github, _ := data["github"].(string)
	category, _ := data["category"].(string)

	var techStack []string
	if v, ok := data["techStack"].([]string); ok {
		techStack = v
	}
	var images []string
	if v, ok := data["images"].([]string); ok {
		images = v
	}

	cover := ""
	if len(images) > 0 {
		cover = images[0]
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT INTO projects (resource_type, name, description, detail, github_url, category, cover, views, loves, collections)
		VALUES ('project', ?, ?, ?, ?, ?, ?, 0, 0, 0)
	`, name, description, detail, github, category, cover)
	if err != nil {
		return nil, fmt.Errorf("failed to insert project: %v", err)
	}
	pid64, _ := res.LastInsertId()
	projectID := int(pid64)

	// authors：默认提交者就是作者
	if _, err := tx.ExecContext(ctx, `
		INSERT IGNORE INTO project_authors (project_id, user_id)
		VALUES (?, ?)
	`, projectID, userID); err != nil {
		return nil, fmt.Errorf("failed to insert project author: %v", err)
	}

	// tech stack
	if _, err := tx.ExecContext(ctx, `DELETE FROM project_tech_stack WHERE project_id = ?`, projectID); err != nil {
		return nil, fmt.Errorf("failed to clear project tech stack: %v", err)
	}
	for _, t := range techStack {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO project_tech_stack (project_id, tech)
			VALUES (?, ?)
		`, projectID, t); err != nil {
			return nil, fmt.Errorf("failed to insert project tech: %v", err)
		}
	}

	// images
	if _, err := tx.ExecContext(ctx, `DELETE FROM project_images WHERE project_id = ?`, projectID); err != nil {
		return nil, fmt.Errorf("failed to clear project images: %v", err)
	}
	for i, url := range images {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO project_images (project_id, image_url, sort_order)
			VALUES (?, ?, ?)
		`, projectID, url, i); err != nil {
			return nil, fmt.Errorf("failed to insert project image: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	return map[string]interface{}{
		"resourceId":   projectID,
		"resourceType": "project",
		"resource":     github,
		"auditStatus":  "pending",
		"submitTime":   time.Now().Format("2006-01-02 15:04:05"),
		"auditTime":    nil,
		"rejectReason": nil,
	}, nil
}

func (r *projectRepository) Update(ctx context.Context, userID int, projectID string, data map[string]interface{}) (map[string]interface{}, error) {
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id")
	}

	name, _ := data["name"].(string)
	description, _ := data["description"].(string)
	detail, _ := data["detail"].(string)
	github, _ := data["github"].(string)
	category, _ := data["category"].(string)

	var techStack []string
	if v, ok := data["techStack"].([]string); ok {
		techStack = v
	}
	var images []string
	if v, ok := data["images"].([]string); ok {
		images = v
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	// 简单权限：只有作者能更新（最小可用）
	var isAuthor int
	if err := tx.QueryRowContext(ctx, `
		SELECT 1
		FROM project_authors
		WHERE project_id = ? AND user_id = ?
		LIMIT 1
	`, pid, userID).Scan(&isAuthor); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no permission")
		}
		return nil, fmt.Errorf("failed to check author permission: %v", err)
	}

	coverSQL := ""
	args := []interface{}{name, description, detail, github, category}
	if len(images) > 0 && strings.TrimSpace(images[0]) != "" {
		coverSQL = ", cover = ?"
		args = append(args, strings.TrimSpace(images[0]))
	}
	args = append(args, pid)

	query := fmt.Sprintf(`
		UPDATE projects
		SET name = ?, description = ?, detail = ?, github_url = ?, category = ?%s, updated_at = NOW()
		WHERE project_id = ?
	`, coverSQL)

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("failed to update project: %v", err)
	}

	// tech stack
	if _, err := tx.ExecContext(ctx, `DELETE FROM project_tech_stack WHERE project_id = ?`, pid); err != nil {
		return nil, fmt.Errorf("failed to clear project tech stack: %v", err)
	}
	for _, t := range techStack {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO project_tech_stack (project_id, tech)
			VALUES (?, ?)
		`, pid, t); err != nil {
			return nil, fmt.Errorf("failed to insert project tech: %v", err)
		}
	}

	// images（传了就覆盖；没传就保留原图）
	if images != nil {
		if _, err := tx.ExecContext(ctx, `DELETE FROM project_images WHERE project_id = ?`, pid); err != nil {
			return nil, fmt.Errorf("failed to clear project images: %v", err)
		}
		for i, url := range images {
			url = strings.TrimSpace(url)
			if url == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO project_images (project_id, image_url, sort_order)
				VALUES (?, ?, ?)
			`, pid, url, i); err != nil {
				return nil, fmt.Errorf("failed to insert project image: %v", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	return map[string]interface{}{
		"resourceId":   pid,
		"resourceType": "project",
		"resource":     github,
		"auditStatus":  "pending",
		"submitTime":   time.Now().Format("2006-01-02 15:04:05"),
		"auditTime":    nil,
		"rejectReason": nil,
	}, nil
}

func (r *projectRepository) LikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT IGNORE INTO likes (user_id, resource_type, resource_id)
		VALUES (?, 'project', ?)
	`, userID, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to insert like: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE projects SET loves = loves + 1 WHERE project_id = ?`, pid); err != nil {
			return nil, fmt.Errorf("failed to update loves: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	var loves int
	if err := r.db.QueryRowContext(ctx, `SELECT loves FROM projects WHERE project_id = ?`, pid).Scan(&loves); err != nil {
		return nil, fmt.Errorf("failed to read loves: %v", err)
	}

	return map[string]interface{}{
		"likecounts": loves,
		"isliked":    true,
	}, nil
}

func (r *projectRepository) UnlikeProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM likes
		WHERE user_id = ? AND resource_type = 'project' AND resource_id = ?
	`, userID, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete like: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE projects SET loves = GREATEST(loves - 1, 0) WHERE project_id = ?`, pid); err != nil {
			return nil, fmt.Errorf("failed to decrement loves: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	var loves int
	if err := r.db.QueryRowContext(ctx, `SELECT loves FROM projects WHERE project_id = ?`, pid).Scan(&loves); err != nil {
		return nil, fmt.Errorf("failed to read loves: %v", err)
	}

	return map[string]interface{}{
		"likecounts": loves,
		"isliked":    false,
	}, nil
}

func (r *projectRepository) AddComment(ctx context.Context, userID int, projectID, content string) (map[string]interface{}, error) {
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id")
	}

	res, err := r.db.ExecContext(ctx, `
		INSERT INTO comments (resource_type, resource_id, parent_id, user_id, content)
		VALUES ('project', ?, NULL, ?, ?)
	`, pid, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to insert comment: %v", err)
	}
	commentID64, _ := res.LastInsertId()
	commentID := int(commentID64)

	nickname, avatar, err := fetchUserDisplayProject(ctx, r.db, userID)
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

func (r *projectRepository) DeleteComment(ctx context.Context, userID int, projectID string, commentID string) (map[string]interface{}, error) {
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id")
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
			WHERE comment_id = ? AND resource_type = 'project' AND resource_id = ? AND user_id = ? AND parent_id IS NULL AND deleted_at IS NULL
			LIMIT 1
		`, targetID, pid, userID).Scan(&targetID, &content)
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
			WHERE resource_type = 'project' AND resource_id = ? AND user_id = ? AND parent_id IS NULL AND deleted_at IS NULL
			ORDER BY created_at DESC, comment_id DESC
			LIMIT 1
		`, pid, userID).Scan(&targetID, &content)
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
		WHERE comment_id = ? AND user_id = ? AND resource_type = 'project' AND resource_id = ? AND deleted_at IS NULL
	`, targetID, userID, pid); err != nil {
		return nil, fmt.Errorf("failed to delete comment: %v", err)
	}

	nickname, avatar, err := fetchUserDisplayProject(ctx, r.db, userID)
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

func (r *projectRepository) ReplyComment(ctx context.Context, userID int, projectID, commentID, content string) (map[string]interface{}, error) {
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id")
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
		WHERE comment_id = ? AND resource_type = 'project' AND resource_id = ? AND parent_id IS NULL AND deleted_at IS NULL
		LIMIT 1
	`, parentID, pid).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("parent comment not found")
		}
		return nil, fmt.Errorf("failed to check parent comment: %v", err)
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO comments (resource_type, resource_id, parent_id, user_id, content)
		VALUES ('project', ?, ?, ?, ?)
	`, pid, parentID, userID, content)
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

	nickname, avatar, err := fetchUserDisplayProject(ctx, r.db, userID)
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

func (r *projectRepository) DeleteReply(ctx context.Context, userID int, projectID, commentID string) (map[string]interface{}, error) {
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id")
	}
	rid, err := strconv.Atoi(commentID)
	if err != nil {
		return nil, fmt.Errorf("invalid reply id")
	}

	var parentID sql.NullInt64
	err = r.db.QueryRowContext(ctx, `
		SELECT parent_id
		FROM comments
		WHERE comment_id = ? AND resource_type = 'project' AND resource_id = ? AND user_id = ? AND deleted_at IS NULL
		LIMIT 1
	`, rid, pid, userID).Scan(&parentID)
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
		WHERE comment_id = ? AND user_id = ? AND resource_type = 'project' AND resource_id = ? AND deleted_at IS NULL
	`, rid, userID, pid); err != nil {
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

	nickname, avatar, err := fetchUserDisplayProject(ctx, r.db, userID)
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

func (r *projectRepository) AddView(ctx context.Context, projectID string) (int, error) {
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return 0, fmt.Errorf("invalid project id")
	}
	if _, err := r.db.ExecContext(ctx, `UPDATE projects SET views = views + 1 WHERE project_id = ?`, pid); err != nil {
		return 0, fmt.Errorf("failed to update views: %v", err)
	}
	var views int
	if err := r.db.QueryRowContext(ctx, `SELECT views FROM projects WHERE project_id = ?`, pid).Scan(&views); err != nil {
		return 0, fmt.Errorf("failed to read views: %v", err)
	}
	return views, nil
}

func (r *projectRepository) CollectProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		INSERT IGNORE INTO collections (user_id, resource_type, resource_id)
		VALUES (?, 'project', ?)
	`, userID, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to insert collection: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE projects SET collections = collections + 1 WHERE project_id = ?`, pid); err != nil {
			return nil, fmt.Errorf("failed to update collections: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	var collections int
	if err := r.db.QueryRowContext(ctx, `SELECT collections FROM projects WHERE project_id = ?`, pid).Scan(&collections); err != nil {
		return nil, fmt.Errorf("failed to read collections: %v", err)
	}

	return map[string]interface{}{
		"iscollected": true,
		"collections": collections,
	}, nil
}

func (r *projectRepository) UncollectProject(ctx context.Context, userID int, projectID string) (map[string]interface{}, error) {
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM collections
		WHERE user_id = ? AND resource_type = 'project' AND resource_id = ?
	`, userID, pid)
	if err != nil {
		return nil, fmt.Errorf("failed to delete collection: %v", err)
	}
	affected, _ := res.RowsAffected()
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE projects SET collections = GREATEST(collections - 1, 0) WHERE project_id = ?`, pid); err != nil {
			return nil, fmt.Errorf("failed to decrement collections: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %v", err)
	}

	var collections int
	if err := r.db.QueryRowContext(ctx, `SELECT collections FROM projects WHERE project_id = ?`, pid).Scan(&collections); err != nil {
		return nil, fmt.Errorf("failed to read collections: %v", err)
	}

	return map[string]interface{}{
		"iscollected": false,
		"collections": collections,
	}, nil
}

func fetchUserDisplayProject(ctx context.Context, db *Database, userID int) (string, string, error) {
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

func (r *projectRepository) isProjectLiked(ctx context.Context, userID int, projectID int) (bool, error) {
	var one int
	err := r.db.QueryRowContext(ctx, `
		SELECT 1
		FROM likes
		WHERE user_id = ? AND resource_type = 'project' AND resource_id = ?
		LIMIT 1
	`, userID, projectID).Scan(&one)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check isliked: %v", err)
	}
	return true, nil
}

func (r *projectRepository) isProjectCollected(ctx context.Context, userID int, projectID int) (bool, error) {
	var one int
	err := r.db.QueryRowContext(ctx, `
		SELECT 1
		FROM collections
		WHERE user_id = ? AND resource_type = 'project' AND resource_id = ?
		LIMIT 1
	`, userID, projectID).Scan(&one)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check iscollected: %v", err)
	}
	return true, nil
}

type projectCommentRow struct {
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

func (r *projectRepository) fetchProjectComments(ctx context.Context, projectID int) ([]map[string]interface{}, int, error) {
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
		WHERE c.resource_type = 'project' AND c.resource_id = ? AND c.deleted_at IS NULL
		ORDER BY c.created_at ASC, c.comment_id ASC
	`, projectID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query comments: %v", err)
	}
	defer rows.Close()

	var all []projectCommentRow
	for rows.Next() {
		var row projectCommentRow
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

	top := make([]projectCommentRow, 0)
	repliesByParent := make(map[int][]projectCommentRow)
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

func (r *projectRepository) GetPending(ctx context.Context, cursor, limit int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{
			"submitor":     "user2",
			"submitDate":   "2023-12-01 10:00:00",
			"reourceId":    1,
			"resourceType": "project",
			"resourcename": "新项目",
			"catagory":     "实训项目",
			"link":         "https://github.com/example/project",
			"description":  "项目描述",
			"tags":         []string{"Go", "React"},
			"file":         "project.zip",
		},
	}, nil
}

func (r *projectRepository) fetchProjectTechStack(ctx context.Context, projectID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT tech FROM project_tech_stack WHERE project_id = ? ORDER BY tech ASC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query project tech stack: %v", err)
	}
	defer rows.Close()

	var tech []string
	for rows.Next() {
		var t sql.NullString
		if err := rows.Scan(&t); err != nil {
			return nil, fmt.Errorf("failed to scan tech: %v", err)
		}
		if s := nullString(t); s != "" {
			tech = append(tech, s)
		}
	}
	return tech, rows.Err()
}

func (r *projectRepository) fetchProjectImages(ctx context.Context, projectID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT image_url FROM project_images WHERE project_id = ? ORDER BY sort_order ASC, id ASC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query project images: %v", err)
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var url sql.NullString
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("failed to scan project image: %v", err)
		}
		if s := nullString(url); s != "" {
			images = append(images, s)
		}
	}
	return images, rows.Err()
}

func (r *projectRepository) fetchProjectAuthors(ctx context.Context, projectID int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.username
		FROM project_authors pa
		JOIN users u ON u.id = pa.user_id
		WHERE pa.project_id = ?
		ORDER BY u.username ASC
	`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query project authors: %v", err)
	}
	defer rows.Close()

	var authors []string
	for rows.Next() {
		var username sql.NullString
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("failed to scan project author: %v", err)
		}
		if s := nullString(username); s != "" {
			authors = append(authors, s)
		}
	}
	return authors, rows.Err()
}
