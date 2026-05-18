-- 软件工程平台测试数据
-- 适合软件工程大三学生使用

USE softeng;

-- ==================== 用户数据 ====================
-- 注意：密码都是 123456，使用bcrypt加密
-- 实际使用时请修改密码

-- 管理员用户（如果已存在则跳过）
-- 注意：所有用户密码都是 123456，使用bcrypt加密
INSERT IGNORE INTO users (username, nickname, email, password, avatar, description, role, created_at) VALUES
('admin', '管理员', 'admin@softeng.edu.cn', '$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O', 'https://via.placeholder.com/150?text=Admin', '平台管理员，负责审核和运维', 'admin', '2024-01-01 08:00:00');

-- 教师用户（如果已存在则跳过）
INSERT IGNORE INTO users (username, nickname, email, password, avatar, description, role, created_at) VALUES
('teacher_zhang', '张教授', 'zhang@softeng.edu.cn', '$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O', 'https://via.placeholder.com/150?text=Teacher', '软件工程系教授，主讲软件工程导论', 'user', '2024-01-02 09:00:00'),
('teacher_li', '李老师', 'li@softeng.edu.cn', '$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O', 'https://via.placeholder.com/150?text=Teacher2', '数据库系统课程教师', 'user', '2024-01-02 10:00:00'),
('teacher_wang', '王老师', 'wang@softeng.edu.cn', '$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O', 'https://via.placeholder.com/150?text=Teacher3', 'Web开发技术课程教师', 'user', '2024-01-02 11:00:00');

-- 学生用户（软件工程大三学生）（如果已存在则跳过）
INSERT IGNORE INTO users (username, nickname, email, password, avatar, description, role, created_at) VALUES
('student_zhang', '张小三', 'zhangxiaosan@student.edu.cn', '$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O', 'https://via.placeholder.com/150?text=Student1', '软件工程大三学生，热爱前后端开发', 'user', '2024-01-10 08:00:00'),
('student_li', '李四', 'lisi@student.edu.cn', '$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O', 'https://via.placeholder.com/150?text=Student2', '软件工程大三学生，专注后端开发', 'user', '2024-01-10 09:00:00'),
('student_wang', '王五', 'wangwu@student.edu.cn', '$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O', 'https://via.placeholder.com/150?text=Student3', '软件工程大三学生，前端开发爱好者', 'user', '2024-01-10 10:00:00'),
('student_zhao', '赵六', 'zhaoliu@student.edu.cn', '$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O', 'https://via.placeholder.com/150?text=Student4', '软件工程大三学生，全栈开发', 'user', '2024-01-10 11:00:00'),
('student_sun', '孙七', 'sunqi@student.edu.cn', '$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O', 'https://via.placeholder.com/150?text=Student5', '软件工程大三学生，AI方向', 'user', '2024-01-10 12:00:00');

-- ==================== 工具数据 ====================

