-- 演示用种子数据：清空业务表后按固定主键插入，保证外键与列表/详情/评论/收藏/点赞一致。
-- 所有账号密码均为：123456（bcrypt 见下方常量）
-- 使用步骤：
--   1) 新建库：mysql -u root -p < database/schema.sql
--   2) 若是旧库升级：mysql -u root -p softeng < database/patch_existing_to_demo.sql
--   3) 灌演示数据：mysql -u root -p softeng < database/seed_demo.sql

USE softeng;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE comment_likes;
TRUNCATE TABLE comments;
TRUNCATE TABLE collections;
TRUNCATE TABLE likes;
TRUNCATE TABLE resource_status_logs;
TRUNCATE TABLE tool_tags;
TRUNCATE TABLE tool_images;
TRUNCATE TABLE tool_contributors;
TRUNCATE TABLE tools;
TRUNCATE TABLE course_resources_web;
TRUNCATE TABLE course_resources_upload;
TRUNCATE TABLE course_contributors;
TRUNCATE TABLE course_teachers;
TRUNCATE TABLE course_categories;
TRUNCATE TABLE courses;
TRUNCATE TABLE project_tech_stack;
TRUNCATE TABLE project_images;
TRUNCATE TABLE project_authors;
TRUNCATE TABLE projects;
TRUNCATE TABLE users;

SET FOREIGN_KEY_CHECKS = 1;

-- bcrypt('123456')
SET @pwd = '$2a$10$AMUByobB32fA2BFk4zUIIeaGifTZhXLIyMCR1uMR40zVQHvbCti1O';

INSERT INTO users (id, username, nickname, email, password, avatar, description, role, created_at) VALUES
(1, 'admin', '系统管理员', 'admin@softeng.edu.cn', @pwd, 'https://picsum.photos/seed/uadmin/200/200', '负责内容审核与平台运维', 'admin', '2024-01-01 08:00:00'),
(2, 'teacher_zhang', '张教授', 'zhang@softeng.edu.cn', @pwd, 'https://picsum.photos/seed/uzhang/200/200', '软件工程、需求工程方向', 'user', '2024-01-02 09:00:00'),
(3, 'teacher_li', '李老师', 'li@softeng.edu.cn', @pwd, 'https://picsum.photos/seed/uli/200/200', '数据库与分布式系统', 'user', '2024-01-02 10:00:00'),
(4, 'teacher_wang', '王老师', 'wang@softeng.edu.cn', @pwd, 'https://picsum.photos/seed/uwang/200/200', 'Web 全栈与工程实践', 'user', '2024-01-02 11:00:00'),
(5, 'student_zhang', '张晨', 'zhangchen@student.edu.cn', @pwd, 'https://picsum.photos/seed/s1/200/200', '大三软工，偏前端', 'user', '2024-01-10 08:00:00'),
(6, 'student_li', '李悦', 'liyue@student.edu.cn', @pwd, 'https://picsum.photos/seed/s2/200/200', '大三软工，偏后端', 'user', '2024-01-10 09:00:00'),
(7, 'student_wang', '王可', 'wangke@student.edu.cn', @pwd, 'https://picsum.photos/seed/s3/200/200', '大三软工，全栈', 'user', '2024-01-10 10:00:00'),
(8, 'student_zhao', '赵琪', 'zhaoqi@student.edu.cn', @pwd, 'https://picsum.photos/seed/s4/200/200', '对 DevOps 感兴趣', 'user', '2024-01-10 11:00:00'),
(9, 'student_sun', '孙浩', 'sunhao@student.edu.cn', @pwd, 'https://picsum.photos/seed/s5/200/200', '算法与数据结构爱好者', 'user', '2024-01-10 12:00:00'),
(10, 'student_zhou', '周宁', 'zhouning@student.edu.cn', @pwd, 'https://picsum.photos/seed/s6/200/200', '移动开发与 HCI', 'user', '2024-01-10 13:00:00'),
(11, 'student_wu', '吴岚', 'wulan@student.edu.cn', @pwd, 'https://picsum.photos/seed/s7/200/200', '测试与质量保障', 'user', '2024-01-10 14:00:00'),
(12, 'student_zheng', '郑远', 'zhengyuan@student.edu.cn', @pwd, 'https://picsum.photos/seed/s8/200/200', '开源与社区贡献', 'user', '2024-01-10 15:00:00');

