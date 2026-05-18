#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
删除不属于专必/专选/公必/公选的课程（简化版本）
"""

import mysql.connector
from mysql.connector import Error
import sys
import io

# 设置输出编码为UTF-8
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

# 数据库配置
DB_CONFIG = {
    'host': 'localhost',
    'user': 'root',
    'password': 'Wan05609',
    'database': 'softeng',
    'charset': 'utf8mb4',
    'collation': 'utf8mb4_unicode_ci'
}

# 合法的课程分类
VALID_CATEGORIES = ['专必', '专选', '公必', '公选']

def connect_db():
    """连接数据库"""
    try:
        conn = mysql.connector.connect(**DB_CONFIG)
        return conn
    except Error as e:
        print(f"数据库连接错误: {e}")
        return None

def main():
    """主函数"""
    print("=" * 60)
    print("删除不属于专必/专选/公必/公选的课程")
    print("=" * 60)
    print(f"合法分类: {', '.join(VALID_CATEGORIES)}\n")
    
    conn = connect_db()
    if not conn:
        print("无法连接到数据库")
        sys.exit(1)
    
    cursor = conn.cursor()
    
    try:
        # 查询需要删除的课程
        query = """
            SELECT c.course_id, c.name, 
                   GROUP_CONCAT(cc.category ORDER BY cc.category SEPARATOR ', ') as categories
            FROM courses c
            LEFT JOIN course_categories cc ON c.course_id = cc.course_id
            GROUP BY c.course_id, c.name
            HAVING categories IS NULL 
                OR (categories NOT LIKE '%专必%' 
                    AND categories NOT LIKE '%专选%' 
                    AND categories NOT LIKE '%公必%' 
                    AND categories NOT LIKE '%公选%')
        """
        cursor.execute(query)
        courses_to_delete = cursor.fetchall()
        
        if not courses_to_delete:
            print("所有课程都符合要求（都属于专必/专选/公必/公选），无需删除")
            print("\n检查所有课程分类...")
            
            # 列出所有课程及其分类
            cursor.execute("""
                SELECT c.course_id, c.name, 
                       COALESCE(GROUP_CONCAT(cc.category ORDER BY cc.category SEPARATOR ', '), '无分类') as categories
                FROM courses c
                LEFT JOIN course_categories cc ON c.course_id = cc.course_id
                GROUP BY c.course_id, c.name
                ORDER BY c.course_id
            """)
            all_courses = cursor.fetchall()
            
            for course_id, name, categories in all_courses:
                print(f"  课程ID {course_id}: {name} - 分类: {categories}")
            
            return
        
        # 显示需要删除的课程
        print(f"找到 {len(courses_to_delete)} 门需要删除的课程:\n")
        course_ids = []
        for course_id, name, categories in courses_to_delete:
            print(f"  课程ID {course_id}: {name} - 分类: {categories or '无分类'}")
            course_ids.append(course_id)
        
        # 确认删除
        print(f"\n将要删除以上 {len(course_ids)} 门课程及其所有关联数据（评论、收藏、点赞、资源、教师、分类等）")
        confirm = input("\n确认删除? (输入 'yes' 确认，其他任意键取消): ")
        
        if confirm.lower() != 'yes':
            print("操作已取消")
            return
        
        # 执行删除
        placeholders = ','.join(['%s'] * len(course_ids))
        
        # 删除课程评论
        cursor.execute(f"""
            DELETE FROM comments 
            WHERE resource_type = 'course' 
            AND resource_id IN ({placeholders})
        """, course_ids)
        comments_deleted = cursor.rowcount
        print(f"已删除 {comments_deleted} 条课程评论")
        
        # 删除课程相关的收藏
        cursor.execute(f"""
            DELETE FROM collections 
            WHERE resource_type = 'course' 
            AND resource_id IN ({placeholders})
        """, course_ids)
        collections_deleted = cursor.rowcount
        print(f"已删除 {collections_deleted} 条收藏记录")
        
        # 删除课程相关的点赞
        cursor.execute(f"""
            DELETE FROM likes 
            WHERE resource_type = 'course' 
            AND resource_id IN ({placeholders})
        """, course_ids)
        likes_deleted = cursor.rowcount
        print(f"已删除 {likes_deleted} 条点赞记录")
        
        # 删除课程（由于外键约束，会自动删除相关的教师、分类、资源、贡献者等）
        cursor.execute(f"""
            DELETE FROM courses 
            WHERE course_id IN ({placeholders})
        """, course_ids)
        courses_deleted = cursor.rowcount
        print(f"已删除 {courses_deleted} 门课程及其关联数据")
        
        # 提交事务
        conn.commit()
        print(f"\n成功！已删除 {courses_deleted} 门课程及其所有关联数据")
        
        # 显示剩余课程数量
        cursor.execute("SELECT COUNT(*) FROM courses")
        remaining = cursor.fetchone()[0]
        print(f"数据库中剩余课程: {remaining} 门")
        
    except Error as e:
        print(f"\n错误: {e}")
        conn.rollback()
    except KeyboardInterrupt:
        print("\n\n操作被用户中断")
        conn.rollback()
    finally:
        cursor.close()
        if conn.is_connected():
            conn.close()

if __name__ == "__main__":
    main()

