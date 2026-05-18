-- 课程详情数据导入脚本
-- 将课程详情、资源和评论数据导入到现有数据库中

USE softeng;

-- ========================================================
-- 步骤1: 检查并添加courses表缺少的字段
-- ========================================================
-- 如果courses表没有description字段，则添加
SET @col_exists = (
    SELECT COUNT(*)
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = 'softeng'
      AND TABLE_NAME = 'courses'
      AND COLUMN_NAME = 'description'
);

SET @sql = IF(@col_exists = 0,
    'ALTER TABLE courses ADD COLUMN description TEXT COMMENT ''课程描述'' AFTER cover',
    'SELECT ''description column already exists'' AS message'
);

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- ========================================================
-- 步骤2: 更新现有课程数据（根据course_id匹配）
-- 注意：这里假设course_id从1开始，如果现有数据ID不同，需要调整
-- ========================================================

-- 更新课程基本信息（loves对应loves字段，resource_count用于后续更新resource_count统计）
-- 由于数据中有些course_id可能已经存在，我们需要先检查，避免冲突

-- 先更新已存在的课程（如果course_id在1001-1020范围内）
UPDATE courses SET
    name = CASE course_id
        WHEN 1001 THEN '高等数学(上)'
        WHEN 1002 THEN 'C语言程序设计'
        WHEN 1003 THEN '计算机导论'
        WHEN 1004 THEN '高等数学(下)'
        WHEN 1005 THEN '离散数学'
        WHEN 1006 THEN '面向对象程序设计(Java)'
        WHEN 1007 THEN '数据结构与算法'
        WHEN 1008 THEN '计算机组成原理'
        WHEN 1009 THEN 'Python脚本编程'
        WHEN 1010 THEN '操作系统'
        WHEN 1011 THEN '计算机网络'
        WHEN 1012 THEN '数据库系统原理'
        WHEN 1013 THEN 'Linux环境编程'
        WHEN 1014 THEN '软件工程导论'
        WHEN 1015 THEN 'Web前端开发'
        WHEN 1016 THEN '算法分析与设计'
        WHEN 1017 THEN '软件测试与质量保证'
        WHEN 1018 THEN '移动应用开发(Android)'
        WHEN 1019 THEN '机器学习基础'
        WHEN 1020 THEN '编译原理'
    END,
    semester = CASE course_id
        WHEN 1001 THEN '1-1'
        WHEN 1002 THEN '1-1'
        WHEN 1003 THEN '1-1'
        WHEN 1004 THEN '1-2'
        WHEN 1005 THEN '1-2'
        WHEN 1006 THEN '1-2'
        WHEN 1007 THEN '2-1'
        WHEN 1008 THEN '2-1'
        WHEN 1009 THEN '2-1'
        WHEN 1010 THEN '2-2'
        WHEN 1011 THEN '2-2'
        WHEN 1012 THEN '2-2'
        WHEN 1013 THEN '2-2'
        WHEN 1014 THEN '3-1'
        WHEN 1015 THEN '3-1'
        WHEN 1016 THEN '3-1'
        WHEN 1017 THEN '3-2'
        WHEN 1018 THEN '3-2'
        WHEN 1019 THEN '3-2'
        WHEN 1020 THEN '3-2'
    END,
    credit = CASE course_id
        WHEN 1001 THEN 5
        WHEN 1002 THEN 4
        WHEN 1003 THEN 2
        WHEN 1004 THEN 5
        WHEN 1005 THEN 4
        WHEN 1006 THEN 4
        WHEN 1007 THEN 4
        WHEN 1008 THEN 4
        WHEN 1009 THEN 2
        WHEN 1010 THEN 4
        WHEN 1011 THEN 4
        WHEN 1012 THEN 3
        WHEN 1013 THEN 3
        WHEN 1014 THEN 3
        WHEN 1015 THEN 3
        WHEN 1016 THEN 3
        WHEN 1017 THEN 2
        WHEN 1018 THEN 3
        WHEN 1019 THEN 2
        WHEN 1020 THEN 4
    END,
    cover = CASE course_id
        WHEN 1001 THEN 'https://images.unsplash.com/photo-1635070041078-e363dbe005cb?w=500&auto=format&fit=crop'
        WHEN 1002 THEN 'https://images.unsplash.com/photo-1515879218367-8466d910aaa4?w=500&auto=format&fit=crop'
        WHEN 1003 THEN 'https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=500&auto=format&fit=crop'
        WHEN 1004 THEN 'https://images.unsplash.com/photo-1596495578065-6e0763fa1178?w=500&auto=format&fit=crop'
        WHEN 1005 THEN 'https://images.unsplash.com/photo-1509228468518-180dd4864904?w=500&auto=format&fit=crop'
        WHEN 1006 THEN 'https://images.unsplash.com/photo-1526379095098-d400fd0bf935?w=500&auto=format&fit=crop'
        WHEN 1007 THEN 'https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=500&auto=format&fit=crop'
        WHEN 1008 THEN 'https://images.unsplash.com/photo-1591453089816-0fbb971b454c?w=500&auto=format&fit=crop'
        WHEN 1009 THEN 'https://images.unsplash.com/photo-1526379879527-8559ecfcaec0?w=500&auto=format&fit=crop'
        WHEN 1010 THEN 'https://images.unsplash.com/photo-1518432031352-d6fc5c10da5a?w=500&auto=format&fit=crop'
        WHEN 1011 THEN 'https://images.unsplash.com/photo-1544197150-b99a580bbcbf?w=500&auto=format&fit=crop'
        WHEN 1012 THEN 'https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=500&auto=format&fit=crop'
        WHEN 1013 THEN 'https://images.unsplash.com/photo-1629654297299-c8506221ca97?w=500&auto=format&fit=crop'
        WHEN 1014 THEN 'https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=500&auto=format&fit=crop'
        WHEN 1015 THEN 'https://images.unsplash.com/photo-1587620962725-abab7fe55159?w=500&auto=format&fit=crop'
        WHEN 1016 THEN 'https://images.unsplash.com/photo-1550751827-4bd374c3f58b?w=500&auto=format&fit=crop'
        WHEN 1017 THEN 'https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=500&auto=format&fit=crop'
        WHEN 1018 THEN 'https://images.unsplash.com/photo-1610433571932-d1964175b9f1?w=500&auto=format&fit=crop'
        WHEN 1019 THEN 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=500&auto=format&fit=crop'
        WHEN 1020 THEN 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=500&auto=format&fit=crop'
    END,
    description = CASE course_id
        WHEN 1001 THEN '理工科基础数学课程，重点讲解极限与连续、导数与微分、中值定理与导数的应用、不定积分、定积分及其应用。是后续学习专业课的基石。'
        WHEN 1002 THEN '编程入门第一课。从零开始讲解计算机编程，涵盖数据类型、控制结构、数组、指针、结构体等C语言核心语法，培养计算思维。'
        WHEN 1003 THEN '计算机科学的全景概览。介绍计算机发展史、基本硬件组成、操作系统原理概论、网络基础及计算机伦理，帮助新生建立专业认知。'
        WHEN 1004 THEN '微积分进阶，涵盖空间解析几何、多元函数微分法、重积分、曲线积分与曲面积分、无穷级数等内容。'
        WHEN 1005 THEN '计算机科学的数学基础。包含集合论、数理逻辑、图论、代数结构。这是数据结构和算法分析的理论基础，非常重要！'
        WHEN 1006 THEN '深入讲解Java语言与面向对象思想（封装、继承、多态），涵盖Java SE核心库、异常处理、IO流及多线程编程。'
        WHEN 1007 THEN '程序设计的灵魂，考研面试必考。内容包括线性表、栈与队列、树与二叉树、图、查找与排序。本课程难度较大，需要大量代码实践。'
        WHEN 1008 THEN '深入理解计算机硬件系统工作原理：数据的表示、运算方法、存储系统、指令系统、CPU设计、总线与I/O系统。'
        WHEN 1009 THEN '人生苦短，我用Python。快速掌握Python语法，学习爬虫基础、数据分析库(Pandas/Numpy)及自动化办公脚本编写。'
        WHEN 1010 THEN '管理计算机硬件与软件资源的系统软件。重点讲解进程管理、内存管理、文件系统、设备管理。理解并发、锁、死锁等核心概念。'
        WHEN 1011 THEN '自顶向下方法讲解网络协议栈：HTTP、TCP/IP、路由算法、局域网技术。理解互联网是如何连接世界的。'
        WHEN 1012 THEN '深入讲解关系数据库系统的基本概念、理论和设计方法，包括ER图设计、SQL语言高阶应用、事务处理及并发控制等核心内容。'
        WHEN 1013 THEN '熟悉Linux指令与Shell脚本，掌握Vim使用、系统调用、进程间通信。后端开发必备技能。'
        WHEN 1014 THEN '系统地介绍软件工程的概念、原理、方法和技术。涵盖需求分析、UML建模、软件设计模式、敏捷开发(Scrum)、DevOps概念入门。'
        WHEN 1015 THEN '现代前端技术栈 Vue3 + TS 实战开发。从HTML/CSS基础到现代前端框架Vue3的深度解析，包含组件化开发、状态管理Pinia、路由Vue Router等。'
        WHEN 1016 THEN '解决复杂问题的核心思维。涵盖分治法、动态规划、贪心算法、回溯法等经典算法策略，结合LeetCode真题进行实战讲解。'
        WHEN 1017 THEN '确保软件质量的关键环节。介绍黑盒测试、白盒测试、单元测试(JUnit)、自动化测试工具(Selenium)的使用。'
        WHEN 1018 THEN '开发你的第一个手机App。学习Kotlin语言基础，Activity生命周期，UI布局，网络请求Retrofit，本地存储Room。'
        WHEN 1019 THEN '人工智能入门。通俗易懂地讲解机器学习基本原理，介绍监督学习、非监督学习、线性回归、神经网络及Python scikit-learn实战。'
        WHEN 1020 THEN '计算机专业的"天书"。涵盖词法分析、语法分析、语义分析、代码生成与优化。理解编译器是如何翻译代码的。'
    END,
    loves = CASE course_id
        WHEN 1001 THEN 45
        WHEN 1002 THEN 120
        WHEN 1003 THEN 30
        WHEN 1004 THEN 42
        WHEN 1005 THEN 88
        WHEN 1006 THEN 210
        WHEN 1007 THEN 350
        WHEN 1008 THEN 150
        WHEN 1009 THEN 180
        WHEN 1010 THEN 220
        WHEN 1011 THEN 200
        WHEN 1012 THEN 160
        WHEN 1013 THEN 140
        WHEN 1014 THEN 130
        WHEN 1015 THEN 400
        WHEN 1016 THEN 280
        WHEN 1017 THEN 90
        WHEN 1018 THEN 150
        WHEN 1019 THEN 310
        WHEN 1020 THEN 60
    END
