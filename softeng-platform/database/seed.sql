-- seed.sql: 插入示例数据（对应之前代码里的“死数据/示例返回”）
-- 使用方法：在 MySQL 客户端里执行本文件（确保 schema.sql 已执行过）

USE softeng;

-- 为了可重复执行：先清空（空库首次执行也没问题）
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE tool_contributors;
TRUNCATE TABLE tool_tags;
TRUNCATE TABLE tool_images;
TRUNCATE TABLE tools;

TRUNCATE TABLE course_contributors;
TRUNCATE TABLE course_resources_upload;
TRUNCATE TABLE course_resources_web;
TRUNCATE TABLE course_categories;
TRUNCATE TABLE course_teachers;
TRUNCATE TABLE courses;

TRUNCATE TABLE project_authors;
TRUNCATE TABLE project_images;
TRUNCATE TABLE project_tech_stack;
TRUNCATE TABLE projects;

TRUNCATE TABLE users;
SET FOREIGN_KEY_CHECKS = 1;

-- 所有用户统一密码：123456（bcrypt 哈希）
-- 你也可以后续通过 /auth/login 之类接口验证（注意：你项目目前 JWT 未完全打通）
SET @PWD_123456 = '$2a$10$oHhSw3NUtFzQB1JFK7Tw4.OHCS1kjKXYQ20cdGsP96myzM3PUPKAe';

-- ===================== users =====================
INSERT INTO users (id, username, nickname, email, password, avatar, description, face_photo, role)
VALUES
  (1, 'user1', '用户1', 'user1@example.com', @PWD_123456, 'https://example.com/avatar1.jpg', '我是 user1', 'https://example.com/cover1.jpg', 'user'),
  (2, 'user2', '用户2', 'user2@example.com', @PWD_123456, 'https://example.com/avatar2.jpg', '我是 user2', 'https://example.com/cover2.jpg', 'user'),
  (3, 'user3', '用户3', 'user3@example.com', @PWD_123456, 'https://example.com/avatar3.jpg', '我是 user3', 'https://example.com/cover3.jpg', 'user'),
  (4, 'admin', '管理员', 'admin@example.com', @PWD_123456, NULL, NULL, NULL, 'admin');

-- ===================== tools =====================
INSERT INTO tools (
  resource_id, resource_type, resource_name, resource_link, description, description_detail,
  category, views, collections, loves, status, submitter_id
) VALUES
  (1, 'tool', 'Sample Tool', 'https://example.com',
   'This is a sample tool', 'Detailed description of the tool',
   '软件开发', 100, 50, 25, 'approved', 1),
  (2, 'tool', 'Search Result Tool', 'https://example.com/search',
   'This tool matches the search', 'This tool is used for search demo',
   '论文阅读', 80, 40, 20, 'approved', 3);

INSERT INTO tool_images (tool_id, image_url, sort_order) VALUES
  (1, 'https://example.com/image.jpg', 0),
  (1, 'https://example.com/image1.jpg', 1),
  (1, 'https://example.com/image2.jpg', 2),
  (2, 'https://example.com/search.jpg', 0);

INSERT INTO tool_tags (tool_id, tag) VALUES
  (1, '免费'),
  (1, 'AI工具'),
  (1, '高效'),
  (2, '搜索'),
  (2, '工具');

INSERT INTO tool_contributors (tool_id, user_id) VALUES
  (1, 1),
  (1, 2),
  (2, 3);

-- ===================== courses =====================
INSERT INTO courses (
  course_id, resource_type, name, semester, credit, cover, views, loves, collections
) VALUES
  (1, 'course', '软件工程导论', '大二上', 3, 'https://example.com/course1.jpg', 1000, 200, 150),
  (2, 'course', '高级软件工程', '大三上', 2, 'https://example.com/course2.jpg', 800, 150, 100);

INSERT INTO course_teachers (course_id, teacher_name) VALUES
  (1, '张教授'),
  (1, '李教授'),
  (2, '王教授');

INSERT INTO course_categories (course_id, category) VALUES
  (1, '专必'),
  (1, '有签到'),
  (2, '专选'),
  (2, '无签到');

INSERT INTO course_resources_web (course_id, resource_intro, resource_url, sort_order) VALUES
  (1, '课程视频', 'https://example.com/video1', 0);

INSERT INTO course_resources_upload (course_id, resource_intro, resource_upload, sort_order) VALUES
  (1, '课程讲义', 'https://example.com/lecture1.pdf', 0);

INSERT INTO course_contributors (course_id, user_id) VALUES
  (1, 1),
  (1, 2),
  (2, 3);

-- ===================== projects =====================
INSERT INTO projects (
  project_id, resource_type, name, description, detail, github_url, category, cover, views, loves, collections
) VALUES
  (1, 'project', '校园社交平台', '基于Go和React的校园社交平台',
   '这是一个示例项目（支持markdown）。', 'https://github.com/example/social-platform',
   '实训项目', 'https://example.com/project1.jpg', 200, 45, 30),
  (2, 'project', '电商平台', '完整的电商平台解决方案',
   '这是另一个示例项目（支持markdown）。', 'https://github.com/example/ecommerce-platform',
   '课程设计', 'https://example.com/project2.jpg', 150, 35, 25);

INSERT INTO project_tech_stack (project_id, tech) VALUES
  (1, 'Go'),
  (1, 'React'),
  (1, 'PostgreSQL'),
  (2, 'Java'),
  (2, 'Spring Boot'),
  (2, 'Vue');

INSERT INTO project_images (project_id, image_url, sort_order) VALUES
  (1, 'https://example.com/image1.jpg', 0),
  (1, 'https://example.com/image2.jpg', 1),
  (2, 'https://example.com/project2_img1.jpg', 0);

INSERT INTO project_authors (project_id, user_id) VALUES
  (1, 1),
  (1, 2),
  (2, 3);

-- 完成提示：现在访问 /tools/profile /courses/profile /projects/profile 应该能看到真实数据了

