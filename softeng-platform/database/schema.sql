-- 软件工程平台数据库表结构
-- 基于API文档数据模型设计

CREATE DATABASE IF NOT EXISTS softeng CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE softeng;

-- ==================== 用户相关表 ====================

-- 用户表（已更新，包含所有字段）
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE COMMENT '用户名',
    nickname VARCHAR(255) COMMENT '昵称',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱',
    password VARCHAR(255) NOT NULL COMMENT '密码（加密后）',
    avatar VARCHAR(500) COMMENT '头像地址',
    description TEXT COMMENT '个人动态描述',
    face_photo VARCHAR(500) COMMENT '封面地址',
    role VARCHAR(50) DEFAULT 'user' COMMENT '角色：user/admin',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ==================== 工具相关表 ====================

-- 工具表
CREATE TABLE IF NOT EXISTS tools (
    resource_id INT AUTO_INCREMENT PRIMARY KEY,
    resource_type VARCHAR(50) DEFAULT 'tool' COMMENT '资源类型',
    resource_name VARCHAR(255) NOT NULL COMMENT '工具名称',
    resource_link VARCHAR(500) COMMENT '工具访问链接',
    description VARCHAR(500) COMMENT '工具简介（20-50字）',
    description_detail TEXT COMMENT '工具详细介绍（50-100字）',
    category VARCHAR(100) COMMENT '工具类别',
    views INT DEFAULT 0 COMMENT '浏览量',
    collections INT DEFAULT 0 COMMENT '收藏量',
    loves INT DEFAULT 0 COMMENT '喜爱量/点赞数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'pending' COMMENT '审核状态：pending/approved/rejected',
    audit_time TIMESTAMP NULL COMMENT '审核时间',
    reject_reason TEXT COMMENT '驳回原因',
    submitter_id INT COMMENT '提交用户ID',
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_submitter (submitter_id),
    FOREIGN KEY (submitter_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工具表';

-- 工具图片表
CREATE TABLE IF NOT EXISTS tool_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tool_id INT NOT NULL COMMENT '工具ID',
    image_url VARCHAR(500) NOT NULL COMMENT '图片URL',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tool_id (tool_id),
    FOREIGN KEY (tool_id) REFERENCES tools(resource_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工具图片表';

-- 工具标签表
CREATE TABLE IF NOT EXISTS tool_tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tool_id INT NOT NULL COMMENT '工具ID',
    tag VARCHAR(50) NOT NULL COMMENT '标签名称',
    INDEX idx_tool_id (tool_id),
    INDEX idx_tag (tag),
    FOREIGN KEY (tool_id) REFERENCES tools(resource_id) ON DELETE CASCADE,
    UNIQUE KEY uk_tool_tag (tool_id, tag)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工具标签表';

-- 工具贡献者表
CREATE TABLE IF NOT EXISTS tool_contributors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tool_id INT NOT NULL COMMENT '工具ID',
    user_id INT NOT NULL COMMENT '用户ID',
    INDEX idx_tool_id (tool_id),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (tool_id) REFERENCES tools(resource_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_tool_user (tool_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工具贡献者表';

-- ==================== 课程相关表 ====================

-- 课程表
CREATE TABLE IF NOT EXISTS courses (
    course_id INT AUTO_INCREMENT PRIMARY KEY,
    resource_type VARCHAR(50) DEFAULT 'course' COMMENT '资源类型',
    name VARCHAR(255) NOT NULL COMMENT '课程名称',
    semester VARCHAR(50) COMMENT '学期',
    credit INT COMMENT '学分',
    cover VARCHAR(500) COMMENT '课程封面图',
    views INT DEFAULT 0 COMMENT '浏览量',
    loves INT DEFAULT 0 COMMENT '点赞数',
    collections INT DEFAULT 0 COMMENT '收藏量',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_semester (semester),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程表';

-- 课程教师表
CREATE TABLE IF NOT EXISTS course_teachers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    course_id INT NOT NULL COMMENT '课程ID',
    teacher_name VARCHAR(100) NOT NULL COMMENT '教师姓名',
    INDEX idx_course_id (course_id),
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程教师表';

-- 课程分类表
CREATE TABLE IF NOT EXISTS course_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    course_id INT NOT NULL COMMENT '课程ID',
    category VARCHAR(100) NOT NULL COMMENT '分类名称',
    INDEX idx_course_id (course_id),
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程分类表';

-- 课程资源表（URL资源）
CREATE TABLE IF NOT EXISTS course_resources_web (
    resource_id INT AUTO_INCREMENT PRIMARY KEY,
    course_id INT NOT NULL COMMENT '课程ID',
    resource_intro VARCHAR(255) NOT NULL COMMENT '资源说明',
    resource_url VARCHAR(500) NOT NULL COMMENT '资源网址',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_course_id (course_id),
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程URL资源表';

-- 课程资源表（上传资源/课本）
CREATE TABLE IF NOT EXISTS course_resources_upload (
    resource_id INT AUTO_INCREMENT PRIMARY KEY,
    course_id INT NOT NULL COMMENT '课程ID',
    resource_intro VARCHAR(255) NOT NULL COMMENT '资源说明',
    resource_upload VARCHAR(500) NOT NULL COMMENT '上传文件URL',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_course_id (course_id),
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程上传资源表';

-- 课程贡献者表
CREATE TABLE IF NOT EXISTS course_contributors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    course_id INT NOT NULL COMMENT '课程ID',
    user_id INT NOT NULL COMMENT '用户ID',
    INDEX idx_course_id (course_id),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_course_user (course_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程贡献者表';

-- ==================== 项目相关表 ====================

-- 项目表
CREATE TABLE IF NOT EXISTS projects (
    project_id INT AUTO_INCREMENT PRIMARY KEY,
    resource_type VARCHAR(50) DEFAULT 'project' COMMENT '资源类型',
    name VARCHAR(255) NOT NULL UNIQUE COMMENT '项目名称',
    description VARCHAR(500) COMMENT '项目简介',
    detail TEXT COMMENT '项目详细介绍（支持markdown）',
    github_url VARCHAR(500) COMMENT 'Github仓库链接',
    category VARCHAR(100) COMMENT '项目类别：实训项目/课程设计等',
    cover VARCHAR(500) COMMENT '项目封面图',
    views INT DEFAULT 0 COMMENT '浏览量',
    loves INT DEFAULT 0 COMMENT '点赞数',
    collections INT DEFAULT 0 COMMENT '收藏量',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目表';

-- 项目技术栈表
CREATE TABLE IF NOT EXISTS project_tech_stack (
    id INT AUTO_INCREMENT PRIMARY KEY,
    project_id INT NOT NULL COMMENT '项目ID',
    tech VARCHAR(50) NOT NULL COMMENT '技术栈名称',
    INDEX idx_project_id (project_id),
    FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    UNIQUE KEY uk_project_tech (project_id, tech)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目技术栈表';

-- 项目图片表
CREATE TABLE IF NOT EXISTS project_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    project_id INT NOT NULL COMMENT '项目ID',
    image_url VARCHAR(500) NOT NULL COMMENT '图片URL',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_project_id (project_id),
    FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目图片表';

-- 项目作者表
CREATE TABLE IF NOT EXISTS project_authors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    project_id INT NOT NULL COMMENT '项目ID',
    user_id INT NOT NULL COMMENT '用户ID',
    INDEX idx_project_id (project_id),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_project_user (project_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目作者表';

-- ==================== 评论相关表 ====================

-- 评论表（通用，用于工具/课程/项目）
CREATE TABLE IF NOT EXISTS comments (
    comment_id INT AUTO_INCREMENT PRIMARY KEY,
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型：tool/course/project',
    resource_id INT NOT NULL COMMENT '资源ID（工具ID/课程ID/项目ID）',
    parent_id INT NULL COMMENT '父评论ID（用于回复）',
    user_id INT NOT NULL COMMENT '评论用户ID',
    content TEXT NOT NULL COMMENT '评论内容（不超过800字）',
    love_count INT DEFAULT 0 COMMENT '点赞数',
    reply_total INT DEFAULT 0 COMMENT '回复总数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '评论时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL COMMENT '删除时间（软删除）',
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_user_id (user_id),
    INDEX idx_parent_id (parent_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(comment_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- 评论点赞表
CREATE TABLE IF NOT EXISTS comment_likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    comment_id INT NOT NULL COMMENT '评论ID',
    user_id INT NOT NULL COMMENT '用户ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_comment_id (comment_id),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_comment_user (comment_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论点赞表';

-- ==================== 用户行为表 ====================

-- 收藏表
CREATE TABLE IF NOT EXISTS collections (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL COMMENT '用户ID',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型：tool/course/project',
    resource_id INT NOT NULL COMMENT '资源ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_resource (resource_type, resource_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_resource (user_id, resource_type, resource_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏表';

-- 点赞表
CREATE TABLE IF NOT EXISTS likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL COMMENT '用户ID',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型：tool/course/project',
    resource_id INT NOT NULL COMMENT '资源ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_resource (resource_type, resource_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_resource (user_id, resource_type, resource_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞表';

-- ==================== 审核/状态管理表 ====================

-- 资源状态变更记录表
CREATE TABLE IF NOT EXISTS resource_status_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型',
    resource_id INT NOT NULL COMMENT '资源ID',
    old_status VARCHAR(50) COMMENT '原状态',
    new_status VARCHAR(50) NOT NULL COMMENT '新状态',
    operator_id INT COMMENT '操作用户ID',
    operate_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    INDEX idx_resource (resource_type, resource_id),
    FOREIGN KEY (operator_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资源状态变更记录表';

-- ==================== 初始化数据 ====================

-- 插入一个管理员用户（密码需要在使用时设置）
-- INSERT INTO users (username, nickname, email, password, role) 
-- VALUES ('admin', '管理员', 'admin@example.com', '需要设置密码', 'admin');