-- 已审核通过的工具
-- 注意：submitter_id对应实际用户ID：16=student_zhang, 17=student_li, 18=student_wang, 19=student_zhao, 20=student_sun
INSERT INTO tools (resource_type, resource_name, resource_link, description, description_detail, category, views, collections, loves, status, audit_time, submitter_id, created_at) VALUES
('tool', 'VS Code', 'https://code.visualstudio.com/', '微软开发的免费代码编辑器', 'Visual Studio Code是一个轻量级但功能强大的源代码编辑器，支持多种编程语言和丰富的插件生态。适合前端、后端、全栈开发使用。', '开发工具', 1250, 89, 156, 'approved', '2024-01-15 10:00:00', 17, '2024-01-14 09:00:00'),
('tool', 'Postman', 'https://www.postman.com/', 'API开发和测试工具', 'Postman是一个功能强大的API测试工具，支持REST、GraphQL等API类型。可以快速测试接口、管理API文档、进行自动化测试。', '测试工具', 980, 67, 112, 'approved', '2024-01-16 10:00:00', 18, '2024-01-15 10:00:00'),
('tool', 'Git', 'https://git-scm.com/', '分布式版本控制系统', 'Git是最流行的版本控制系统，帮助开发者管理代码版本、协同开发、分支管理。是软件工程必备工具。', '版本控制', 2100, 145, 198, 'approved', '2024-01-17 10:00:00', 19, '2024-01-16 11:00:00'),
('tool', 'MySQL Workbench', 'https://www.mysql.com/products/workbench/', 'MySQL数据库管理工具', 'MySQL官方提供的可视化数据库管理工具，支持数据库设计、SQL开发、服务器配置等功能。适合数据库课程学习使用。', '数据库工具', 750, 52, 78, 'approved', '2024-01-18 10:00:00', 20, '2024-01-17 14:00:00'),
('tool', 'Docker', 'https://www.docker.com/', '容器化平台', 'Docker是一个开源的应用容器引擎，可以轻松打包、部署和运行应用程序。是现代软件开发的重要工具。', '部署工具', 650, 45, 89, 'approved', '2024-01-19 10:00:00', 16, '2024-01-18 15:00:00'),
('tool', 'Figma', 'https://www.figma.com/', '协作式界面设计工具', 'Figma是一个基于浏览器的协作界面设计工具，支持实时协作、原型设计、组件库管理。适合前端开发和UI设计。', '设计工具', 580, 38, 65, 'approved', '2024-01-20 10:00:00', 17, '2024-01-19 16:00:00'),
('tool', 'JetBrains IDEA', 'https://www.jetbrains.com/idea/', 'Java开发IDE', 'IntelliJ IDEA是功能强大的Java集成开发环境，支持Spring、Hibernate等框架，智能代码提示和重构。', '开发工具', 890, 61, 95, 'approved', '2024-01-21 10:00:00', 18, '2024-01-20 10:00:00'),
('tool', 'Chrome DevTools', 'https://developer.chrome.com/docs/devtools/', '浏览器开发者工具', 'Chrome内置的开发者工具，用于调试JavaScript、分析性能、检查网络请求等。前端开发必备。', '调试工具', 1100, 78, 134, 'approved', '2024-01-22 10:00:00', 19, '2024-01-21 11:00:00');

-- 待审核的工具
INSERT INTO tools (resource_type, resource_name, resource_link, description, description_detail, category, views, collections, loves, status, submitter_id, created_at) VALUES
('tool', 'Swagger UI', 'https://swagger.io/tools/swagger-ui/', 'API文档生成工具', 'Swagger UI可以自动生成美观的API交互式文档，支持在线测试接口。适合API开发项目使用。', '文档工具', 0, 0, 0, 'pending', 17, NOW()),
('tool', 'Jira', 'https://www.atlassian.com/software/jira', '项目管理工具', 'Jira是Atlassian开发的项目管理和问题跟踪工具，支持敏捷开发、看板、任务分配等功能。', '项目管理', 0, 0, 0, 'pending', 18, NOW()),
('tool', 'Redis Desktop Manager', 'https://redisdesktop.com/', 'Redis可视化工具', 'Redis Desktop Manager是一个Redis数据库的可视化管理工具，方便查看和管理Redis数据。', '数据库工具', 0, 0, 0, 'pending', 19, NOW());

-- 工具标签 (使用实际的工具ID：17=VS Code, 18=Postman, 19=Git, 20=MySQL Workbench, 21=Docker, 22=Figma, 23=JetBrains IDEA, 24=Chrome DevTools, 25=Swagger UI, 26=Jira, 27=Redis Desktop Manager)
INSERT IGNORE INTO tool_tags (tool_id, tag) VALUES
(17, '编辑器'), (17, '前端'), (17, '后端'), (17, '免费'),
(18, 'API'), (18, '测试'), (18, '开发'),
(19, '版本控制'), (19, 'Git'), (19, '必备'),
(20, '数据库'), (20, 'MySQL'), (20, '可视化'),
(21, '容器'), (21, 'Docker'), (21, '部署'),
(22, '设计'), (22, 'UI'), (22, '协作'),
(23, 'IDE'), (23, 'Java'), (23, '企业级'),
(24, '调试'), (24, '前端'), (24, '浏览器'),
(25, '文档'), (25, 'API'),
(26, '项目管理'),
(27, '数据库');