WHERE course_id BETWEEN 1001 AND 1020;

-- 如果课程不存在，则插入新课程
INSERT INTO courses (course_id, name, semester, credit, cover, description, loves, resource_type, created_at, updated_at)
SELECT * FROM (
    SELECT 1001 AS course_id, '高等数学(上)' AS name, '1-1' AS semester, 5 AS credit, 'https://images.unsplash.com/photo-1635070041078-e363dbe005cb?w=500&auto=format&fit=crop' AS cover, '理工科基础数学课程，重点讲解极限与连续、导数与微分、中值定理与导数的应用、不定积分、定积分及其应用。是后续学习专业课的基石。' AS description, 45 AS loves, 'course' AS resource_type, NOW() AS created_at, NOW() AS updated_at
    UNION ALL SELECT 1002, 'C语言程序设计', '1-1', 4, 'https://images.unsplash.com/photo-1515879218367-8466d910aaa4?w=500&auto=format&fit=crop', '编程入门第一课。从零开始讲解计算机编程，涵盖数据类型、控制结构、数组、指针、结构体等C语言核心语法，培养计算思维。', 120, 'course', NOW(), NOW()
    UNION ALL SELECT 1003, '计算机导论', '1-1', 2, 'https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=500&auto=format&fit=crop', '计算机科学的全景概览。介绍计算机发展史、基本硬件组成、操作系统原理概论、网络基础及计算机伦理，帮助新生建立专业认知。', 30, 'course', NOW(), NOW()
    UNION ALL SELECT 1004, '高等数学(下)', '1-2', 5, 'https://images.unsplash.com/photo-1596495578065-6e0763fa1178?w=500&auto=format&fit=crop', '微积分进阶，涵盖空间解析几何、多元函数微分法、重积分、曲线积分与曲面积分、无穷级数等内容。', 42, 'course', NOW(), NOW()
    UNION ALL SELECT 1005, '离散数学', '1-2', 4, 'https://images.unsplash.com/photo-1509228468518-180dd4864904?w=500&auto=format&fit=crop', '计算机科学的数学基础。包含集合论、数理逻辑、图论、代数结构。这是数据结构和算法分析的理论基础，非常重要！', 88, 'course', NOW(), NOW()
    UNION ALL SELECT 1006, '面向对象程序设计(Java)', '1-2', 4, 'https://images.unsplash.com/photo-1526379095098-d400fd0bf935?w=500&auto=format&fit=crop', '深入讲解Java语言与面向对象思想（封装、继承、多态），涵盖Java SE核心库、异常处理、IO流及多线程编程。', 210, 'course', NOW(), NOW()
    UNION ALL SELECT 1007, '数据结构与算法', '2-1', 4, 'https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=500&auto=format&fit=crop', '程序设计的灵魂，考研面试必考。内容包括线性表、栈与队列、树与二叉树、图、查找与排序。本课程难度较大，需要大量代码实践。', 350, 'course', NOW(), NOW()
    UNION ALL SELECT 1008, '计算机组成原理', '2-1', 4, 'https://images.unsplash.com/photo-1591453089816-0fbb971b454c?w=500&auto=format&fit=crop', '深入理解计算机硬件系统工作原理：数据的表示、运算方法、存储系统、指令系统、CPU设计、总线与I/O系统。', 150, 'course', NOW(), NOW()
    UNION ALL SELECT 1009, 'Python脚本编程', '2-1', 2, 'https://images.unsplash.com/photo-1526379879527-8559ecfcaec0?w=500&auto=format&fit=crop', '人生苦短，我用Python。快速掌握Python语法，学习爬虫基础、数据分析库(Pandas/Numpy)及自动化办公脚本编写。', 180, 'course', NOW(), NOW()
    UNION ALL SELECT 1010, '操作系统', '2-2', 4, 'https://images.unsplash.com/photo-1518432031352-d6fc5c10da5a?w=500&auto=format&fit=crop', '管理计算机硬件与软件资源的系统软件。重点讲解进程管理、内存管理、文件系统、设备管理。理解并发、锁、死锁等核心概念。', 220, 'course', NOW(), NOW()
    UNION ALL SELECT 1011, '计算机网络', '2-2', 4, 'https://images.unsplash.com/photo-1544197150-b99a580bbcbf?w=500&auto=format&fit=crop', '自顶向下方法讲解网络协议栈：HTTP、TCP/IP、路由算法、局域网技术。理解互联网是如何连接世界的。', 200, 'course', NOW(), NOW()
    UNION ALL SELECT 1012, '数据库系统原理', '2-2', 3, 'https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=500&auto=format&fit=crop', '深入讲解关系数据库系统的基本概念、理论和设计方法，包括ER图设计、SQL语言高阶应用、事务处理及并发控制等核心内容。', 160, 'course', NOW(), NOW()
    UNION ALL SELECT 1013, 'Linux环境编程', '2-2', 3, 'https://images.unsplash.com/photo-1629654297299-c8506221ca97?w=500&auto=format&fit=crop', '熟悉Linux指令与Shell脚本，掌握Vim使用、系统调用、进程间通信。后端开发必备技能。', 140, 'course', NOW(), NOW()
    UNION ALL SELECT 1014, '软件工程导论', '3-1', 3, 'https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=500&auto=format&fit=crop', '系统地介绍软件工程的概念、原理、方法和技术。涵盖需求分析、UML建模、软件设计模式、敏捷开发(Scrum)、DevOps概念入门。', 130, 'course', NOW(), NOW()
    UNION ALL SELECT 1015, 'Web前端开发', '3-1', 3, 'https://images.unsplash.com/photo-1587620962725-abab7fe55159?w=500&auto=format&fit=crop', '现代前端技术栈 Vue3 + TS 实战开发。从HTML/CSS基础到现代前端框架Vue3的深度解析，包含组件化开发、状态管理Pinia、路由Vue Router等。', 400, 'course', NOW(), NOW()
    UNION ALL SELECT 1016, '算法分析与设计', '3-1', 3, 'https://images.unsplash.com/photo-1550751827-4bd374c3f58b?w=500&auto=format&fit=crop', '解决复杂问题的核心思维。涵盖分治法、动态规划、贪心算法、回溯法等经典算法策略，结合LeetCode真题进行实战讲解。', 280, 'course', NOW(), NOW()
    UNION ALL SELECT 1017, '软件测试与质量保证', '3-2', 2, 'https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=500&auto=format&fit=crop', '确保软件质量的关键环节。介绍黑盒测试、白盒测试、单元测试(JUnit)、自动化测试工具(Selenium)的使用。', 90, 'course', NOW(), NOW()
    UNION ALL SELECT 1018, '移动应用开发(Android)', '3-2', 3, 'https://images.unsplash.com/photo-1610433571932-d1964175b9f1?w=500&auto=format&fit=crop', '开发你的第一个手机App。学习Kotlin语言基础，Activity生命周期，UI布局，网络请求Retrofit，本地存储Room。', 150, 'course', NOW(), NOW()
    UNION ALL SELECT 1019, '机器学习基础', '3-2', 2, 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=500&auto=format&fit=crop', '人工智能入门。通俗易懂地讲解机器学习基本原理，介绍监督学习、非监督学习、线性回归、神经网络及Python scikit-learn实战。', 310, 'course', NOW(), NOW()
    UNION ALL SELECT 1020, '编译原理', '3-2', 4, 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=500&auto=format&fit=crop', '计算机专业的"天书"。涵盖词法分析、语法分析、语义分析、代码生成与优化。理解编译器是如何翻译代码的。', 60, 'course', NOW(), NOW()
) AS new_courses
WHERE NOT EXISTS (SELECT 1 FROM courses WHERE courses.course_id = new_courses.course_id);