INSERT INTO tools (resource_id, resource_type, resource_name, resource_link, description, description_detail, category, views, collections, loves, status, audit_time, submitter_id, created_at) VALUES
(1, 'tool', 'Visual Studio Code', 'https://code.visualstudio.com/', '轻量可扩展的代码编辑器', '支持多语言、调试、Git 集成与海量插件，适合课程实验与团队开发。', '开发工具', 1850, 120, 210, 'approved', '2024-01-15 10:00:00', 5, '2024-01-14 09:00:00'),
(2, 'tool', 'Postman', 'https://www.postman.com/', 'API 设计与联调工具', '集合、环境变量、Mock 与自动化测试一体化，接口课必备。', '测试工具', 1420, 95, 168, 'approved', '2024-01-16 10:00:00', 6, '2024-01-15 10:00:00'),
(3, 'tool', 'Git', 'https://git-scm.com/', '分布式版本控制', '分支模型、Code Review 与 CI 的基础，贯穿软件工程全流程。', '版本控制', 2400, 180, 260, 'approved', '2024-01-17 10:00:00', 7, '2024-01-16 11:00:00'),
(4, 'tool', 'Docker', 'https://www.docker.com/', '容器化交付', '镜像、Compose 与 DevContainers，统一开发/测试/生产环境。', '部署工具', 980, 72, 134, 'approved', '2024-01-18 10:00:00', 8, '2024-01-17 14:00:00'),
(5, 'tool', 'Notion', 'https://www.notion.so/', '知识库与项目管理', '文档、看板、数据库视图，适合小组作业与需求跟踪。', '协作工具', 760, 55, 92, 'approved', '2024-01-19 10:00:00', 9, '2024-01-18 15:00:00'),
(6, 'tool', 'Obsidian', 'https://obsidian.md/', '本地优先笔记', '双向链接与图谱，适合整理课程笔记与架构决策记录。', '效率工具', 620, 48, 81, 'approved', '2024-01-20 10:00:00', 10, '2024-01-19 16:00:00'),
(7, 'tool', 'draw.io', 'https://www.drawio.com/', '流程图与 UML', '快速绘制用例图、时序图与部署图，文档课常用。', '设计工具', 890, 63, 105, 'approved', '2024-01-21 10:00:00', 11, '2024-01-20 10:00:00'),
(8, 'tool', 'Jupyter', 'https://jupyter.org/', '交互式笔记本', '数据探索、可视化与教学演示，数据分析与机器学习入门。', '数据科学', 540, 40, 67, 'approved', '2024-01-22 10:00:00', 12, '2024-01-21 11:00:00'),
(9, 'tool', 'Redis Insight', 'https://redis.io/insight/', 'Redis 可视化管理', '浏览键、慢查询与内存分析，缓存与消息队列实验好帮手。', '数据库工具', 430, 28, 49, 'approved', '2024-01-23 10:00:00', 5, '2024-01-22 12:00:00'),
(10, 'tool', 'SonarQube', 'https://www.sonarsource.com/products/sonarqube/', '代码质量与漏洞扫描', '静态分析、技术债看板，适合质量保障专题演示。', '质量工具', 0, 0, 0, 'pending', NULL, 6, NOW()),
(11, 'tool', 'Linear', 'https://linear.app/', '敏捷需求与迭代', 'Roadmap、Issue 与 Git 集成，轻量替代重型项目管理。', '项目管理', 0, 0, 0, 'pending', NULL, 7, NOW()),
(12, 'tool', 'Penpot', 'https://penpot.app/', '开源 UI 设计', 'Figma 风格协作，适合开源课程与界面原型。', '设计工具', 0, 0, 0, 'pending', NULL, 8, NOW());

