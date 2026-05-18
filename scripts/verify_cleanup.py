#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""验证数据清理效果"""

import mysql.connector
import os

# 从环境变量获取数据库配置
db_user = os.getenv('DB_USER', 'root')
db_password = os.getenv('DB_PASSWORD', 'Wan05609')
db_host = os.getenv('DB_HOST', '127.0.0.1')
db_name = os.getenv('DB_NAME', 'softeng')

try:
    conn = mysql.connector.connect(
        user=db_user,
        password=db_password,
        host=db_host,
        database=db_name
    )
    cursor = conn.cursor()
    
    print("=" * 50)
    print("数据清理效果验证")
    print("=" * 50)
    
    # 检查"豆包"工具
    cursor.execute("SELECT COUNT(*) FROM tools WHERE resource_name = '豆包'")
    doubao_count = cursor.fetchone()[0]
    print(f"\n1. 名为'豆包'的工具数量: {doubao_count}")
    
    # 检查用户提交的工具
    cursor.execute("SELECT COUNT(*) FROM tools WHERE submitter_id IS NOT NULL")
    tool_submissions = cursor.fetchone()[0]
    print(f"2. 用户提交的工具数量: {tool_submissions}")
    
    # 检查用户提交的课程
    cursor.execute("""
        SELECT COUNT(*) FROM courses 
        WHERE course_id IN (SELECT DISTINCT course_id FROM course_contributors)
    """)
    course_submissions = cursor.fetchone()[0]
    print(f"3. 用户提交的课程数量: {course_submissions}")
    
    # 检查用户提交的项目
    cursor.execute("""
        SELECT COUNT(*) FROM projects 
        WHERE project_id IN (SELECT DISTINCT project_id FROM project_authors)
    """)
    project_submissions = cursor.fetchone()[0]
    print(f"4. 用户提交的项目数量: {project_submissions}")
    
    # 检查关联数据
    cursor.execute("SELECT COUNT(*) FROM tool_images WHERE tool_id IN (SELECT resource_id FROM tools WHERE resource_name = '豆包')")
    doubao_images = cursor.fetchone()[0]
    print(f"5. '豆包'工具的图片数量: {doubao_images}")
    
    cursor.execute("SELECT COUNT(*) FROM tool_tags WHERE tool_id IN (SELECT resource_id FROM tools WHERE resource_name = '豆包')")
    doubao_tags = cursor.fetchone()[0]
    print(f"6. '豆包'工具的标签数量: {doubao_tags}")
    
    print("\n" + "=" * 50)
    if doubao_count == 0 and tool_submissions == 0 and course_submissions == 0 and project_submissions == 0:
        print("[SUCCESS] 验证通过：所有数据已成功清理！")
    else:
        print("[WARNING] 警告：仍有部分数据未清理")
    print("=" * 50)
    
    cursor.close()
    conn.close()
    
except mysql.connector.Error as err:
    print(f"数据库错误: {err}")
except Exception as e:
    print(f"发生错误: {e}")