-- ========================================================
-- 步骤3: 插入教师信息到course_teachers表
-- ========================================================
INSERT IGNORE INTO course_teachers (course_id, teacher_name)
VALUES
(1001, '陈高数'), (1002, '刘伟'), (1003, '王芳'),
(1004, '陈高数'), (1005, '张逻辑'), (1006, '赵强'),
(1007, '严蔚敏'), (1008, '李硬件'), (1009, 'Alice'),
(1010, 'Andrew'), (1011, '谢希仁'), (1012, '王DB'),
(1013, 'Linus'), (1014, '张架构'), (1015, '尤雨溪'),
(1016, 'AlgorithmGod'), (1017, '李测试'), (1018, 'Google'),
(1019, 'AI Master'), (1020, '陈龙书');

-- ========================================================
-- 步骤4: 插入课程分类到course_categories表
-- type字段：公必、专必、专选、公选
-- ========================================================
INSERT IGNORE INTO course_categories (course_id, category)
VALUES
(1001, '公必'), (1002, '专必'), (1003, '专必'),
(1004, '公必'), (1005, '专必'), (1006, '专必'),
(1007, '专必'), (1008, '专必'), (1009, '专选'),
(1010, '专必'), (1011, '专必'), (1012, '专必'),
(1013, '专选'), (1014, '专必'), (1015, '专选'),
(1016, '专必'), (1017, '专必'), (1018, '专选'),
(1019, '公选'), (1020, '专必');

