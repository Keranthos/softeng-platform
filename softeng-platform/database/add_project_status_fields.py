import mysql.connector

conn = mysql.connector.connect(
    host='localhost',
    user='root',
    password='Wan05609',
    database='softeng'
)
cursor = conn.cursor()

try:
    # 检查 status 字段是否存在
    cursor.execute("SHOW COLUMNS FROM projects LIKE 'status'")
    status_exists = cursor.fetchone()
    
    if not status_exists:
        print('添加 status 字段...')
        cursor.execute("""
            ALTER TABLE projects 
            ADD COLUMN status VARCHAR(50) DEFAULT 'approved' COMMENT '审核状态：pending/approved/rejected' AFTER collections
        """)
        conn.commit()
        print('[OK] 已添加 status 字段')
    else:
        print('[SKIP] status 字段已存在')
    
    # 检查 audit_time 字段是否存在
    cursor.execute("SHOW COLUMNS FROM projects LIKE 'audit_time'")
    audit_time_exists = cursor.fetchone()
    
    if not audit_time_exists:
        print('添加 audit_time 字段...')
        cursor.execute("""
            ALTER TABLE projects 
            ADD COLUMN audit_time TIMESTAMP NULL COMMENT '审核时间' AFTER status
        """)
        conn.commit()
        print('[OK] 已添加 audit_time 字段')
    else:
        print('[SKIP] audit_time 字段已存在')
    
    # 检查 reject_reason 字段是否存在
    cursor.execute("SHOW COLUMNS FROM projects LIKE 'reject_reason'")
    reject_reason_exists = cursor.fetchone()
    
    if not reject_reason_exists:
        print('添加 reject_reason 字段...')
        cursor.execute("""
            ALTER TABLE projects 
            ADD COLUMN reject_reason TEXT NULL COMMENT '拒绝原因' AFTER audit_time
        """)
        conn.commit()
        print('[OK] 已添加 reject_reason 字段')
    else:
        print('[SKIP] reject_reason 字段已存在')
    
    # 检查并添加索引
    cursor.execute("SHOW INDEX FROM projects WHERE Key_name = 'idx_status'")
    index_exists = cursor.fetchone()
    
    if not index_exists:
        print('添加 status 索引...')
        cursor.execute("CREATE INDEX idx_status ON projects(status)")
        conn.commit()
        print('[OK] 已添加 idx_status 索引')
    else:
        print('[SKIP] idx_status 索引已存在')
    
    # 将现有项目的状态设置为 approved
    print('更新现有项目的状态...')
    cursor.execute("UPDATE projects SET status = 'approved' WHERE status IS NULL OR status = ''")
    updated_rows = cursor.rowcount
    conn.commit()
    print(f'[OK] 已将 {updated_rows} 个项目的状态设置为 approved')
    
    # 验证结果
    cursor.execute("SELECT COUNT(*) FROM projects WHERE status = 'approved'")
    approved_count = cursor.fetchone()[0]
    cursor.execute("SELECT COUNT(*) FROM projects")
    total_count = cursor.fetchone()[0]
    print(f'\n验证结果:')
    print(f'  总项目数: {total_count}')
    print(f'  已审核通过: {approved_count}')
    
except mysql.connector.Error as err:
    print(f'[ERROR] 数据库操作错误: {err}')
    conn.rollback()
finally:
    if conn.is_connected():
        cursor.close()
        conn.close()

