#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import pymysql
import os

# 数据库配置
db_config = {
    'host': 'localhost',
    'user': 'root',
    'password': 'Wan05609',
    'database': 'softeng',
    'charset': 'utf8mb4'
}

# 读取SQL文件
sql_file = os.path.join(os.path.dirname(__file__), 'course_data.sql')
with open(sql_file, 'r', encoding='utf-8') as f:
    sql_content = f.read()

# 连接数据库并执行SQL
try:
    connection = pymysql.connect(**db_config)
    cursor = connection.cursor()
    
    # 分割SQL语句（以分号和换行分割）
    statements = [s.strip() for s in sql_content.split(';') if s.strip() and not s.strip().startswith('--')]
    
    for statement in statements:
        if statement:
            try:
                cursor.execute(statement)
                print(f"执行成功: {statement[:50]}...")
            except Exception as e:
                if "Duplicate column name" not in str(e) and "already exists" not in str(e):
                    print(f"执行错误: {e}")
                    print(f"语句: {statement[:100]}")
    
    connection.commit()
    print("数据导入完成！")
    
except Exception as e:
    print(f"数据库错误: {e}")
finally:
    if 'connection' in locals():
        connection.close()