-- ========================================================
-- 步骤5: 插入课程资源到course_resources_web表
-- 根据type判断：video/doc/tool -> web资源，上传的pdf/doc等 -> upload资源
-- ========================================================
-- C语言 (1002) 资源
INSERT IGNORE INTO course_resources_web (course_id, resource_intro, resource_url, sort_order, created_at)
VALUES
(1002, 'C语言常用函数速查手册', 'https://example.com/c_func.pdf', 0, NOW()),
(1002, 'Dev-C++ 5.11 安装包', 'https://sourceforge.net/', 1, NOW());

-- 离散数学 (1005) 资源
INSERT IGNORE INTO course_resources_web (course_id, resource_intro, resource_url, sort_order, created_at)
VALUES
(1005, '离散数学全套教学视频(30讲)', 'https://bilibili.com/video/xxx', 0, NOW()),
(1005, '图论习题集及答案', 'https://example.com/graph_theory.pdf', 1, NOW());

-- Java (1006) 资源
INSERT IGNORE INTO course_resources_web (course_id, resource_intro, resource_url, sort_order, created_at)
VALUES
(1006, 'IntelliJ IDEA 教育版激活教程', 'https://jetbrains.com', 0, NOW()),
(1006, 'Java核心卷1读书笔记', 'https://github.com/', 1, NOW());

