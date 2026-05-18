import mysql.connector

conn = mysql.connector.connect(
    host='localhost',
    user='root',
    password='Wan05609',
    database='softeng'
)
cursor = conn.cursor()

print('验证合并结果：\n')

# 检查是否还有重复工具
cursor.execute('SELECT resource_name, COUNT(*) as cnt FROM tools GROUP BY resource_name HAVING cnt > 1')
duplicates = cursor.fetchall()

if duplicates:
    print('[WARNING] 仍然存在重复工具：')
    for name, cnt in duplicates:
        print(f'  {name}: {cnt} 个记录')
else:
    print('[OK] 没有发现重复工具')

# 检查合并后的工具数据
print('\n合并后的工具数据：')
cursor.execute('SELECT resource_id, resource_name, category, views, loves FROM tools WHERE resource_name IN (\'ChatGPT\', \'Figma\', \'GitHub Copilot\') ORDER BY resource_name')
tools = cursor.fetchall()
for tool_id, name, cat, views, loves in tools:
    cursor.execute('SELECT COUNT(*) FROM comments WHERE resource_type=\'tool\' AND resource_id=%s AND deleted_at IS NULL', (tool_id,))
    comment_count = cursor.fetchone()[0]
    cursor.execute('SELECT COUNT(*) FROM collections WHERE resource_type=\'tool\' AND resource_id=%s', (tool_id,))
    collections = cursor.fetchone()[0]
    print(f'  ID={tool_id}, 名称={name}, 分类={cat}, 浏览量={views}, 点赞={loves}, 评论={comment_count}, 收藏={collections}')

# 检查工具总数和分类分布
print('\n工具总数和分类分布：')
cursor.execute('SELECT COUNT(*) FROM tools')
total = cursor.fetchone()[0]
print(f'总工具数: {total}')

cursor.execute('SELECT category, COUNT(*) FROM tools GROUP BY category')
print('分类分布：')
for cat, cnt in cursor.fetchall():
    print(f'  {cat}: {cnt}')

conn.close()

