# -*- coding: utf-8 -*-
"""验证导入的项目数据，确保没有乱码"""
import mysql.connector
import sys

# 确保输出使用UTF-8编码
if sys.stdout.encoding != 'utf-8':
    sys.stdout.reconfigure(encoding='utf-8')

try:
    conn = mysql.connector.connect(
        host='localhost',
        user='root',
        password='Wan05609',
        database='softeng',
        charset='utf8mb4',
        collation='utf8mb4_unicode_ci'
    )
    cursor = conn.cursor()
    
    print('验证项目数据（前5条）:')
    print('=' * 80)
    cursor.execute('SELECT project_id, name, category, description, status FROM projects ORDER BY project_id LIMIT 5')
    projects = cursor.fetchall()
    for row in projects:
        project_id, name, category, description, status = row
        desc = (description[:50] + '...') if description and len(description) > 50 else (description or '')
        print(f'ID: {project_id}, 名称: {name}, 分类: {category}, 状态: {status}')
        print(f'  简介: {desc}')
        print()
    
    print('=' * 80)
    print('\n验证技术栈数据:')
    cursor.execute('SELECT project_id, tech FROM project_tech_stack ORDER BY project_id LIMIT 10')
    techs = cursor.fetchall()
    for project_id, tech in techs:
        print(f'  项目ID={project_id}, 技术栈={tech}')
    
    print('\n' + '=' * 80)
    print('\n验证评论数据（前3条）:')
    cursor.execute('SELECT comment_id, resource_id, content, love_count FROM comments WHERE resource_type="project" ORDER BY comment_id LIMIT 3')
    comments = cursor.fetchall()
    for row in comments:
        comment_id, resource_id, content, love_count = row
        content_preview = (content[:40] + '...') if content and len(content) > 40 else (content or '')
        print(f'评论ID: {comment_id}, 项目ID: {resource_id}, 点赞数: {love_count}')
        print(f'  内容: {content_preview}')
        print()
    
    print('=' * 80)
    print('\n数据统计:')
    cursor.execute('SELECT COUNT(*) FROM projects')
    project_count = cursor.fetchone()[0]
    cursor.execute('SELECT COUNT(*) FROM project_tech_stack')
    tech_count = cursor.fetchone()[0]
    cursor.execute('SELECT COUNT(*) FROM project_images')
    image_count = cursor.fetchone()[0]
    cursor.execute('SELECT COUNT(*) FROM project_authors')
    author_count = cursor.fetchone()[0]
    cursor.execute('SELECT COUNT(*) FROM comments WHERE resource_type="project"')
    comment_count = cursor.fetchone()[0]
    
    print(f'  项目数: {project_count}')
    print(f'  技术栈记录数: {tech_count}')
    print(f'  项目图片数: {image_count}')
    print(f'  项目作者数: {author_count}')
    print(f'  评论数: {comment_count}')
    print('=' * 80)

except mysql.connector.Error as err:
    print(f'数据库查询错误: {err}')
finally:
    if 'conn' in locals() and conn.is_connected():
        cursor.close()
        conn.close()