-- 数据结构 (1007) 资源
INSERT IGNORE INTO course_resources_web (course_id, resource_intro, resource_url, sort_order, created_at)
VALUES
(1007, '严蔚敏数据结构课件PPT', 'https://pan.baidu.com/', 0, NOW()),
(1007, '链表反转动画演示', 'https://visualgo.net/', 1, NOW());

-- 数据结构上传资源（PDF文档）
INSERT IGNORE INTO course_resources_upload (course_id, resource_intro, resource_upload, sort_order, created_at)
VALUES
(1007, '期末重点考点押题', 'https://example.com/ds_final.docx', 2, NOW());

-- 操作系统 (1010) 资源
INSERT IGNORE INTO course_resources_web (course_id, resource_intro, resource_url, sort_order, created_at)
VALUES
(1010, 'PV操作经典例题讲解', 'https://example.com/pv.pdf', 0, NOW());

-- 计算机网络 (1011) 资源
INSERT IGNORE INTO course_resources_web (course_id, resource_intro, resource_url, sort_order, created_at)
VALUES
(1011, 'Wireshark 抓包工具', 'https://www.wireshark.org/', 0, NOW());

-- Web前端 (1015) 资源
INSERT IGNORE INTO course_resources_web (course_id, resource_intro, resource_url, sort_order, created_at)
VALUES
(1015, 'Vue3 组合式API实战教程', 'https://bilibili.com/vue3', 0, NOW()),
(1015, '前端面试题汇总2024版', 'https://github.com/interview', 1, NOW());

