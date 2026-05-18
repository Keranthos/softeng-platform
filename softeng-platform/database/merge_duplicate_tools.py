"""
合并重复工具脚本
策略：
1. 对于同名工具，保留数据最丰富的记录（基于评分：评论数×1000 + 浏览量 + 点赞数×10 + 收藏数×5）
2. 将被删除的记录的数据合并到保留的记录中：
   - 浏览量：取最大值
   - 点赞数：求和
   - 收藏数：从collections表重新统计（已经是实时的）
   - 评论：将评论的resource_id更新为保留的记录ID
3. 更新collections表和likes表中的resource_id
4. 删除重复的工具记录及其关联数据
"""

import mysql.connector

conn = mysql.connector.connect(
    host='localhost',
    user='root',
    password='Wan05609',
    database='softeng'
)
cursor = conn.cursor()

print('开始合并重复工具...\n')

# 定义重复工具（工具名称: [要保留的ID, 要删除的ID列表]）
# 基于数据分析，选择数据最丰富的记录作为主记录
duplicate_tools = {
    'ChatGPT': {
        'keep': 121,  # 软件开发分类，有2条评论，数据更丰富
        'remove': [113]  # 个人提升分类
    },
    'Figma': {
        'keep': 114,  # 项目协作分类，浏览量更高
        'remove': [110]  # 软件开发分类
    },
    'GitHub Copilot': {
        'keep': 107,  # 软件开发分类，点赞数更高
        'remove': [122]  # 个人提升分类
    }
}

def get_tool_data(tool_id):
    """获取工具的数据"""
    cursor.execute('''
        SELECT resource_id, category, views, collections, loves, description, description_detail, resource_link
        FROM tools WHERE resource_id = %s
    ''', (tool_id,))
    return cursor.fetchone()

def merge_tool_data(keep_id, remove_ids):
    """合并工具数据"""
    print(f'\n处理工具 ID={keep_id} 和 {remove_ids}:')
    
    # 获取保留记录的数据
    keep_data = get_tool_data(keep_id)
    if not keep_data:
        print(f'  错误：找不到 ID={keep_id} 的工具')
        return False
    
    keep_id_val, keep_category, keep_views, keep_collections, keep_loves, keep_desc, keep_desc_detail, keep_link = keep_data
    print(f'  保留记录: ID={keep_id_val}, 分类={keep_category}, 浏览量={keep_views}, 点赞={keep_loves}')
    
    max_views = keep_views
    total_loves = keep_loves
    
    # 合并被删除记录的数据
    for remove_id in remove_ids:
        remove_data = get_tool_data(remove_id)
        if not remove_data:
            print(f'  警告：找不到 ID={remove_id} 的工具，跳过')
            continue
            
        remove_id_val, remove_category, remove_views, remove_collections, remove_loves, remove_desc, remove_desc_detail, remove_link = remove_data
        print(f'  删除记录: ID={remove_id_val}, 分类={remove_category}, 浏览量={remove_views}, 点赞={remove_loves}')
        
        # 更新最大浏览量
        max_views = max(max_views, remove_views)
        # 累加点赞数
        total_loves += remove_loves
        
        # 更新评论的resource_id
        cursor.execute('''
            UPDATE comments 
            SET resource_id = %s 
            WHERE resource_type = 'tool' AND resource_id = %s
        ''', (keep_id_val, remove_id_val))
        comment_updated = cursor.rowcount
        print(f'    已更新 {comment_updated} 条评论的resource_id')
        
        # 更新collections表的resource_id
        cursor.execute('''
            UPDATE collections 
            SET resource_id = %s 
            WHERE resource_type = 'tool' AND resource_id = %s
        ''', (keep_id_val, remove_id_val))
        collection_updated = cursor.rowcount
        print(f'    已更新 {collection_updated} 条收藏记录的resource_id')
        
        # 更新likes表的resource_id
        cursor.execute('''
            UPDATE likes 
            SET resource_id = %s 
            WHERE resource_type = 'tool' AND resource_id = %s
        ''', (keep_id_val, remove_id_val))
        like_updated = cursor.rowcount
        print(f'    已更新 {like_updated} 条点赞记录的resource_id')
        
        # 更新comment_likes表中的comment_id（如果comment_likes表存在）
        # 注意：comment_likes是通过comment_id关联的，所以需要在更新comments之后处理
        # 但评论的comment_id不会改变，只是resource_id改变了，所以comment_likes不需要更新
        # 只需要确保评论的resource_id已经更新即可
        
        # 删除关联数据（注意：评论、收藏、点赞已经更新resource_id，所以不需要删除）
        # 删除tool_tags
        cursor.execute('DELETE FROM tool_tags WHERE tool_id = %s', (remove_id_val,))
        tag_deleted = cursor.rowcount
        print(f'    已删除 {tag_deleted} 条 tool_tags 记录')
        
        # 删除tool_images
        cursor.execute('DELETE FROM tool_images WHERE tool_id = %s', (remove_id_val,))
        image_deleted = cursor.rowcount
        print(f'    已删除 {image_deleted} 条 tool_images 记录')
        
        # 删除tool_contributors
        cursor.execute('DELETE FROM tool_contributors WHERE tool_id = %s', (remove_id_val,))
        contrib_deleted = cursor.rowcount
        print(f'    已删除 {contrib_deleted} 条 tool_contributors 记录')
        
        # 删除工具记录本身（注意：外键约束会自动删除关联数据，但我们已经手动更新了评论、收藏、点赞）
        cursor.execute('DELETE FROM tools WHERE resource_id = %s', (remove_id_val,))
        print(f'    已删除工具记录 ID={remove_id_val}')
    
    # 更新保留记录的浏览量（取最大值）和点赞数（累加）
    cursor.execute('''
        UPDATE tools 
        SET views = %s, loves = %s
        WHERE resource_id = %s
    ''', (max_views, total_loves, keep_id_val))
    
    # 重新统计收藏数（从collections表）
    cursor.execute('''
        UPDATE tools t
        SET t.collections = (
            SELECT COUNT(*) 
            FROM collections c 
            WHERE c.resource_type = 'tool' AND c.resource_id = t.resource_id
        )
        WHERE t.resource_id = %s
    ''', (keep_id_val,))
    
    print(f'  已更新保留记录的浏览量={max_views}，点赞数={total_loves}')
    print(f'  [OK] 合并完成')
    
    return True

# 执行合并
try:
    for tool_name, config in duplicate_tools.items():
        keep_id = config['keep']
        remove_ids = config['remove']
        merge_tool_data(keep_id, remove_ids)
    
    # 提交事务
    conn.commit()
    print('\n[SUCCESS] 所有重复工具已成功合并！')
    
except Exception as e:
    # 回滚事务
    conn.rollback()
    print(f'\n[ERROR] 合并失败，已回滚: {e}')
    raise

finally:
    conn.close()