-- 工具贡献者 (user_id: 16-20为学生)
INSERT IGNORE INTO tool_contributors (tool_id, user_id) VALUES
(17, 17), (18, 18), (19, 19), (20, 20), (21, 16), (22, 17), (23, 18), (24, 19),
(25, 17), (26, 18), (27, 19);

-- ==================== 课程数据 ====================

INSERT INTO courses (resource_type, name, semester, credit, cover, views, loves, collections, created_at) VALUES
('course', '软件工程导论', '2024春季', 3, 'https://via.placeholder.com/400x300?text=软件工程导论', 520, 45, 38, '2024-02-01 08:00:00'),
('course', '数据库系统原理', '2024春季', 3, 'https://via.placeholder.com/400x300?text=数据库系统原理', 480, 42, 35, '2024-02-02 09:00:00'),
('course', 'Web开发技术', '2024春季', 2, 'https://via.placeholder.com/400x300?text=Web开发技术', 680, 58, 52, '2024-02-03 10:00:00'),
('course', '软件项目管理', '2024春季', 2, 'https://via.placeholder.com/400x300?text=软件项目管理', 390, 32, 28, '2024-02-04 11:00:00'),
('course', '软件测试与质量保证', '2024春季', 2, 'https://via.placeholder.com/400x300?text=软件测试', 350, 28, 25, '2024-02-05 12:00:00'),
('course', '数据结构与算法', '2023秋季', 4, 'https://via.placeholder.com/400x300?text=数据结构', 950, 78, 65, '2023-09-01 08:00:00'),
('course', '操作系统', '2023秋季', 3, 'https://via.placeholder.com/400x300?text=操作系统', 720, 56, 48, '2023-09-02 09:00:00'),
('course', '计算机网络', '2023秋季', 3, 'https://via.placeholder.com/400x300?text=计算机网络', 650, 52, 45, '2023-09-03 10:00:00');

-- 待审核的课程（courses表没有status字段，这些课程也会正常显示）
INSERT INTO courses (resource_type, name, semester, credit, cover, views, loves, collections, created_at) VALUES
('course', '移动应用开发', '2024秋季', 2, 'https://via.placeholder.com/400x300?text=移动开发', 0, 0, 0, NOW()),
('course', '人工智能基础', '2024秋季', 3, 'https://via.placeholder.com/400x300?text=人工智能', 0, 0, 0, NOW());

-- 课程教师
INSERT INTO course_teachers (course_id, teacher_name) VALUES
(1, '张教授'), (1, '李老师'),
(2, '李老师'),
(3, '王老师'),
(4, '张教授'),
(5, '李老师'),
(6, '张教授'),
(7, '王老师'),
(8, '张教授'),
(9, '王老师'),
(10, '李老师');

-- 课程分类
INSERT INTO course_categories (course_id, category) VALUES
(1, '专业必修'), (1, '软件工程'),
(2, '专业必修'), (2, '数据库'),
(3, '专业选修'), (3, 'Web开发'),
(4, '专业选修'), (4, '项目管理'),
(5, '专业选修'), (5, '软件测试'),
(6, '专业必修'), (6, '算法'),
(7, '专业必修'), (7, '系统'),
(8, '专业必修'), (8, '网络'),
(9, '专业选修'), (9, '移动开发'),
(10, '专业选修'), (10, '人工智能');

-- 课程贡献者 (user_id: 13-15为教师，course_id需要根据实际插入后的ID调整)
-- 使用子查询根据课程名称查找实际的course_id
INSERT IGNORE INTO course_contributors (course_id, user_id) 
SELECT c.course_id, 13 FROM courses c WHERE c.name = '软件工程导论'
UNION ALL SELECT c.course_id, 13 FROM courses c WHERE c.name = '软件工程导论'
UNION ALL SELECT c.course_id, 14 FROM courses c WHERE c.name = '数据库系统原理'
UNION ALL SELECT c.course_id, 15 FROM courses c WHERE c.name = 'Web开发技术'
UNION ALL SELECT c.course_id, 13 FROM courses c WHERE c.name = '软件项目管理'
UNION ALL SELECT c.course_id, 14 FROM courses c WHERE c.name = '软件测试与质量保证'
UNION ALL SELECT c.course_id, 13 FROM courses c WHERE c.name = '数据结构与算法'
UNION ALL SELECT c.course_id, 15 FROM courses c WHERE c.name = '操作系统'
UNION ALL SELECT c.course_id, 13 FROM courses c WHERE c.name = '计算机网络'
UNION ALL SELECT c.course_id, 15 FROM courses c WHERE c.name = '移动应用开发'
UNION ALL SELECT c.course_id, 14 FROM courses c WHERE c.name = '人工智能基础';

