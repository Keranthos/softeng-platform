#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
课程详情数据导入脚本
将课程详情、资源和评论数据导入到现有数据库中
"""

import mysql.connector
from mysql.connector import Error
import random
from datetime import datetime

# 数据库配置
# 默认配置（从 config.go 中获取）
# 如需修改，可以直接编辑下面的值，或设置环境变量
import os

DB_CONFIG = {
    'host': os.getenv('DB_HOST', 'localhost'),
    'user': os.getenv('DB_USER', 'root'),
    'password': os.getenv('DB_PASSWORD', 'Wan05609'),  # 默认密码，可从环境变量覆盖
    'database': os.getenv('DB_NAME', 'softeng'),
    'charset': 'utf8mb4',
    'collation': 'utf8mb4_unicode_ci'
}

def connect_db():
    """连接数据库"""
    try:
        conn = mysql.connector.connect(**DB_CONFIG)
        return conn
    except Error as e:
        print(f"数据库连接错误: {e}")
        return None

def get_user_ids(conn):
    """获取所有用户ID列表"""
    cursor = conn.cursor()
    try:
        cursor.execute("SELECT id FROM users ORDER BY id")
        user_ids = [row[0] for row in cursor.fetchall()]
        return user_ids if user_ids else [1]  # 如果没有用户，至少返回[1]（假设存在）
    except Error as e:
        print(f"查询用户ID错误: {e}")
        return [1]
    finally:
        cursor.close()

def add_description_column(conn):
    """检查并添加description字段"""
    cursor = conn.cursor()
    try:
        # 检查字段是否存在
        cursor.execute("""
            SELECT COUNT(*)
            FROM information_schema.COLUMNS
            WHERE TABLE_SCHEMA = 'softeng'
              AND TABLE_NAME = 'courses'
              AND COLUMN_NAME = 'description'
        """)
        exists = cursor.fetchone()[0]
        
        if not exists:
            cursor.execute("ALTER TABLE courses ADD COLUMN description TEXT COMMENT '课程描述' AFTER cover")
            conn.commit()
            print("已添加description字段")
        else:
            print("description字段已存在")
    except Error as e:
        print(f"添加description字段错误: {e}")
        conn.rollback()
    finally:
        cursor.close()

def import_courses(conn):
    """导入课程数据"""
    cursor = conn.cursor()
    
    courses_data = [
        (1001, '高等数学(上)', '1-1', 5, 'https://images.unsplash.com/photo-1635070041078-e363dbe005cb?w=500&auto=format&fit=crop', 
         '理工科基础数学课程，重点讲解极限与连续、导数与微分、中值定理与导数的应用、不定积分、定积分及其应用。是后续学习专业课的基石。', 45, '陈高数', '公必'),
        (1002, 'C语言程序设计', '1-1', 4, 'https://images.unsplash.com/photo-1515879218367-8466d910aaa4?w=500&auto=format&fit=crop',
         '编程入门第一课。从零开始讲解计算机编程，涵盖数据类型、控制结构、数组、指针、结构体等C语言核心语法，培养计算思维。', 120, '刘伟', '专必'),
        (1003, '计算机导论', '1-1', 2, 'https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=500&auto=format&fit=crop',
         '计算机科学的全景概览。介绍计算机发展史、基本硬件组成、操作系统原理概论、网络基础及计算机伦理，帮助新生建立专业认知。', 30, '王芳', '专必'),
        (1004, '高等数学(下)', '1-2', 5, 'https://images.unsplash.com/photo-1596495578065-6e0763fa1178?w=500&auto=format&fit=crop',
         '微积分进阶，涵盖空间解析几何、多元函数微分法、重积分、曲线积分与曲面积分、无穷级数等内容。', 42, '陈高数', '公必'),
        (1005, '离散数学', '1-2', 4, 'https://images.unsplash.com/photo-1509228468518-180dd4864904?w=500&auto=format&fit=crop',
         '计算机科学的数学基础。包含集合论、数理逻辑、图论、代数结构。这是数据结构和算法分析的理论基础，非常重要！', 88, '张逻辑', '专必'),
        (1006, '面向对象程序设计(Java)', '1-2', 4, 'https://images.unsplash.com/photo-1526379095098-d400fd0bf935?w=500&auto=format&fit=crop',
         '深入讲解Java语言与面向对象思想（封装、继承、多态），涵盖Java SE核心库、异常处理、IO流及多线程编程。', 210, '赵强', '专必'),
        (1007, '数据结构与算法', '2-1', 4, 'https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=500&auto=format&fit=crop',
         '程序设计的灵魂，考研面试必考。内容包括线性表、栈与队列、树与二叉树、图、查找与排序。本课程难度较大，需要大量代码实践。', 350, '严蔚敏', '专必'),
        (1008, '计算机组成原理', '2-1', 4, 'https://images.unsplash.com/photo-1591453089816-0fbb971b454c?w=500&auto=format&fit=crop',
         '深入理解计算机硬件系统工作原理：数据的表示、运算方法、存储系统、指令系统、CPU设计、总线与I/O系统。', 150, '李硬件', '专必'),
        (1009, 'Python脚本编程', '2-1', 2, 'https://images.unsplash.com/photo-1526379879527-8559ecfcaec0?w=500&auto=format&fit=crop',
         '人生苦短，我用Python。快速掌握Python语法，学习爬虫基础、数据分析库(Pandas/Numpy)及自动化办公脚本编写。', 180, 'Alice', '专选'),
        (1010, '操作系统', '2-2', 4, 'https://images.unsplash.com/photo-1518432031352-d6fc5c10da5a?w=500&auto=format&fit=crop',
         '管理计算机硬件与软件资源的系统软件。重点讲解进程管理、内存管理、文件系统、设备管理。理解并发、锁、死锁等核心概念。', 220, 'Andrew', '专必'),
        (1011, '计算机网络', '2-2', 4, 'https://images.unsplash.com/photo-1544197150-b99a580bbcbf?w=500&auto=format&fit=crop',
         '自顶向下方法讲解网络协议栈：HTTP、TCP/IP、路由算法、局域网技术。理解互联网是如何连接世界的。', 200, '谢希仁', '专必'),
        (1012, '数据库系统原理', '2-2', 3, 'https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=500&auto=format&fit=crop',
         '深入讲解关系数据库系统的基本概念、理论和设计方法，包括ER图设计、SQL语言高阶应用、事务处理及并发控制等核心内容。', 160, '王DB', '专必'),
        (1013, 'Linux环境编程', '2-2', 3, 'https://images.unsplash.com/photo-1629654297299-c8506221ca97?w=500&auto=format&fit=crop',
         '熟悉Linux指令与Shell脚本，掌握Vim使用、系统调用、进程间通信。后端开发必备技能。', 140, 'Linus', '专选'),
        (1014, '软件工程导论', '3-1', 3, 'https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=500&auto=format&fit=crop',
         '系统地介绍软件工程的概念、原理、方法和技术。涵盖需求分析、UML建模、软件设计模式、敏捷开发(Scrum)、DevOps概念入门。', 130, '张架构', '专必'),
        (1015, 'Web前端开发', '3-1', 3, 'https://images.unsplash.com/photo-1587620962725-abab7fe55159?w=500&auto=format&fit=crop',
         '现代前端技术栈 Vue3 + TS 实战开发。从HTML/CSS基础到现代前端框架Vue3的深度解析，包含组件化开发、状态管理Pinia、路由Vue Router等。', 400, '尤雨溪', '专选'),
        (1016, '算法分析与设计', '3-1', 3, 'https://images.unsplash.com/photo-1550751827-4bd374c3f58b?w=500&auto=format&fit=crop',
         '解决复杂问题的核心思维。涵盖分治法、动态规划、贪心算法、回溯法等经典算法策略，结合LeetCode真题进行实战讲解。', 280, 'AlgorithmGod', '专必'),
        (1017, '软件测试与质量保证', '3-2', 2, 'https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=500&auto=format&fit=crop',
         '确保软件质量的关键环节。介绍黑盒测试、白盒测试、单元测试(JUnit)、自动化测试工具(Selenium)的使用。', 90, '李测试', '专必'),
        (1018, '移动应用开发(Android)', '3-2', 3, 'https://images.unsplash.com/photo-1610433571932-d1964175b9f1?w=500&auto=format&fit=crop',
         '开发你的第一个手机App。学习Kotlin语言基础，Activity生命周期，UI布局，网络请求Retrofit，本地存储Room。', 150, 'Google', '专选'),
        (1019, '机器学习基础', '3-2', 2, 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=500&auto=format&fit=crop',
         '人工智能入门。通俗易懂地讲解机器学习基本原理，介绍监督学习、非监督学习、线性回归、神经网络及Python scikit-learn实战。', 310, 'AI Master', '公选'),
        (1020, '编译原理', '3-2', 4, 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=500&auto=format&fit=crop',
         '计算机专业的"天书"。涵盖词法分析、语法分析、语义分析、代码生成与优化。理解编译器是如何翻译代码的。', 60, '陈龙书', '专必'),
    ]
    
    try:
        for course_id, name, semester, credit, cover, description, loves, teacher, category_type in courses_data:
            # 先尝试更新
            cursor.execute("""
                UPDATE courses 
                SET name = %s, semester = %s, credit = %s, cover = %s, description = %s, loves = %s
                WHERE course_id = %s
            """, (name, semester, credit, cover, description, loves, course_id))
            
            # 如果更新失败（课程不存在），则插入
            if cursor.rowcount == 0:
                cursor.execute("""
                    INSERT INTO courses (course_id, name, semester, credit, cover, description, loves, resource_type, created_at, updated_at)
                    VALUES (%s, %s, %s, %s, %s, %s, %s, 'course', NOW(), NOW())
                """, (course_id, name, semester, credit, cover, description, loves))
            
            # 插入教师信息
            cursor.execute("""
                INSERT IGNORE INTO course_teachers (course_id, teacher_name)
                VALUES (%s, %s)
            """, (course_id, teacher))
            
            # 插入分类信息
            cursor.execute("""
                INSERT IGNORE INTO course_categories (course_id, category)
                VALUES (%s, %s)
            """, (course_id, category_type))
        
        conn.commit()
        print(f"成功导入/更新 {len(courses_data)} 门课程")
    except Error as e:
        print(f"导入课程数据错误: {e}")
        conn.rollback()
    finally:
        cursor.close()

def import_resources(conn):
    """导入课程资源数据"""
    cursor = conn.cursor()
    
    resources_data = [
        # (course_id, resource_type, resource_intro, url, is_upload)
        # C语言 (1002)
        (1002, 'doc', 'C语言常用函数速查手册', 'https://example.com/c_func.pdf', False),
        (1002, 'tool', 'Dev-C++ 5.11 安装包', 'https://sourceforge.net/', False),
        # 离散数学 (1005)
        (1005, 'video', '离散数学全套教学视频(30讲)', 'https://bilibili.com/video/xxx', False),
        (1005, 'doc', '图论习题集及答案', 'https://example.com/graph_theory.pdf', False),
        # Java (1006)
        (1006, 'tool', 'IntelliJ IDEA 教育版激活教程', 'https://jetbrains.com', False),
        (1006, 'doc', 'Java核心卷1读书笔记', 'https://github.com/', False),
        # 数据结构 (1007)
        (1007, 'doc', '严蔚敏数据结构课件PPT', 'https://pan.baidu.com/', False),
        (1007, 'video', '链表反转动画演示', 'https://visualgo.net/', False),
        (1007, 'doc', '期末重点考点押题', 'https://example.com/ds_final.docx', True),
        # 操作系统 (1010)
        (1010, 'doc', 'PV操作经典例题讲解', 'https://example.com/pv.pdf', False),
        # 计算机网络 (1011)
        (1011, 'tool', 'Wireshark 抓包工具', 'https://www.wireshark.org/', False),
        # Web前端 (1015)
        (1015, 'video', 'Vue3 组合式API实战教程', 'https://bilibili.com/vue3', False),
        (1015, 'doc', '前端面试题汇总2024版', 'https://github.com/interview', False),
        # 算法 (1016)
        (1016, 'doc', 'LeetCode Hot 100 题解', 'https://leetcode.cn', False),
    ]
    
    try:
        sort_order = 0
        last_course_id = None
        
        for course_id, resource_type, resource_intro, url, is_upload in resources_data:
            # 如果课程ID改变，重置排序
            if last_course_id != course_id:
                sort_order = 0
                last_course_id = course_id
            
            if is_upload:
                cursor.execute("""
                    INSERT IGNORE INTO course_resources_upload (course_id, resource_intro, resource_upload, sort_order, created_at)
                    VALUES (%s, %s, %s, %s, NOW())
                """, (course_id, resource_intro, url, sort_order))
            else:
                cursor.execute("""
                    INSERT IGNORE INTO course_resources_web (course_id, resource_intro, resource_url, sort_order, created_at)
                    VALUES (%s, %s, %s, %s, NOW())
                """, (course_id, resource_intro, url, sort_order))
            
            sort_order += 1
        
        conn.commit()
        print(f"成功导入 {len(resources_data)} 个课程资源")
    except Error as e:
        print(f"导入资源数据错误: {e}")
        conn.rollback()
    finally:
        cursor.close()

def import_comments(conn, user_ids):
    """导入课程评论数据"""
    cursor = conn.cursor()
    
    comments_data = [
        # (course_id, content, love_count, comment_time)
        (1001, '这也太难了吧，完全听不懂极限定义...', 50, '2023-09-10 10:00:00'),
        (1001, '陈老师讲得很细，课后习题一定要自己做。', 12, '2023-09-12 14:00:00'),
        (1002, '指针那里真的晕了，求大神指点！', 8, '2023-10-01 09:00:00'),
        (1002, '多写代码，多调试，C语言其实不难。', 25, '2023-10-02 11:00:00'),
        (1003, '这门课比较轻松，主要了解概念。', 5, '2023-09-15 16:00:00'),
        (1003, '王老师人很好，期末给分高。', 30, '2023-12-20 10:00:00'),
        (1004, '重积分算到手断...', 44, '2024-03-01 10:00:00'),
        (1004, '空间几何要注意画图辅助理解。', 10, '2024-03-05 12:00:00'),
        (1005, '真值表画错了，痛失10分。', 2, '2024-04-01 10:00:00'),
        (1005, '图论部分挺有意思的。', 6, '2024-04-10 14:00:00'),
        (1006, '面向对象思想真的很重要，比面向过程好维护多了。', 100, '2024-05-01 09:00:00'),
        (1006, '推荐直接用 IDEA，别用 Eclipse 了。', 200, '2024-05-02 10:00:00'),
        (1007, '数据结构是考研专业课最难的一门，必须死磕。', 300, '2023-10-10 20:00:00'),
        (1007, 'KMP算法看了一周才看懂...', 50, '2023-10-15 21:00:00'),
        (1008, '补码反码原码搞得头晕。', 15, '2023-11-01 10:00:00'),
        (1008, '流水线技术是提升CPU性能的关键。', 20, '2023-11-05 10:00:00'),
        (1009, 'Python 真的简洁，写起来很爽。', 40, '2023-09-20 14:00:00'),
        (1009, '爬虫作业要注意设置 User-Agent，不然会被封IP。', 60, '2023-09-25 15:00:00'),
        (1010, '哲学家就餐问题很有趣。', 33, '2024-03-10 10:00:00'),
        (1010, '一定要理解虚拟内存的概念。', 22, '2024-03-12 11:00:00'),
        (1011, '三次握手和四次挥手背下来，面试必问。', 120, '2024-04-20 09:00:00'),
        (1011, '子网掩码计算有点绕。', 10, '2024-04-22 10:00:00'),
        (1012, '范式理论虽然枯燥，但是设计数据库必须遵守。', 45, '2024-05-01 14:00:00'),
        (1012, '左连接右连接傻傻分不清楚。', 20, '2024-05-05 16:00:00'),
        (1013, '如何在 Vim 中退出？在线等，挺急的。', 999, '2024-03-15 12:00:00'),
        (1013, 'rm -rf /* 慎用！！', 500, '2024-03-16 13:00:00'),
        (1014, '文档写得头大，还是写代码爽。', 15, '2023-10-10 09:00:00'),
        (1014, '敏捷开发确实比瀑布模型灵活。', 20, '2023-10-12 10:00:00'),
        (1015, 'Vue3 真的好用，Composition API 很香。', 88, '2023-11-20 10:00:00'),
        (1015, '垂直居中为什么这么难调？？', 200, '2023-11-21 11:00:00'),
        (1016, '动态规划的核心是状态转移方程。', 150, '2023-12-01 20:00:00'),
        (1016, '刷了200题，感觉稍微入门了。', 60, '2023-12-05 21:00:00'),
        (1017, '开发不仅要会写代码，还要会写测试用例。', 30, '2024-03-20 10:00:00'),
        (1017, '自动化测试是趋势。', 25, '2024-03-22 11:00:00'),
        (1018, 'AS 模拟器太吃内存了，建议用真机调试。', 40, '2024-04-10 15:00:00'),
        (1018, 'Kotlin 语法糖很甜。', 35, '2024-04-12 16:00:00'),
        (1019, '数学不好学 AI 有点吃力啊。', 55, '2024-05-15 10:00:00'),
        (1019, 'PyTorch 比 TensorFlow 容易上手。', 45, '2024-05-18 11:00:00'),
        (1020, '这就是传说中的劝退课吗？太难了！', 200, '2024-06-01 10:00:00'),
        (1020, '写出一个编译器后的成就感无与伦比。', 10, '2024-06-05 12:00:00'),
    ]
    
    try:
        # 循环分配用户ID
        user_index = 0
        
        for course_id, content, love_count, comment_time in comments_data:
            user_id = user_ids[user_index % len(user_ids)]
            user_index += 1
            
            cursor.execute("""
                INSERT INTO comments (resource_type, resource_id, user_id, content, love_count, created_at, parent_id, reply_total)
                VALUES ('course', %s, %s, %s, %s, %s, NULL, 0)
            """, (course_id, user_id, content, love_count, comment_time))
        
        conn.commit()
        print(f"成功导入 {len(comments_data)} 条评论")
    except Error as e:
        print(f"导入评论数据错误: {e}")
        conn.rollback()
    finally:
        cursor.close()

def main():
    """主函数"""
    print("开始导入课程详情数据...")
    
    conn = connect_db()
    if not conn:
        print("无法连接到数据库，请检查配置")
        return
    
    try:
        # 添加description字段（如果需要）
        add_description_column(conn)
        
        # 获取用户ID列表
        user_ids = get_user_ids(conn)
        if not user_ids:
            print("警告: 数据库中没有任何用户，评论将无法关联用户。请先创建用户。")
            user_ids = [1]  # 假设至少有一个用户ID为1
        
        print(f"找到 {len(user_ids)} 个用户，将用于分配评论")
        
        # 导入课程数据
        print("\n导入课程数据...")
        import_courses(conn)
        
        # 导入资源数据
        print("\n导入资源数据...")
        import_resources(conn)
        
        # 导入评论数据
        print("\n导入评论数据...")
        import_comments(conn, user_ids)
        
        print("\n数据导入完成！")
        
    except Error as e:
        print(f"导入过程中发生错误: {e}")
        conn.rollback()
    finally:
        if conn.is_connected():
            conn.close()
            print("数据库连接已关闭")

if __name__ == "__main__":
    main()