-- 算法 (1016) 资源
INSERT IGNORE INTO course_resources_web (course_id, resource_intro, resource_url, sort_order, created_at)
VALUES
(1016, 'LeetCode Hot 100 题解', 'https://leetcode.cn', 0, NOW());

-- ========================================================
-- 步骤6: 插入课程评论到comments表
-- 注意：评论需要关联user_id，这里我们先获取现有用户ID列表
-- 然后为每条评论随机分配一个用户ID
-- ========================================================

-- 首先，获取所有用户ID（假设至少有一些用户存在）
-- 如果用户表为空，我们需要先创建一些测试用户

-- 为评论分配用户：从现有用户中循环分配
-- 假设用户ID从1开始，我们使用模运算来分配

-- 插入评论（resource_type='course', resource_id=course_id）
-- 每条评论的user_id从1开始循环分配，love_count对应likes字段
-- 注意：这里假设至少有一些用户存在，如果用户表为空，需要先创建用户

INSERT INTO comments (resource_type, resource_id, user_id, content, love_count, created_at, parent_id)
SELECT 'course', course_id, 
    -- 使用课程ID和评论序号来分配用户ID（假设有用户ID 1-10，循环使用）
    ((ROW_NUMBER() OVER (ORDER BY course_id) - 1) % 10 + 1) AS user_id,
    content,
    love_count,
    comment_time,
    NULL AS parent_id