-- 课程资源（URL资源）
INSERT INTO course_resources_web (course_id, resource_intro, resource_url, sort_order) VALUES
(1, '软件工程概述视频', 'https://www.bilibili.com/video/BV1example1', 1),
(1, '需求分析课程', 'https://www.bilibili.com/video/BV1example2', 2),
(2, 'MySQL基础教程', 'https://www.bilibili.com/video/BV1example3', 1),
(2, 'SQL高级查询', 'https://www.bilibili.com/video/BV1example4', 2),
(3, 'HTML/CSS基础', 'https://www.bilibili.com/video/BV1example5', 1),
(3, 'JavaScript进阶', 'https://www.bilibili.com/video/BV1example6', 2),
(3, 'Vue.js框架', 'https://www.bilibili.com/video/BV1example7', 3);

-- 课程资源（上传资源）
INSERT INTO course_resources_upload (course_id, resource_intro, resource_upload, sort_order) VALUES
(1, '软件工程导论课件PPT', 'https://example.com/resources/software-engineering-ppt.pdf', 1),
(1, '需求规格说明文档模板', 'https://example.com/resources/requirement-template.docx', 2),
(2, '数据库设计实验指导书', 'https://example.com/resources/database-lab-guide.pdf', 1),
(2, 'MySQL练习题集', 'https://example.com/resources/mysql-exercises.pdf', 2),
(3, 'Web开发项目案例', 'https://example.com/resources/web-project-case.zip', 1);

-- ==================== 项目数据 ====================

INSERT INTO projects (resource_type, name, description, detail, github_url, category, cover, views, loves, collections, created_at) VALUES
('project', '在线学习平台', '基于Vue3和Go开发的在线学习管理系统，支持课程管理、作业提交、在线考试等功能。', '# 项目简介\n这是一个完整的在线学习平台，采用前后端分离架构。\n\n## 技术栈\n- 前端：Vue3 + Element Plus + Axios\n- 后端：Go + Gin + MySQL\n- 部署：Docker + Nginx\n\n## 功能特性\n1. 用户管理（学生、教师、管理员）\n2. 课程管理（创建、编辑、发布）\n3. 作业系统（布置、提交、批改）\n4. 在线考试系统\n5. 成绩统计和导出', 'https://github.com/example/learning-platform', '实训项目', 'https://via.placeholder.com/600x400?text=在线学习平台', 320, 28, 35, '2024-02-15 10:00:00'),
('project', '校园二手交易平台', '基于React和Node.js开发的校园二手商品交易平台，支持商品发布、搜索、聊天等功能。', '# 项目简介\n为校园学生提供便捷的二手商品交易服务。\n\n## 技术栈\n- 前端：React + Ant Design + Redux\n- 后端：Node.js + Express + MongoDB\n- 实时通信：Socket.io\n\n## 功能特性\n1. 商品发布和管理\n2. 智能搜索和分类\n3. 站内消息系统\n4. 用户信用评价', 'https://github.com/example/campus-market', '课程设计', 'https://via.placeholder.com/600x400?text=二手交易平台', 280, 24, 30, '2024-02-16 11:00:00'),
('project', '个人博客系统', '基于Spring Boot和Vue3开发的个人博客系统，支持Markdown编辑、评论、标签管理。', '# 项目简介\n功能完整的个人博客系统，适合技术分享和个人记录。\n\n## 技术栈\n- 前端：Vue3 + Vite + Element Plus\n- 后端：Spring Boot + MyBatis + MySQL\n- 搜索引擎：Elasticsearch\n\n## 功能特性\n1. Markdown编辑器\n2. 文章分类和标签\n3. 评论系统\n4. 全文搜索', 'https://github.com/example/personal-blog', '个人项目', 'https://via.placeholder.com/600x400?text=个人博客', 450, 38, 42, '2024-02-17 12:00:00'),
('project', '任务管理系统', '基于Python Django和Vue开发的团队协作任务管理工具，支持看板、甘特图、时间跟踪。', '# 项目简介\n敏捷开发团队协作工具，类似Trello和Jira。\n\n## 技术栈\n- 前端：Vue3 + TypeScript + Vuetify\n- 后端：Django + Django REST Framework\n- 数据库：PostgreSQL\n\n## 功能特性\n1. 看板式任务管理\n2. 甘特图项目规划\n3. 团队协作和权限管理\n4. 时间跟踪和报告', 'https://github.com/example/task-manager', '实训项目', 'https://via.placeholder.com/600x400?text=任务管理', 390, 32, 38, '2024-02-18 13:00:00'),
('project', '图书管理系统', '基于Java Swing开发的桌面图书管理系统，适合数据库课程设计使用。', '# 项目简介\n经典的图书管理系统，功能完整，适合课程设计。\n\n## 技术栈\n- 前端：Java Swing\n- 后端：Java + JDBC\n- 数据库：MySQL\n\n## 功能特性\n1. 图书信息管理\n2. 借阅和归还\n3. 读者管理\n4. 统计报表', 'https://github.com/example/library-system', '课程设计', 'https://via.placeholder.com/600x400?text=图书管理', 520, 45, 48, '2024-02-19 14:00:00'),
('project', '在线聊天室', '基于WebSocket的实时聊天应用，支持多房间、私聊、文件传输。', '# 项目简介\n实时在线聊天应用，支持多人在线聊天。\n\n## 技术栈\n- 前端：React + Socket.io-client\n- 后端：Node.js + Express + Socket.io\n- 数据库：Redis（消息队列）\n\n## 功能特性\n1. 多房间聊天\n2. 私聊功能\n3. 文件上传和传输\n4. 用户状态显示', 'https://github.com/example/chat-room', '个人项目', 'https://via.placeholder.com/600x400?text=聊天室', 210, 18, 22, '2024-02-20 15:00:00');

