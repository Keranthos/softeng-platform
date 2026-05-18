# -*- coding: utf-8 -*-
"""
导入项目模拟数据到数据库
确保使用UTF-8编码以避免乱码
"""
import mysql.connector
import sys
from datetime import datetime

# 确保输出使用UTF-8编码
if sys.stdout.encoding != 'utf-8':
    sys.stdout.reconfigure(encoding='utf-8')

# 模拟数据（从 mockData.js 中提取）
mock_projects = [
    {
        'id': 1,
        'name': 'Vue3 Admin Template',
        'description': '基于Vue3 + Element Plus的后台管理模板',
        'details': '这是一个功能完善的后台管理模板，包含用户管理、权限控制、图表展示等模块。采用最新的Vue3 Composition API，代码结构清晰，易于二次开发。',
        'logo': 'https://vuejs.org/images/logo.png',
        'coverImage': 'https://images.unsplash.com/photo-1551650975-87deedd944c3?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
        'category': '前端项目',
        'githubUrl': 'https://github.com/vuejs/vue',
        'demoUrl': 'https://vuejs.org',
        'technologies': ['Vue', 'JavaScript', 'Element Plus', 'Vue Router', 'Pinia'],
        'tags': ['frontend', 'vue', 'admin', 'template'],
        'license': 'MIT',
        'status': 'active',
        'views': 1250,
        'stars': 320,
        'contributors': [
            {'id': 1, 'name': '张三', 'avatar': 'https://i.pravatar.cc/150?img=1', 'role': '主要开发者'},
            {'id': 2, 'name': '李四', 'avatar': 'https://i.pravatar.cc/150?img=2', 'role': 'UI设计师'}
        ],
        'createdAt': '2024-01-15T10:30:00Z'
    },
    {
        'id': 2,
        'name': 'React E-commerce',
        'description': '使用React构建的电子商务平台',
        'details': '完整的电商解决方案，包含商品展示、购物车、支付接口、用户订单管理等功能。采用微服务架构，支持高并发访问。',
        'logo': 'https://reactjs.org/logo-og.png',
        'coverImage': 'https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
        'category': '全栈项目',
        'githubUrl': 'https://github.com/facebook/react',
        'demoUrl': 'https://reactjs.org',
        'technologies': ['React', 'Node.js', 'MongoDB', 'Express', 'Redux'],
        'tags': ['fullstack', 'react', 'ecommerce', 'nodejs'],
        'license': 'MIT',
        'status': 'active',
        'views': 890,
        'stars': 210,
        'contributors': [
            {'id': 3, 'name': '王五', 'avatar': 'https://i.pravatar.cc/150?img=3', 'role': '全栈工程师'},
            {'id': 4, 'name': '赵六', 'avatar': 'https://i.pravatar.cc/150?img=4', 'role': '后端开发'}
        ],
        'createdAt': '2024-02-20T14:45:00Z'
    },
    {
        'id': 3,
        'name': 'Python Data Analyzer',
        'description': 'Python数据分析和可视化工具',
        'details': '提供数据清洗、统计分析、可视化图表生成等功能。支持多种数据格式导入，内置机器学习算法，可进行预测分析。',
        'logo': 'https://www.python.org/static/community_logos/python-logo.png',
        'coverImage': 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
        'category': '后端项目',
        'githubUrl': 'https://github.com/python/cpython',
        'demoUrl': None,
        'technologies': ['Python', 'Pandas', 'NumPy', 'Matplotlib', 'Scikit-learn'],
        'tags': ['backend', 'python', 'data-science', 'ml'],
        'license': 'GPL-3.0',
        'status': 'active',
        'views': 670,
        'stars': 180,
        'contributors': [
            {'id': 5, 'name': '钱七', 'avatar': 'https://i.pravatar.cc/150?img=5', 'role': '数据科学家'}
        ],
        'createdAt': '2024-03-10T09:15:00Z'
    },
    {
        'id': 4,
        'name': 'Flutter Travel App',
        'description': 'Flutter跨平台旅行应用',
        'details': '包含旅行计划、景点推荐、酒店预订、地图导航等功能的移动应用。支持iOS和Android平台，使用Firebase作为后端。',
        'logo': 'https://storage.googleapis.com/cms-storage-bucket/0dbfcc7a59cd1cf16282.png',
        'coverImage': 'https://images.unsplash.com/photo-1544551763-46a013bb70d5?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
        'category': '移动应用',
        'githubUrl': 'https://github.com/flutter/flutter',
        'demoUrl': 'https://flutter.dev',
        'technologies': ['Flutter', 'Dart', 'Firebase', 'Google Maps'],
        'tags': ['mobile', 'flutter', 'travel', 'firebase'],
        'license': 'BSD-3-Clause',
        'status': 'active',
        'views': 540,
        'stars': 150,
        'contributors': [
            {'id': 6, 'name': '孙八', 'avatar': 'https://i.pravatar.cc/150?img=6', 'role': '移动开发'},
            {'id': 7, 'name': '周九', 'avatar': 'https://i.pravatar.cc/150?img=7', 'role': 'UI设计师'}
        ],
        'createdAt': '2024-04-05T16:20:00Z'
    },
    {
        'id': 5,
        'name': 'Electron Music Player',
        'description': '使用Electron开发的桌面音乐播放器',
        'details': '支持多种音频格式，具有播放列表管理、歌词显示、音效调节等功能。界面美观，操作流畅，支持快捷键操作。',
        'logo': 'https://electronjs.org/images/electron-logo.svg',
        'coverImage': 'https://images.unsplash.com/photo-1511379938547-c1f69419868d?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
        'category': '桌面应用',
        'githubUrl': 'https://github.com/electron/electron',
        'demoUrl': None,
        'technologies': ['Electron', 'JavaScript', 'HTML', 'CSS', 'Node.js'],
        'tags': ['desktop', 'electron', 'music', 'player'],
        'license': 'MIT',
        'status': 'archived',
        'views': 320,
        'stars': 95,
        'contributors': [
            {'id': 8, 'name': '吴十', 'avatar': 'https://i.pravatar.cc/150?img=8', 'role': '桌面开发'}
        ],
        'createdAt': '2023-12-01T11:00:00Z'
    },
    {
        'id': 6,
        'name': 'Spring Boot Microservices',
        'description': '基于Spring Boot的微服务架构示例',
        'details': '完整的微服务解决方案，包含服务注册与发现、配置中心、API网关、熔断机制等。使用Docker容器化部署。',
        'logo': 'https://spring.io/images/spring-logo-9146a4d3298760c2e7e49595184e1975.svg',
        'coverImage': 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
        'category': '后端项目',
        'githubUrl': 'https://github.com/spring-projects/spring-boot',
        'demoUrl': None,
        'technologies': ['Java', 'Spring Boot', 'Spring Cloud', 'Docker', 'Kubernetes'],
        'tags': ['backend', 'java', 'microservices', 'spring'],
        'license': 'Apache-2.0',
        'status': 'active',
        'views': 780,
        'stars': 240,
        'contributors': [
            {'id': 9, 'name': '郑十一', 'avatar': 'https://i.pravatar.cc/150?img=9', 'role': '后端架构师'},
            {'id': 10, 'name': '王十二', 'avatar': 'https://i.pravatar.cc/150?img=10', 'role': 'DevOps工程师'}
        ],
        'createdAt': '2024-03-25T13:30:00Z'
    },
    {
        'id': 7,
        'name': 'Vue3 UI Library',
        'description': '企业级Vue3 UI组件库',
        'details': '包含丰富的UI组件，如表格、表单、图表、弹窗等。完全使用TypeScript开发，提供完整的类型定义。支持主题定制和按需加载。',
        'logo': 'https://vuejs.org/images/logo.png',
        'coverImage': 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
        'category': '前端项目',
        'githubUrl': 'https://github.com/vuejs/vue',
        'demoUrl': 'https://vuejs.org',
        'technologies': ['Vue', 'TypeScript', 'Vite', 'SCSS', 'Jest'],
        'tags': ['frontend', 'vue', 'ui', 'typescript', 'library'],
        'license': 'MIT',
        'status': 'active',
        'views': 920,
        'stars': 310,
        'contributors': [
            {'id': 1, 'name': '张三', 'avatar': 'https://i.pravatar.cc/150?img=1', 'role': '主要开发者'},
            {'id': 11, 'name': '刘十三', 'avatar': 'https://i.pravatar.cc/150?img=11', 'role': '测试工程师'}
        ],
        'createdAt': '2024-02-10T08:45:00Z'
    },
    {
        'id': 8,
        'name': 'Node.js API Framework',
        'description': '高性能Node.js REST API框架',
        'details': '基于Express的增强框架，内置JWT认证、权限控制、请求验证、日志记录、性能监控等功能。提供CLI工具快速生成项目结构。',
        'logo': 'https://nodejs.org/static/images/logo.svg',
        'coverImage': 'https://images.unsplash.com/photo-1618401471353-b98afee0b2eb?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80',
        'category': '后端项目',
        'githubUrl': 'https://github.com/nodejs/node',
        'demoUrl': None,
        'technologies': ['Node.js', 'Express', 'JWT', 'MongoDB', 'Redis'],
        'tags': ['backend', 'nodejs', 'api', 'framework', 'express'],
        'license': 'MIT',
        'status': 'active',
        'views': 650,
        'stars': 175,
        'contributors': [
            {'id': 12, 'name': '陈十四', 'avatar': 'https://i.pravatar.cc/150?img=12', 'role': '后端开发'}
        ],
        'createdAt': '2024-01-30T15:20:00Z'
    }
]

