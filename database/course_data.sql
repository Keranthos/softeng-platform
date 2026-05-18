-- 插入课程路线数据
-- 注意：course_id 如果已存在则跳过，使用 INSERT IGNORE 避免重复插入

-- 先检查并添加 code 字段（如果不存在）
-- MySQL不支持 IF NOT EXISTS，需要先检查
SET @dbname = DATABASE();
SET @tablename = 'courses';
SET @columnname = 'code';
SET @preparedStatement = (SELECT IF(
  (
    SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE
      (TABLE_SCHEMA = @dbname)
      AND (TABLE_NAME = @tablename)
      AND (COLUMN_NAME = @columnname)
  ) > 0,
  'SELECT 1',
  CONCAT('ALTER TABLE ', @tablename, ' ADD COLUMN ', @columnname, ' VARCHAR(50)')
));
PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

-- 插入课程数据
INSERT IGNORE INTO courses (course_id, name, code, semester, credit, resource_type, cover, views, loves, collections) VALUES
-- 大一上学期
(101, '高等数学 I', 'MATH1001', '1-1', 5, 'course', 'https://placehold.co/600x450/3b82f6/ffffff?text=MATH', 0, 45, 12),
(102, '程序设计基础', 'CS1001', '1-1', 4, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 102, 28),
(103, '思想道德修养', 'POLIO1001', '1-1', 2, 'course', 'https://placehold.co/600x450/a855f7/ffffff?text=POLIO', 0, 10, 5),
(104, '当代文化研究', 'PUB1001', '1-1', 2, 'course', 'https://placehold.co/600x450/f97316/ffffff?text=PUB', 0, 80, 5),

-- 大一下学期
(201, '高等数学 II', 'MATH1002', '1-2', 5, 'course', 'https://placehold.co/600x450/3b82f6/ffffff?text=MATH', 0, 38, 15),
(202, '线性代数', 'MATH1003', '1-2', 3, 'course', 'https://placehold.co/600x450/3b82f6/ffffff?text=MATH', 0, 88, 20),
(203, '离散数学', 'CS1002', '1-2', 4, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 150, 35),
(204, '体育2', 'PE1002', '1-2', 1, 'course', 'https://placehold.co/600x450/ef4444/ffffff?text=PE', 0, 180, 5),

-- 大二上学期
(301, '数据结构与算法', 'CS2001', '2-1', 5, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 230, 56),
(302, '计算机组成原理', 'CS2002', '2-1', 4, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 95, 30),
(303, 'Python应用开发', 'CS2005', '2-1', 2, 'course', 'https://placehold.co/600x450/14b8a6/ffffff?text=PYTHON', 0, 67, 18),

-- 大二下学期
(401, '操作系统', 'CS2003', '2-2', 4, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 180, 42),
(402, '计算机网络', 'CS2004', '2-2', 4, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 160, 38),

-- 大三上学期
(501, '计算机网络', 'CS3001', '3-1', 3, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 120, 18),
(502, '数据库系统概论', 'CS3002', '3-1', 3, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 140, 25),

-- 大三下学期
(601, '软件工程导论', 'CS3003', '3-2', 2, 'course', 'https://placehold.co/600x450/14b8a6/ffffff?text=SE', 0, 80, 22),

-- 大四上学期
(701, '人工智能导论', 'CS4001', '4-1', 2, 'course', 'https://placehold.co/600x450/14b8a6/ffffff?text=AI', 0, 90, 15),

-- 大四下学期
(702, '毕业设计', 'CS4002', '4-2', 6, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=GRAD', 0, 50, 10);

-- 插入教师信息到 course_teachers 表
INSERT IGNORE INTO course_teachers (course_id, teacher_name) VALUES
(101, '张老师'),
(102, '李老师'),
(103, '王老师'),
(104, '张老师'),
(201, '张老师'),
(202, '赵老师'),
(203, '钱老师'),
(204, '张老师'),
(301, '孙老师'),
(302, '周老师'),
(303, '吴老师'),
(401, '郑老师'),
(402, '冯老师'),
(501, '马老师'),
(502, '刘老师'),
(601, '毛老师'),
(701, '林老师'),
(702, '何老师');

-- 插入课程类别到 course_categories 表
INSERT IGNORE INTO course_categories (course_id, category) VALUES
(101, '公必'),
(102, '专必'),
(103, '公必'),
(104, '公选'),
(201, '公必'),
(202, '公必'),
(203, '专必'),
(204, '公必'),
(301, '专必'),
(302, '专必'),
(303, '专选'),
(401, '专必'),
(402, '专必'),
(501, '专必'),
(502, '专必'),
(601, '专选'),
(701, '专选'),
(702, '专必');