-- 待审核的项目（projects表没有status字段，这些项目也会正常显示）
INSERT IGNORE INTO projects (resource_type, name, description, detail, github_url, category, cover, views, loves, collections, created_at) VALUES
('project', '在线考试系统', '基于Spring Cloud微服务架构的在线考试平台，支持多种题型、自动判卷、防作弊。', '微服务架构的在线考试系统，包含用户服务、考试服务、题库服务等。', 'https://github.com/example/exam-system', '实训项目', 'https://via.placeholder.com/600x400?text=考试系统', 0, 0, 0, NOW()),
('project', '智能排课系统', '基于遗传算法的智能排课系统，解决高校课程安排问题。', '使用遗传算法优化课程安排，避免时间冲突，合理分配教室资源。', 'https://github.com/example/course-scheduling', '课程设计', 'https://via.placeholder.com/600x400?text=排课系统', 0, 0, 0, NOW());

-- 项目技术栈
INSERT INTO project_tech_stack (project_id, tech) VALUES
(1, 'Vue3'), (1, 'Go'), (1, 'Gin'), (1, 'MySQL'), (1, 'Docker'),
(2, 'React'), (2, 'Node.js'), (2, 'Express'), (2, 'MongoDB'), (2, 'Socket.io'),
(3, 'Vue3'), (3, 'Spring Boot'), (3, 'MyBatis'), (3, 'MySQL'), (3, 'Elasticsearch'),
(4, 'Vue3'), (4, 'Django'), (4, 'PostgreSQL'), (4, 'TypeScript'),
(5, 'Java'), (5, 'Swing'), (5, 'JDBC'), (5, 'MySQL'),
(6, 'React'), (6, 'Node.js'), (6, 'Socket.io'), (6, 'Redis'),
(7, 'Spring Cloud'), (7, 'Java'), (7, 'MySQL'), (7, 'Redis'),
(8, 'Python'), (8, '遗传算法'), (8, 'Flask'), (8, 'MySQL');