INSERT INTO tool_tags (tool_id, tag) VALUES
(1, '编辑器'), (1, '前端'), (1, '插件'),
(2, 'API'), (2, '测试'), (2, 'HTTP'),
(3, 'Git'), (3, '版本控制'), (3, '协作'),
(4, '容器'), (4, 'DevOps'), (4, '部署'),
(5, '文档'), (5, '看板'), (5, '协作'),
(6, '笔记'), (6, '知识管理'),
(7, 'UML'), (7, '建模'),
(8, 'Python'), (8, '数据科学'),
(9, 'Redis'), (9, '缓存'),
(10, '静态分析'), (11, '敏捷'), (12, '开源');

INSERT INTO tool_images (tool_id, image_url, sort_order) VALUES
(1, 'https://picsum.photos/seed/t1/800/450', 0),
(2, 'https://picsum.photos/seed/t2/800/450', 0),
(3, 'https://picsum.photos/seed/t3/800/450', 0);

INSERT INTO tool_contributors (tool_id, user_id) VALUES
(1, 5), (1, 6), (2, 6), (3, 7), (4, 8), (5, 9), (6, 10), (7, 11), (8, 12), (9, 5);

INSERT INTO courses (course_id, resource_type, name, semester, credit, description, cover, views, loves, collections, created_at) VALUES
(1, 'course', '软件工程导论', '2024春季', 3, '覆盖生命周期、敏捷与质量基础，配合案例研讨。', 'https://picsum.photos/seed/c1/800/500', 820, 68, 52, '2024-02-01 08:00:00'),
(2, 'course', '数据库系统原理', '2024春季', 3, '关系模型、SQL、事务与索引，含实验与课程设计衔接。', 'https://picsum.photos/seed/c2/800/500', 760, 62, 48, '2024-02-02 09:00:00'),
(3, 'course', 'Web 开发技术', '2024春季', 2, 'HTML/CSS/JS 与现代前端框架入门，完成小型全栈项目。', 'https://picsum.photos/seed/c3/800/500', 920, 74, 61, '2024-02-03 10:00:00'),
(4, 'course', '软件项目管理', '2024春季', 2, '进度、成本、风险与干系人管理，结合看板实践。', 'https://picsum.photos/seed/c4/800/500', 540, 44, 36, '2024-02-04 11:00:00'),
(5, 'course', '软件测试与质量保证', '2024春季', 2, '黑盒/白盒、自动化与持续测试流水线概念。', 'https://picsum.photos/seed/c5/800/500', 510, 40, 33, '2024-02-05 12:00:00'),
(6, 'course', '数据结构与算法', '2023秋季', 4, '复杂度分析、常用结构与算法策略，为面试与竞赛打基础。', 'https://picsum.photos/seed/c6/800/500', 1200, 96, 78, '2023-09-01 08:00:00'),
(7, 'course', '操作系统', '2023秋季', 3, '进程、内存、文件系统与并发，配套实验镜像。', 'https://picsum.photos/seed/c7/800/500', 880, 70, 55, '2023-09-02 09:00:00'),
(8, 'course', '计算机网络', '2023秋季', 3, '分层模型、TCP/IP 与常见应用协议实验。', 'https://picsum.photos/seed/c8/800/500', 790, 58, 47, '2023-09-03 10:00:00'),
(9, 'course', '移动应用开发', '2024秋季', 2, '跨平台与原生路线对比，完成可演示的移动原型。', 'https://picsum.photos/seed/c9/800/500', 120, 10, 8, NOW()),
(10, 'course', '人工智能基础', '2024秋季', 3, '机器学习流程、常见模型与伦理讨论入门。', 'https://picsum.photos/seed/c10/800/500', 95, 8, 6, NOW());

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

