#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import mysql.connector

db_config = {
    'host': 'localhost',
    'user': 'root',
    'password': 'Wan05609',
    'database': 'softeng',
    'charset': 'utf8mb4'
}

courses_data = [
    (201, '高等数学 II', 'MATH1002', '1-2', 5, 'course', 'https://placehold.co/600x450/3b82f6/ffffff?text=MATH', 0, 38, 15),
    (202, '线性代数', 'MATH1003', '1-2', 3, 'course', 'https://placehold.co/600x450/3b82f6/ffffff?text=MATH', 0, 88, 20),
    (203, '离散数学', 'CS1002', '1-2', 4, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 150, 35),
    (204, '体育2', 'PE1002', '1-2', 1, 'course', 'https://placehold.co/600x450/ef4444/ffffff?text=PE', 0, 180, 5),
    (301, '数据结构与算法', 'CS2001', '2-1', 5, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 230, 56),
    (302, '计算机组成原理', 'CS2002', '2-1', 4, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 95, 30),
    (303, 'Python应用开发', 'CS2005', '2-1', 2, 'course', 'https://placehold.co/600x450/14b8a6/ffffff?text=PYTHON', 0, 67, 18),
    (401, '操作系统', 'CS2003', '2-2', 4, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 180, 42),
    (402, '计算机网络', 'CS2004', '2-2', 4, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 160, 38),
    (501, '计算机网络', 'CS3001', '3-1', 3, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 120, 18),
    (502, '数据库系统概论', 'CS3002', '3-1', 3, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=CS', 0, 140, 25),
    (601, '软件工程导论', 'CS3003', '3-2', 2, 'course', 'https://placehold.co/600x450/14b8a6/ffffff?text=SE', 0, 80, 22),
    (701, '人工智能导论', 'CS4001', '4-1', 2, 'course', 'https://placehold.co/600x450/14b8a6/ffffff?text=AI', 0, 90, 15),
    (702, '毕业设计', 'CS4002', '4-2', 6, 'course', 'https://placehold.co/600x450/10b981/ffffff?text=GRAD', 0, 50, 10),
]

teachers_data = [
    (201, '张老师'), (202, '赵老师'), (203, '钱老师'), (204, '张老师'),
    (301, '孙老师'), (302, '周老师'), (303, '吴老师'),
    (401, '郑老师'), (402, '冯老师'),
    (501, '马老师'), (502, '刘老师'),
    (601, '毛老师'),
    (701, '林老师'),
    (702, '何老师'),
]

categories_data = [
    (201, '公必'), (202, '公必'), (203, '专必'), (204, '公必'),
    (301, '专必'), (302, '专必'), (303, '专选'),
    (401, '专必'), (402, '专必'),
    (501, '专必'), (502, '专必'),
    (601, '专选'),
    (701, '专选'),
    (702, '专必'),
]

try:
    conn = mysql.connector.connect(**db_config)
    cursor = conn.cursor()
    
    # 插入课程
    sql_course = "INSERT IGNORE INTO courses (course_id, name, code, semester, credit, resource_type, cover, views, loves, collections) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"
    cursor.executemany(sql_course, courses_data)
    
    # 插入教师
    sql_teacher = "INSERT IGNORE INTO course_teachers (course_id, teacher_name) VALUES (%s, %s)"
    cursor.executemany(sql_teacher, teachers_data)
    
    # 插入分类
    sql_category = "INSERT IGNORE INTO course_categories (course_id, category) VALUES (%s, %s)"
    cursor.executemany(sql_category, categories_data)
    
    conn.commit()
    print(f"成功插入 {len(courses_data)} 条课程数据")
    cursor.close()
    conn.close()
except Exception as e:
    print(f"错误: {e}")