-- 项目作者 (user_id: 16-20为学生)
INSERT INTO project_authors (project_id, user_id) VALUES
(1, 17), (1, 18),
(2, 19), (2, 20),
(3, 16),
(4, 17), (4, 18),
(5, 19),
(6, 20),
(7, 16),
(8, 17);

-- ==================== 评论数据 ====================

-- 工具评论 (user_id: 16-20为学生)
INSERT INTO comments (resource_type, resource_id, user_id, content, love_count, reply_total, created_at) VALUES
('tool', 1, 17, 'VS Code真的太好用了！插件生态丰富，界面简洁，强烈推荐！', 12, 2, '2024-01-20 10:00:00'),
('tool', 1, 18, '正在使用中，代码提示很智能，调试功能也很强大。', 8, 1, '2024-01-20 14:30:00'),
('tool', 2, 19, 'Postman对API测试来说必不可少，界面友好，功能强大。', 10, 1, '2024-01-21 09:15:00'),
('tool', 3, 20, 'Git是程序员必备技能，版本控制的核心工具！', 15, 3, '2024-01-22 11:20:00');

-- 课程评论
INSERT INTO comments (resource_type, resource_id, user_id, content, love_count, reply_total, created_at) VALUES
('course', 1, 17, '张教授讲课很清晰，软件工程的基础知识讲得很透彻！', 18, 2, '2024-02-10 09:00:00'),
('course', 1, 18, '课程内容很实用，对实际项目开发很有帮助。', 14, 1, '2024-02-10 15:30:00'),
('course', 2, 19, '数据库课程难度适中，李老师讲解详细，实验指导书很全面。', 16, 2, '2024-02-11 10:20:00'),
('course', 3, 20, 'Web开发课程很有趣，学到了很多前端技术！', 12, 1, '2024-02-12 11:45:00');

-- 项目评论
INSERT INTO comments (resource_type, resource_id, user_id, content, love_count, reply_total, created_at) VALUES
('project', 1, 18, '这个项目架构很清晰，代码质量很高，值得学习！', 20, 3, '2024-02-25 10:00:00'),
('project', 1, 19, '在线学习平台功能很完整，适合作为毕业设计参考。', 15, 2, '2024-02-25 14:20:00'),
('project', 3, 20, '博客系统的Markdown编辑器很好用，UI设计也很美观。', 13, 1, '2024-02-26 09:30:00');

-- 评论回复 (parent_id是评论ID，会在插入后自动生成，这里使用相对ID)
INSERT INTO comments (resource_type, resource_id, parent_id, user_id, content, love_count, created_at) VALUES
('tool', 1, 1, 18, '同意！特别是Live Server插件，开发效率提升很多。', 5, '2024-01-20 16:00:00'),
('tool', 1, 1, 19, '还有GitLens插件也很实用，可以直接看到每行代码的作者。', 4, '2024-01-20 17:30:00'),
('tool', 1, 2, 17, '是的，而且支持很多编程语言，配置也很简单。', 3, '2024-01-20 18:00:00'),
('course', 1, 5, 19, '同感！特别是需求分析那一章，讲得很透彻。', 6, '2024-02-10 20:00:00'),
('course', 1, 5, 20, '作业设计也很有挑战性，能学到很多东西。', 5, '2024-02-11 08:30:00');

-- 评论点赞 (comment_id会在插入后自动生成，这里使用相对ID，实际插入时需要调整)
INSERT INTO comment_likes (comment_id, user_id, created_at) VALUES
(1, 18), (1, 19), (1, 20), (1, 16),
(2, 17), (2, 19),
(3, 17), (3, 18), (3, 20),
(4, 17), (4, 18), (4, 19), (4, 20), (4, 16),
(5, 18), (5, 19), (5, 20),
(6, 17), (6, 19),
(7, 17), (7, 18), (7, 20),
(8, 18), (8, 19), (8, 16),
(9, 19), (9, 20), (9, 16),
(10, 18), (10, 20),
(11, 19);

-- ==================== 收藏数据 ====================

