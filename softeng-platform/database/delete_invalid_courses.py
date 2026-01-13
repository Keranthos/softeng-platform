#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
删除不属于专必/专选/公必/公选的课程
"""

import mysql.connector
from mysql.connector import Error
import sys

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

def get_courses_to_delete(conn):
    """获取需要删除的课程ID列表"""
    cursor = conn.cursor()
    courses_to_delete = []
    
    try:
        # 查询所有课程及其分类
        query = """
            SELECT c.course_id, c.name, 
                   GROUP_CONCAT(cc.category ORDER BY cc.category SEPARATOR ', ') as categories
            FROM courses c
            LEFT JOIN course_categories cc ON c.course_id = cc.course_id
            GROUP BY c.course_id, c.name
        """
        cursor.execute(query)
        all_courses = cursor.fetchall()
        
        print(f"总共找到 {len(all_courses)} 门课程\n")
        
        for course_id, name, categories in all_courses:
            # 如果没有分类，或者分类中没有任何一个合法分类，则标记为删除
            if not categories:
                print(f"课程ID {course_id}: {name} - 无分类 (标记删除)")
                courses_to_delete.append(course_id)
            else:
                # 检查分类列表是否包含合法分类
                category_list = [cat.strip() for cat in categories.split(',')]
                has_valid_category = any(cat in VALID_CATEGORIES for cat in category_list)
                
                if not has_valid_category:
                    print(f"课程ID {course_id}: {name} - 分类: {categories} (标记删除)")
                    courses_to_delete.append(course_id)
                else:
                    print(f"课程ID {course_id}: {name} - 分类: {categories} (保留)")
        
        return courses_to_delete
        
    except Error as e:
        print(f"查询课程错误: {e}")
        return []
    finally:
        cursor.close()

def delete_courses(conn, course_ids):
    """删除指定的课程（由于外键约束，会自动删除相关的教师、分类、资源等）"""
    if not course_ids:
        print("\n没有需要删除的课程")
        return
    
    cursor = conn.cursor()
    
    try:
        # 先删除相关的评论（comments表中的数据）
        print(f"\n准备删除 {len(course_ids)} 门课程...")
        
        # 删除课程评论
        placeholders = ','.join(['%s'] * len(course_ids))
        delete_comments_query = f"""
            DELETE FROM comments 
            WHERE resource_type = 'course' 
            AND resource_id IN ({placeholders})
        """
        cursor.execute(delete_comments_query, course_ids)
        comments_deleted = cursor.rowcount
        print(f"已删除 {comments_deleted} 条课程评论")
        
        # 删除课程相关的收藏（collections表中的数据）
        delete_collections_query = f"""
            DELETE FROM collections 
            WHERE resource_type = 'course' 
            AND resource_id IN ({placeholders})
        """
        cursor.execute(delete_collections_query, course_ids)
        collections_deleted = cursor.rowcount
        print(f"已删除 {collections_deleted} 条收藏记录")
        
        # 删除课程相关的点赞（likes表中的数据）
        delete_likes_query = f"""
            DELETE FROM likes 
            WHERE resource_type = 'course' 
            AND resource_id IN ({placeholders})
        """
        cursor.execute(delete_likes_query, course_ids)
        likes_deleted = cursor.rowcount
        print(f"已删除 {likes_deleted} 条点赞记录")
        
        # 删除课程（由于外键约束，会自动删除相关的教师、分类、资源、贡献者等）
        delete_courses_query = f"""
            DELETE FROM courses 
            WHERE course_id IN ({placeholders})
        """
        cursor.execute(delete_courses_query, course_ids)
        courses_deleted = cursor.rowcount
        print(f"已删除 {courses_deleted} 门课程")
        
        # 提交事务
        conn.commit()
        print(f"\n成功删除 {courses_deleted} 门课程及其所有关联数据！")
        
    except Error as e:
        print(f"删除课程错误: {e}")
        conn.rollback()
        raise
    finally:
        cursor.close()

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
    
    try:
        # 获取需要删除的课程
        courses_to_delete = get_courses_to_delete(conn)
        
        if not courses_to_delete:
            print("\n所有课程都符合要求，无需删除")
            return
        
        # 确认删除
        print(f"\n将要删除以下 {len(courses_to_delete)} 门课程:")
        print(f"课程ID列表: {courses_to_delete}")
        
        confirm = input("\n确认删除? (输入 'yes' 确认，其他任意键取消): ")
        if confirm.lower() != 'yes':
            print("操作已取消")
            return
        
        # 执行删除
        delete_courses(conn, courses_to_delete)
        
        print("\n操作完成！")
        
    except Error as e:
        print(f"\n操作过程中发生错误: {e}")
        conn.rollback()
    except KeyboardInterrupt:
        print("\n\n操作被用户中断")
        conn.rollback()
    finally:
        if conn.is_connected():
            conn.close()
            print("数据库连接已关闭")

if __name__ == "__main__":
    main()
