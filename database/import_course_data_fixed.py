#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
重新导入课程数据（修复乱码问题）
"""
import pymysql
import sys
import io

if sys.platform == 'win32':
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

# 数据库配置
db_config = {
    'host': 'localhost',
    'user': 'root',
    'password': 'Wan05609',
    'database': 'softeng',
    'charset': 'utf8mb4',
    'cursorclass': pymysql.cursors.DictCursor
}

# 课程数据（从course_data.sql提取，确保使用UTF-8编码）
courses_data = [
    # 大一上学期
    {'id': 101, 'name': '高等数学 I', 'code': 'MATH1001', 'semester': '1-1', 'credit': 5, 'teachers': ['张教授', '李教授'], 'categories': ['数学', '基础课程']},
    {'id': 102, 'name': '程序设计基础', 'code': 'CS1001', 'semester': '1-1', 'credit': 4, 'teachers': ['王教授', '赵教授'], 'categories': ['计算机', '编程']},
    {'id': 103, 'name': '思想道德修养', 'code': 'POLIO1001', 'semester': '1-1', 'credit': 2, 'teachers': ['刘教授'], 'categories': ['政治', '通识']},
    {'id': 104, 'name': '当代文化研究', 'code': 'PUB1001', 'semester': '1-1', 'credit': 2, 'teachers': ['陈教授'], 'categories': ['文化', '通识']},
    
    # 大一下学期
    {'id': 201, 'name': '高等数学 II', 'code': 'MATH1002', 'semester': '1-2', 'credit': 5, 'teachers': ['张教授', '李教授'], 'categories': ['数学', '基础课程']},
    {'id': 202, 'name': '线性代数', 'code': 'MATH1003', 'semester': '1-2', 'credit': 3, 'teachers': ['周教授'], 'categories': ['数学', '基础课程']},
    {'id': 203, 'name': '离散数学', 'code': 'CS1002', 'semester': '1-2', 'credit': 4, 'teachers': ['王教授'], 'categories': ['数学', '计算机']},
    {'id': 204, 'name': '体育2', 'code': 'PE1002', 'semester': '1-2', 'credit': 1, 'teachers': ['体育老师'], 'categories': ['体育', '通识']},
    
    # 大二上学期
    {'id': 301, 'name': '数据结构与算法', 'code': 'CS2001', 'semester': '2-1', 'credit': 5, 'teachers': ['王教授', '赵教授'], 'categories': ['计算机', '算法']},
    {'id': 302, 'name': '计算机组成原理', 'code': 'CS2002', 'semester': '2-1', 'credit': 4, 'teachers': ['孙教授'], 'categories': ['计算机', '硬件']},
    {'id': 303, 'name': 'Python应用开发', 'code': 'CS2005', 'semester': '2-1', 'credit': 2, 'teachers': ['钱教授'], 'categories': ['计算机', '编程']},
    
    # 大二下学期
    {'id': 401, 'name': '操作系统', 'code': 'CS2003', 'semester': '2-2', 'credit': 4, 'teachers': ['孙教授', '吴教授'], 'categories': ['计算机', '系统']},
    {'id': 402, 'name': '计算机网络', 'code': 'CS2004', 'semester': '2-2', 'credit': 4, 'teachers': ['郑教授'], 'categories': ['计算机', '网络']},
    
    # 大三上学期
    {'id': 501, 'name': '计算机网络', 'code': 'CS3001', 'semester': '3-1', 'credit': 3, 'teachers': ['郑教授'], 'categories': ['计算机', '网络']},
    {'id': 502, 'name': '数据库系统概论', 'code': 'CS3002', 'semester': '3-1', 'credit': 3, 'teachers': ['冯教授'], 'categories': ['计算机', '数据库']},
    {'id': 503, 'name': '软件工程', 'code': 'CS3003', 'semester': '3-1', 'credit': 3, 'teachers': ['陈教授'], 'categories': ['计算机', '软件工程']},
    {'id': 504, 'name': 'Web前端开发', 'code': 'CS3004', 'semester': '3-1', 'credit': 2, 'teachers': ['周教授'], 'categories': ['计算机', 'Web开发']},
    
    # 大三下学期
    {'id': 601, 'name': '编译原理', 'code': 'CS3005', 'semester': '3-2', 'credit': 3, 'teachers': ['卫教授'], 'categories': ['计算机', '系统']},
    {'id': 602, 'name': '人工智能基础', 'code': 'CS3006', 'semester': '3-2', 'credit': 3, 'teachers': ['蒋教授'], 'categories': ['计算机', '人工智能']},
    {'id': 603, 'name': '移动应用开发', 'code': 'CS3007', 'semester': '3-2', 'credit': 2, 'teachers': ['沈教授'], 'categories': ['计算机', '移动开发']},
    
    # 大四上学期
    {'id': 701, 'name': '毕业设计', 'code': 'CS4001', 'semester': '4-1', 'credit': 8, 'teachers': ['导师组'], 'categories': ['计算机', '实践']},
    {'id': 702, 'name': '企业实习', 'code': 'CS4002', 'semester': '4-1', 'credit': 4, 'teachers': ['企业导师'], 'categories': ['实践', '实习']},
]

def get_or_create_default_user(cursor):
    """获取或创建默认用户作为提交者"""
    cursor.execute("SELECT id FROM users WHERE username = 'admin' LIMIT 1")
    result = cursor.fetchone()
    if result:
        return result['id']
    
    cursor.execute("""
        INSERT INTO users (username, nickname, email, password, role) 
        VALUES ('admin', '管理员', 'admin@example.com', '$2a$10$default', 'admin')
    """)
    cursor.execute("SELECT LAST_INSERT_ID() as id")
    result = cursor.fetchone()
    return result['id']

def import_courses():
    """导入课程数据"""
    connection = None
    try:
        connection = pymysql.connect(**db_config)
        cursor = connection.cursor()
        
        print("=" * 60)
        print("开始导入课程数据...")
        print("=" * 60)
        
        # 获取或创建默认用户
        submitter_id = get_or_create_default_user(cursor)
        print(f"使用提交者ID: {submitter_id}")
        
        # 清空现有课程数据
        print("\n清理现有课程数据...")
        cursor.execute("SET FOREIGN_KEY_CHECKS = 0")
        cursor.execute("DELETE FROM course_resources_web")
        cursor.execute("DELETE FROM course_resources_upload")
        cursor.execute("DELETE FROM course_contributors")
        cursor.execute("DELETE FROM course_categories")
        cursor.execute("DELETE FROM course_teachers")
        cursor.execute("DELETE FROM courses")
        cursor.execute("SET FOREIGN_KEY_CHECKS = 1")
        print("清理完成")
        
        success_count = 0
        error_count = 0
        
        # 导入每个课程
        for course in courses_data:
            try:
                # 插入课程主表
                sql_course = """
                    INSERT INTO courses (
                        course_id, resource_type, name, semester, credit, 
                        cover, views, loves, collections
                    ) VALUES (
                        %s, 'course', %s, %s, %s, 
                        %s, 0, 0, 0
                    )
                """
                
                cover_url = f"https://placehold.co/600x450/3b82f6/ffffff?text={course['code']}"
                
                cursor.execute(sql_course, (
                    course['id'],
                    course['name'],
                    course['semester'],
                    course['credit'],
                    cover_url
                ))
                
                course_id = course['id']
                
                # 插入教师
                if course.get('teachers'):
                    sql_teacher = """
                        INSERT INTO course_teachers (course_id, teacher_name)
                        VALUES (%s, %s)
                    """
                    for teacher in course['teachers']:
                        cursor.execute(sql_teacher, (course_id, teacher))
                
                # 插入分类
                if course.get('categories'):
                    sql_category = """
                        INSERT INTO course_categories (course_id, category)
                        VALUES (%s, %s)
                    """
                    for category in course['categories']:
                        cursor.execute(sql_category, (course_id, category))
                
                success_count += 1
                print(f"[成功] 导入成功: {course['name']} (ID: {course_id})")
                
            except Exception as e:
                error_count += 1
                print(f"[失败] 导入失败: {course['name']} - {str(e)}")
                continue
        
        # 提交事务
        connection.commit()
        
        print("\n" + "=" * 60)
        print(f"导入完成！成功: {success_count}, 失败: {error_count}")
        print("=" * 60)
        
    except Exception as e:
        if connection:
            connection.rollback()
        print(f"\n数据库错误: {e}")
        import traceback
        traceback.print_exc()
    finally:
        if connection:
            connection.close()
            print("\n数据库连接已关闭")

if __name__ == '__main__':
    import_courses()