INSERT INTO collections (user_id, resource_type, resource_id, created_at) VALUES
-- 学生17的收藏 (student_li)
(17, 'tool', 1, '2024-01-25 10:00:00'),
(17, 'tool', 2, '2024-01-25 11:00:00'),
(17, 'tool', 3, '2024-01-26 09:00:00'),
(17, 'course', 1, '2024-02-12 10:00:00'),
(17, 'course', 3, '2024-02-12 14:00:00'),
(17, 'project', 1, '2024-02-26 10:00:00'),
(17, 'project', 3, '2024-02-26 15:00:00'),
-- 学生18的收藏 (student_wang)
(18, 'tool', 1, '2024-01-24 09:00:00'),
(18, 'tool', 4, '2024-01-27 10:00:00'),
(18, 'course', 2, '2024-02-13 10:00:00'),
(18, 'project', 2, '2024-02-27 10:00:00'),
-- 学生19的收藏 (student_zhao)
(19, 'tool', 2, '2024-01-25 15:00:00'),
(19, 'tool', 5, '2024-01-28 10:00:00'),
(19, 'course', 3, '2024-02-14 10:00:00'),
(19, 'project', 1, '2024-02-28 10:00:00'),
(19, 'project', 4, '2024-02-28 14:00:00');

-- ==================== 点赞数据 ====================

INSERT INTO likes (user_id, resource_type, resource_id, created_at) VALUES
-- 工具点赞 (user_id: 16-20为学生)
(17, 'tool', 1, '2024-01-25 10:00:00'),
(18, 'tool', 1, '2024-01-25 11:00:00'),
(19, 'tool', 1, '2024-01-25 12:00:00'),
(17, 'tool', 2, '2024-01-26 10:00:00'),
(18, 'tool', 2, '2024-01-26 11:00:00'),
(19, 'tool', 3, '2024-01-27 10:00:00'),
(20, 'tool', 3, '2024-01-27 11:00:00'),
(16, 'tool', 3, '2024-01-27 12:00:00'),
-- 课程点赞
(17, 'course', 1, '2024-02-12 10:00:00'),
(18, 'course', 1, '2024-02-12 11:00:00'),
(19, 'course', 1, '2024-02-12 12:00:00'),
(17, 'course', 2, '2024-02-13 10:00:00'),
(18, 'course', 2, '2024-02-13 11:00:00'),
-- 项目点赞
(17, 'project', 1, '2024-02-26 10:00:00'),
(18, 'project', 1, '2024-02-26 11:00:00'),
(19, 'project', 1, '2024-02-26 12:00:00'),
(17, 'project', 3, '2024-02-27 10:00:00'),
(18, 'project', 3, '2024-02-27 11:00:00');

-- ==================== 状态日志 ====================

INSERT INTO resource_status_logs (resource_type, resource_id, old_status, new_status, operator_id, operate_time) VALUES
('tool', 1, NULL, 'pending', 17, '2024-01-14 09:00:00'),
('tool', 1, 'pending', 'approved', 1, '2024-01-15 10:00:00'),
('tool', 2, NULL, 'pending', 18, '2024-01-15 10:00:00'),
('tool', 2, 'pending', 'approved', 1, '2024-01-16 10:00:00'),
('tool', 3, NULL, 'pending', 19, '2024-01-16 11:00:00'),
('tool', 3, 'pending', 'approved', 1, '2024-01-17 10:00:00'),
('tool', 4, NULL, 'pending', 20, '2024-01-17 14:00:00'),
('tool', 4, 'pending', 'approved', 1, '2024-01-18 10:00:00'),
('tool', 5, NULL, 'pending', 16, '2024-01-18 15:00:00'),
('tool', 5, 'pending', 'approved', 1, '2024-01-19 10:00:00'),
('tool', 9, NULL, 'pending', 17, NOW()),
('tool', 10, NULL, 'pending', 18, NOW()),
('tool', 11, NULL, 'pending', 19, NOW());

-- ==================== 使用说明 ====================
-- 所有用户密码都是: 123456
-- 可以使用以下账号登录：
-- 管理员：admin / 123456
-- 教师：teacher_zhang / 123456, teacher_li / 123456, teacher_wang / 123456
-- 学生：student_zhang / 123456, student_li / 123456, student_wang / 123456, student_zhao / 123456, student_sun / 123456

-- 注意：密码使用bcrypt加密，哈希值为：$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O