# 项目评论数据
mock_comments = [
    {
        'id': 1,
        'projectId': 1,
        'userId': 1,
        'username': 'user1',
        'nickname': '开发者小明',
        'avatar': 'https://i.pravatar.cc/150?img=1',
        'content': '这个项目结构清晰，代码规范，非常适合学习Vue3！',
        'likes': 12,
        'createdAt': '2024-04-01T10:30:00Z'
    },
    {
        'id': 2,
        'projectId': 1,
        'userId': 2,
        'username': 'user2',
        'nickname': '设计师小红',
        'avatar': 'https://i.pravatar.cc/150?img=2',
        'content': 'UI设计很美观，交互体验很好！',
        'likes': 8,
        'createdAt': '2024-04-02T14:20:00Z'
    },
    {
        'id': 3,
        'projectId': 2,
        'userId': 3,
        'username': 'user3',
        'nickname': '全栈工程师',
        'avatar': 'https://i.pravatar.cc/150?img=3',
        'content': '电商功能很完整，学习了！',
        'likes': 15,
        'createdAt': '2024-04-03T09:15:00Z'
    }
]

try:
    # 连接数据库，指定字符集为utf8mb4
    conn = mysql.connector.connect(
        host='localhost',
        user='root',
        password='Wan05609',
        database='softeng',
        charset='utf8mb4',
        collation='utf8mb4_unicode_ci'
    )
    cursor = conn.cursor()
    
    print('开始导入项目数据...\n')
    
    # 获取现有用户ID列表（用于评论分配）
    cursor.execute('SELECT id FROM users ORDER BY id')
    user_ids = [row[0] for row in cursor.fetchall()]
    if not user_ids:
        print('警告：数据库中没有用户，评论将无法分配用户ID')
        user_ids = [1]  # 使用默认ID
    
    # 导入项目
    for project in mock_projects:
        project_id = project['id']
        
        # 解析创建时间
        created_at = datetime.fromisoformat(project['createdAt'].replace('Z', '+00:00'))
        
        # 映射status：active -> approved, archived -> approved (归档的也显示为已审核)
        status = 'approved' if project['status'] in ['active', 'archived'] else 'pending'
        
        # 插入项目主表
        insert_project_sql = """
            INSERT INTO projects (
                project_id, resource_type, name, description, detail, 
                github_url, category, cover, views, loves, collections, 
                status, created_at
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
            ON DUPLICATE KEY UPDATE
                resource_type = VALUES(resource_type),
                name = VALUES(name),
                description = VALUES(description),
                detail = VALUES(detail),
                github_url = VALUES(github_url),
                category = VALUES(category),
                cover = VALUES(cover),
                views = VALUES(views),
                loves = VALUES(loves),
                collections = VALUES(collections),
                status = VALUES(status),
                created_at = VALUES(created_at)
        """
        
        cursor.execute(insert_project_sql, (
            project_id,
            'project',
            project['name'],
            project['description'],
            project['details'],
            project['githubUrl'],
            project['category'],
            project['coverImage'] or project.get('logo', ''),
            project['views'],
            project['stars'],  # stars -> loves
            0,  # collections 初始为0
            status,
            created_at
        ))
        
        project_inserted_id = cursor.lastrowid if cursor.lastrowid else project_id
        print(f'[OK] 已导入项目: {project["name"]} (ID: {project_id})')
        
        # 插入技术栈
        if project.get('technologies'):
            for tech in project['technologies']:
                try:
                    cursor.execute(
                        'INSERT INTO project_tech_stack (project_id, tech) VALUES (%s, %s) ON DUPLICATE KEY UPDATE tech = VALUES(tech)',
                        (project_id, tech)
                    )
                except mysql.connector.Error as e:
                    print(f'  警告：插入技术栈 {tech} 失败: {e}')
        
        # 插入项目图片（如果有coverImage）
        if project.get('coverImage'):
            try:
                cursor.execute(
                    'INSERT INTO project_images (project_id, image_url, sort_order) VALUES (%s, %s, %s) ON DUPLICATE KEY UPDATE image_url = VALUES(image_url)',
                    (project_id, project['coverImage'], 0)
                )
            except mysql.connector.Error as e:
                print(f'  警告：插入项目图片失败: {e}')
        
        # 插入项目作者（如果有contributors）
        if project.get('contributors'):
            # 获取现有的用户ID列表，如果contributor的id不在列表中，使用第一个用户ID
            for contrib in project['contributors']:
                user_id = contrib['id'] if contrib['id'] in user_ids else user_ids[0]
                try:
                    cursor.execute(
                        'INSERT INTO project_authors (project_id, user_id) VALUES (%s, %s) ON DUPLICATE KEY UPDATE user_id = VALUES(user_id)',
                        (project_id, user_id)
                    )
                except mysql.connector.Error as e:
                    # 如果用户不存在，跳过
                    print(f'  警告：插入项目作者失败 (user_id={user_id}): {e}')
    
    # 导入评论（需要先检查评论表结构）
    print('\n开始导入评论数据...')
    for comment in mock_comments:
        project_id = comment['projectId']
        user_id = comment['userId'] if comment['userId'] in user_ids else user_ids[0]
        
        # 解析创建时间
        created_at = datetime.fromisoformat(comment['createdAt'].replace('Z', '+00:00'))
        
        # 插入评论（comments表结构：comment_id, resource_type, resource_id, parent_id, user_id, content, love_count, reply_total, created_at, updated_at, deleted_at）
        try:
            insert_comment_sql = """
                INSERT INTO comments (
                    resource_type, resource_id, parent_id, user_id, content, 
                    love_count, reply_total, created_at, deleted_at
                ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, NULL)
            """
            cursor.execute(insert_comment_sql, (
                'project',
                project_id,  # resource_id 是整数
                None,  # parent_id 为空（主评论）
                user_id,
                comment['content'],
                comment.get('likes', 0),  # love_count
                0,  # reply_total 初始为0
                created_at
            ))
            comment_id = cursor.lastrowid
            
            print(f'[OK] 已导入评论: {comment["content"][:30]}... (项目ID: {project_id}, 点赞数: {comment.get("likes", 0)})')
        except mysql.connector.Error as e:
            print(f'[ERROR] 导入评论失败: {e}')
            print(f'  评论内容: {comment["content"][:50]}...')
            print(f'  项目ID: {project_id}, 用户ID: {user_id}')
    
    # 提交事务
    conn.commit()
    print('\n[SUCCESS] 所有项目数据导入完成！')
    
    # 验证导入结果
    cursor.execute('SELECT COUNT(*) FROM projects')
    project_count = cursor.fetchone()[0]
    cursor.execute('SELECT COUNT(*) FROM project_tech_stack')
    tech_count = cursor.fetchone()[0]
    cursor.execute('SELECT COUNT(*) FROM comments WHERE resource_type = "project"')
    comment_count = cursor.fetchone()[0]
    
    print(f'\n验证结果:')
    print(f'  项目数: {project_count}')
    print(f'  技术栈记录数: {tech_count}')
    print(f'  评论数: {comment_count}')

except mysql.connector.Error as err:
    print(f'[ERROR] 数据库操作错误: {err}')
    if 'conn' in locals():
        conn.rollback()
finally:
    if 'conn' in locals() and conn.is_connected():
        cursor.close()
        conn.close()
        print('\n数据库连接已关闭')