FROM (
    -- 1001 高数
    SELECT 1001 AS course_id, '这也太难了吧，完全听不懂极限定义...' AS content, 50 AS love_count, '2023-09-10 10:00:00' AS comment_time
    UNION ALL SELECT 1001, '陈老师讲得很细，课后习题一定要自己做。', 12, '2023-09-12 14:00:00'
    -- 1002 C语言
    UNION ALL SELECT 1002, '指针那里真的晕了，求大神指点！', 8, '2023-10-01 09:00:00'
    UNION ALL SELECT 1002, '多写代码，多调试，C语言其实不难。', 25, '2023-10-02 11:00:00'
    -- 1003 导论
    UNION ALL SELECT 1003, '这门课比较轻松，主要了解概念。', 5, '2023-09-15 16:00:00'
    UNION ALL SELECT 1003, '王老师人很好，期末给分高。', 30, '2023-12-20 10:00:00'
    -- 1004 高数下
    UNION ALL SELECT 1004, '重积分算到手断...', 44, '2024-03-01 10:00:00'
    UNION ALL SELECT 1004, '空间几何要注意画图辅助理解。', 10, '2024-03-05 12:00:00'
    -- 1005 离散
    UNION ALL SELECT 1005, '真值表画错了，痛失10分。', 2, '2024-04-01 10:00:00'
    UNION ALL SELECT 1005, '图论部分挺有意思的。', 6, '2024-04-10 14:00:00'
    -- 1006 Java
    UNION ALL SELECT 1006, '面向对象思想真的很重要，比面向过程好维护多了。', 100, '2024-05-01 09:00:00'
    UNION ALL SELECT 1006, '推荐直接用 IDEA，别用 Eclipse 了。', 200, '2024-05-02 10:00:00'
    -- 1007 数据结构
    UNION ALL SELECT 1007, '数据结构是考研专业课最难的一门，必须死磕。', 300, '2023-10-10 20:00:00'
    UNION ALL SELECT 1007, 'KMP算法看了一周才看懂...', 50, '2023-10-15 21:00:00'
    -- 1008 计组
    UNION ALL SELECT 1008, '补码反码原码搞得头晕。', 15, '2023-11-01 10:00:00'
    UNION ALL SELECT 1008, '流水线技术是提升CPU性能的关键。', 20, '2023-11-05 10:00:00'
    -- 1009 Python
    UNION ALL SELECT 1009, 'Python 真的简洁，写起来很爽。', 40, '2023-09-20 14:00:00'
    UNION ALL SELECT 1009, '爬虫作业要注意设置 User-Agent，不然会被封IP。', 60, '2023-09-25 15:00:00'
    -- 1010 OS
    UNION ALL SELECT 1010, '哲学家就餐问题很有趣。', 33, '2024-03-10 10:00:00'
    UNION ALL SELECT 1010, '一定要理解虚拟内存的概念。', 22, '2024-03-12 11:00:00'
    -- 1011 计网
    UNION ALL SELECT 1011, '三次握手和四次挥手背下来，面试必问。', 120, '2024-04-20 09:00:00'
    UNION ALL SELECT 1011, '子网掩码计算有点绕。', 10, '2024-04-22 10:00:00'
    -- 1012 数据库
    UNION ALL SELECT 1012, '范式理论虽然枯燥，但是设计数据库必须遵守。', 45, '2024-05-01 14:00:00'
    UNION ALL SELECT 1012, '左连接右连接傻傻分不清楚。', 20, '2024-05-05 16:00:00'
    -- 1013 Linux
    UNION ALL SELECT 1013, '如何在 Vim 中退出？在线等，挺急的。', 999, '2024-03-15 12:00:00'
    UNION ALL SELECT 1013, 'rm -rf /* 慎用！！', 500, '2024-03-16 13:00:00'
    -- 1014 软工导论
    UNION ALL SELECT 1014, '文档写得头大，还是写代码爽。', 15, '2023-10-10 09:00:00'
    UNION ALL SELECT 1014, '敏捷开发确实比瀑布模型灵活。', 20, '2023-10-12 10:00:00'
    -- 1015 前端
    UNION ALL SELECT 1015, 'Vue3 真的好用，Composition API 很香。', 88, '2023-11-20 10:00:00'
    UNION ALL SELECT 1015, '垂直居中为什么这么难调？？', 200, '2023-11-21 11:00:00'
    -- 1016 算法
    UNION ALL SELECT 1016, '动态规划的核心是状态转移方程。', 150, '2023-12-01 20:00:00'
    UNION ALL SELECT 1016, '刷了200题，感觉稍微入门了。', 60, '2023-12-05 21:00:00'
    -- 1017 软件测试
    UNION ALL SELECT 1017, '开发不仅要会写代码，还要会写测试用例。', 30, '2024-03-20 10:00:00'
    UNION ALL SELECT 1017, '自动化测试是趋势。', 25, '2024-03-22 11:00:00'
    -- 1018 Android
    UNION ALL SELECT 1018, 'AS 模拟器太吃内存了，建议用真机调试。', 40, '2024-04-10 15:00:00'
    UNION ALL SELECT 1018, 'Kotlin 语法糖很甜。', 35, '2024-04-12 16:00:00'
    -- 1019 AI
    UNION ALL SELECT 1019, '数学不好学 AI 有点吃力啊。', 55, '2024-05-15 10:00:00'
    UNION ALL SELECT 1019, 'PyTorch 比 TensorFlow 容易上手。', 45, '2024-05-18 11:00:00'
    -- 1020 编译原理
    UNION ALL SELECT 1020, '这就是传说中的劝退课吗？太难了！', 200, '2024-06-01 10:00:00'
    UNION ALL SELECT 1020, '写出一个编译器后的成就感无与伦比。', 10, '2024-06-05 12:00:00'
) AS comment_data
WHERE EXISTS (SELECT 1 FROM users LIMIT 1); -- 确保至少有一个用户存在

-- ========================================================
-- 步骤7: 更新课程点赞数（确保loves字段正确）
-- ========================================================
-- 这里已经在上面的UPDATE语句中更新了，不需要再次更新

-- ========================================================
-- 完成提示
-- ========================================================
SELECT 'Course detail data imported successfully!' AS message;