INSERT INTO course_contributors (course_id, user_id) VALUES
(1, 2), (2, 3), (3, 4), (4, 2), (5, 3), (6, 2), (7, 4), (8, 2), (9, 4), (10, 3);

INSERT INTO course_resources_web (course_id, resource_intro, resource_url, sort_order) VALUES
(1, '软件工程概述（公开课）', 'https://example.com/softeng/course/se101/intro', 1),
(1, '需求工程导学', 'https://example.com/softeng/course/se101/requirements', 2),
(2, 'MySQL 基础', 'https://example.com/softeng/course/db201/sql-basics', 1),
(3, '现代 JavaScript 入门', 'https://developer.mozilla.org/zh-CN/docs/Web/JavaScript', 1),
(3, 'Vue 官方文档', 'https://cn.vuejs.org/', 2);

INSERT INTO course_resources_upload (course_id, resource_intro, resource_upload, sort_order) VALUES
(1, '课程大纲 PDF', 'https://picsum.photos/seed/pdf1/100/100', 1),
(2, '实验指导书', 'https://picsum.photos/seed/pdf2/100/100', 1),
(3, '全栈实训资料包', 'https://picsum.photos/seed/zip1/100/100', 1);

INSERT INTO projects (project_id, resource_type, name, description, detail, github_url, category, cover, views, loves, collections, status, audit_time, created_at) VALUES
(1, 'project', '在线学习平台', 'Vue3 + Go 的全栈教学管理系统', '# 在线学习平台\n- 课程/作业/成绩\n- Docker 一键部署\n', 'https://github.com/example/learning-platform', '实训项目', 'https://picsum.photos/seed/p1/900/600', 520, 46, 58, 'approved', '2024-02-15 10:00:00', '2024-02-15 09:00:00'),
(2, 'project', '校园二手交易平台', 'React + Node 的校园电商原型', '# 校园二手交易\n- 商品与搜索\n- 站内信\n', 'https://github.com/example/campus-market', '课程设计', 'https://picsum.photos/seed/p2/900/600', 480, 38, 44, 'approved', '2024-02-16 10:00:00', '2024-02-16 09:00:00'),
(3, 'project', '个人博客系统', 'Markdown + 评论 + 标签', '# 博客系统\n- SEO 友好路由\n- 代码高亮\n', 'https://github.com/example/personal-blog', '个人项目', 'https://picsum.photos/seed/p3/900/600', 610, 52, 60, 'approved', '2024-02-17 10:00:00', '2024-02-17 09:00:00'),
(4, 'project', '任务协作看板', 'Django REST + Vue3 的看板应用', '# 看板\n- 列表/标签/成员\n', 'https://github.com/example/task-board', '实训项目', 'https://picsum.photos/seed/p4/900/600', 430, 36, 41, 'approved', '2024-02-18 10:00:00', '2024-02-18 09:00:00'),
(5, 'project', '图书管理系统', 'Java Swing + JDBC 桌面端', '# 图书管理\n- 借还书\n- 报表\n', 'https://github.com/example/library-swing', '课程设计', 'https://picsum.photos/seed/p5/900/600', 560, 48, 55, 'approved', '2024-02-19 10:00:00', '2024-02-19 09:00:00'),
(6, 'project', '实时聊天室', 'WebSocket + Redis 会话', '# 聊天室\n- 多房间\n- 在线状态\n', 'https://github.com/example/chat-room', '个人项目', 'https://picsum.photos/seed/p6/900/600', 360, 30, 35, 'approved', '2024-02-20 10:00:00', '2024-02-20 09:00:00'),
(7, 'project', '微服务订单演示', 'Spring Cloud 最小可运行集合', '# 订单演示\n- 网关/注册中心\n', 'https://github.com/example/order-demo', '实训项目', 'https://picsum.photos/seed/p7/900/600', 290, 24, 28, 'approved', '2024-02-21 10:00:00', '2024-02-21 09:00:00'),
(8, 'project', '智能排课原型', '启发式搜索 + 可视化课表', '# 排课\n- 约束建模\n', 'https://github.com/example/scheduling', '课程设计', 'https://picsum.photos/seed/p8/900/600', 240, 20, 22, 'approved', '2024-02-22 10:00:00', '2024-02-22 09:00:00'),
(9, 'project', '在线考试系统（待审）', '微服务架构在线考试', '待管理员审核后上架。', 'https://github.com/example/exam', '实训项目', 'https://picsum.photos/seed/p9/900/600', 0, 0, 0, 'pending', NULL, NOW()),
(10, 'project', '数据可视化大屏（待审）', 'ECharts + 后端聚合接口', '用于实验室数据展示。', 'https://github.com/example/bi-screen', '课程设计', 'https://picsum.photos/seed/p10/900/600', 0, 0, 0, 'pending', NULL, NOW());

INSERT INTO project_tech_stack (project_id, tech) VALUES
(1, 'Vue3'), (1, 'Go'), (1, 'Gin'), (1, 'MySQL'),
(2, 'React'), (2, 'Node.js'), (2, 'Express'),
(3, 'Vue3'), (3, 'Spring Boot'), (3, 'MySQL'),
(4, 'Vue3'), (4, 'Django'), (4, 'PostgreSQL'),
(5, 'Java'), (5, 'Swing'), (5, 'JDBC'),
(6, 'React'), (6, 'WebSocket'), (6, 'Redis'),
(7, 'Spring Cloud'), (7, 'Java'), (7, 'MySQL'),
(8, 'Python'), (8, 'Vue3'), (8, 'FastAPI'),
(9, 'Java'), (9, 'Spring Boot'),
(10, 'Vue3'), (10, 'ECharts'), (10, 'Go');

INSERT INTO project_images (project_id, image_url, sort_order) VALUES
(1, 'https://picsum.photos/seed/pi1/1200/700', 0),
(2, 'https://picsum.photos/seed/pi2/1200/700', 0),
(3, 'https://picsum.photos/seed/pi3/1200/700', 0);

INSERT INTO project_authors (project_id, user_id) VALUES
(1, 5), (1, 6), (2, 7), (2, 8), (3, 9), (4, 10), (5, 11), (6, 12), (7, 5), (8, 6), (9, 7), (10, 8);

INSERT INTO comments (comment_id, resource_type, resource_id, parent_id, user_id, content, love_count, reply_total, created_at) VALUES
(1, 'tool', 1, NULL, 6, 'VS Code 插件生态太强了，课堂演示很方便。', 14, 2, '2024-01-20 10:00:00'),
(2, 'tool', 1, NULL, 7, 'Remote SSH + Dev Container 组合写实训报告很省时间。', 11, 0, '2024-01-20 14:30:00'),
(3, 'tool', 2, NULL, 8, '接口课基本离不开 Postman，环境变量管理很清晰。', 9, 1, '2024-01-21 09:15:00'),
(4, 'tool', 3, NULL, 9, '分支策略讲清楚之后，小组协作顺畅很多。', 16, 1, '2024-01-22 11:20:00'),
(5, 'course', 1, NULL, 6, '需求章节案例贴近真实项目，作业量适中。', 20, 2, '2024-02-10 09:00:00'),
(6, 'course', 1, NULL, 7, '希望再多一节敏捷回顾实战。', 12, 0, '2024-02-10 15:30:00'),
(7, 'course', 3, NULL, 8, '全栈小项目做完对前后端协作理解深了很多。', 15, 1, '2024-02-12 11:45:00'),
(8, 'project', 1, NULL, 7, '代码结构清晰，Docker Compose 一键起来，适合当模板。', 22, 2, '2024-02-25 10:00:00'),
(9, 'project', 3, NULL, 10, 'Markdown 渲染与目录导航体验不错。', 13, 0, '2024-02-26 09:30:00');

INSERT INTO comments (comment_id, resource_type, resource_id, parent_id, user_id, content, love_count, reply_total, created_at) VALUES
(10, 'tool', 1, 1, 8, '同意，Live Server + ESLint 课堂演示很顺手。', 6, 0, '2024-01-20 16:00:00'),
(11, 'tool', 1, 1, 9, 'GitLens 看行级 blame 讲解代码责任特好用。', 5, 0, '2024-01-20 17:30:00'),
(12, 'tool', 3, 4, 10, 'rebase 交互式整理提交历史我们组都在用。', 4, 0, '2024-01-22 18:00:00'),
(13, 'course', 1, 5, 9, '补充阅读链接如果能放到资源区就更集中了。', 5, 0, '2024-02-11 08:30:00'),
(14, 'project', 1, 8, 6, 'CI 流水线那一段注释写得很细，答辩加分。', 7, 0, '2024-02-25 16:00:00');

INSERT INTO comment_likes (comment_id, user_id) VALUES
(1, 5), (1, 7), (1, 8), (1, 9),
(2, 5), (2, 6),
(3, 5), (3, 6), (3, 7),
(4, 5), (4, 6), (4, 7), (4, 8),
(5, 5), (5, 7), (5, 8), (5, 9),
(8, 5), (8, 6), (8, 9), (8, 10);

INSERT INTO collections (user_id, resource_type, resource_id, created_at) VALUES
(6, 'tool', 1, '2024-01-25 10:00:00'), (6, 'tool', 2, '2024-01-25 11:00:00'), (6, 'course', 1, '2024-02-12 10:00:00'), (6, 'project', 1, '2024-02-26 10:00:00'),
(7, 'tool', 3, '2024-01-26 10:00:00'), (7, 'course', 3, '2024-02-13 10:00:00'), (7, 'project', 3, '2024-02-27 10:00:00'),
(8, 'tool', 4, '2024-01-27 10:00:00'), (8, 'course', 2, '2024-02-14 10:00:00'), (8, 'project', 2, '2024-02-28 10:00:00'),
(9, 'tool', 5, '2024-01-28 10:00:00'), (9, 'project', 4, '2024-02-28 14:00:00'),
(10, 'course', 6, '2024-02-15 10:00:00'), (10, 'project', 5, '2024-02-29 10:00:00');

INSERT INTO likes (user_id, resource_type, resource_id, created_at) VALUES
(6, 'tool', 1, '2024-01-25 10:00:00'), (7, 'tool', 1, '2024-01-25 11:00:00'), (8, 'tool', 1, '2024-01-25 12:00:00'),
(6, 'tool', 2, '2024-01-26 10:00:00'), (7, 'tool', 3, '2024-01-27 10:00:00'), (8, 'tool', 3, '2024-01-27 11:00:00'),
(6, 'course', 1, '2024-02-12 10:00:00'), (7, 'course', 1, '2024-02-12 11:00:00'), (8, 'course', 3, '2024-02-13 10:00:00'),
(6, 'project', 1, '2024-02-26 10:00:00'), (7, 'project', 1, '2024-02-26 11:00:00'), (9, 'project', 3, '2024-02-27 10:00:00');

INSERT INTO resource_status_logs (resource_type, resource_id, old_status, new_status, operator_id, operate_time, reject_reason) VALUES
('tool', 1, NULL, 'pending', 5, '2024-01-14 09:00:00', NULL),
('tool', 1, 'pending', 'approved', 1, '2024-01-15 10:00:00', NULL),
('tool', 2, NULL, 'pending', 6, '2024-01-15 10:00:00', NULL),
('tool', 2, 'pending', 'approved', 1, '2024-01-16 10:00:00', NULL),
('project', 9, NULL, 'pending', 7, NOW(), NULL),
('project', 10, NULL, 'pending', 8, NOW(), NULL);

ALTER TABLE users AUTO_INCREMENT = 100;
ALTER TABLE tools AUTO_INCREMENT = 100;
ALTER TABLE courses AUTO_INCREMENT = 100;
ALTER TABLE projects AUTO_INCREMENT = 100;
ALTER TABLE comments AUTO_INCREMENT = 100;
