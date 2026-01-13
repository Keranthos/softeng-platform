#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
图片本地化迁移脚本
将数据库中已有的外部图片URL下载到本地，并更新数据库中的URL
"""

import sys
import os

# 确保输出使用UTF-8编码（Windows兼容）
if sys.stdout.encoding != 'utf-8':
    try:
        sys.stdout.reconfigure(encoding='utf-8')
    except AttributeError:
        # Python < 3.7 兼容
        import codecs
        sys.stdout = codecs.getwriter('utf-8')(sys.stdout.buffer, 'strict')
        sys.stderr = codecs.getwriter('utf-8')(sys.stderr.buffer, 'strict')
import requests
import hashlib
import time
from datetime import datetime
from pathlib import Path
from urllib.parse import urlparse
import urllib3

# 禁用SSL警告（某些网站SSL证书可能有问题）
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

try:
    import mysql.connector
    USE_MYSQL_CONNECTOR = True
except ImportError:
    try:
        import pymysql
        USE_MYSQL_CONNECTOR = False
    except ImportError:
        print("❌ 错误: 需要安装 mysql-connector-python 或 pymysql")
        print("   安装命令: pip install mysql-connector-python")
        print("   或: pip install pymysql")
        sys.exit(1)

# 数据库配置（从环境变量或配置文件读取）
DB_CONFIG = {
    'host': os.getenv('DB_HOST', '127.0.0.1'),
    'port': int(os.getenv('DB_PORT', 3306)),
    'user': os.getenv('DB_USER', 'root'),
    'password': os.getenv('DB_PASSWORD', 'Wan05609'),
    'database': os.getenv('DB_NAME', 'softeng'),
    'charset': 'utf8mb4'
}

# 上传目录配置（相对于项目根目录，即 softeng-platform/softeng-platform/）
# 注意：迁移脚本在 database/ 目录下运行，但图片应该保存在项目根目录的 uploads/ 下
import os
SCRIPT_DIR = Path(__file__).parent.absolute()
PROJECT_ROOT = SCRIPT_DIR.parent  # database/ 的父目录
UPLOAD_DIR = PROJECT_ROOT / 'uploads' / 'images'
MAX_IMAGE_SIZE = 5 * 1024 * 1024  # 5MB
TIMEOUT = 30  # 下载超时时间（秒）
MAX_RETRIES = 3  # 最大重试次数
RETRY_DELAY = 2  # 重试延迟（秒）

def ensure_upload_dir():
    """确保上传目录存在"""
    now = datetime.now()
    year = now.strftime('%Y')
    month = now.strftime('%m')
    upload_path = Path(UPLOAD_DIR) / year / month
    upload_path.mkdir(parents=True, exist_ok=True)
    return upload_path

def generate_filename(original_url):
    """生成文件名（基于时间戳和URL的MD5）"""
    timestamp = int(time.time() * 1000000)  # 微秒时间戳
    hash_obj = hashlib.md5(f"{timestamp}_{original_url}".encode())
    hash_str = hash_obj.hexdigest()[:16]
    
    # 从URL获取扩展名
    parsed = urlparse(original_url)
    ext = os.path.splitext(parsed.path)[1]
    if not ext or ext not in ['.jpg', '.jpeg', '.png', '.gif', '.webp']:
        ext = '.jpg'  # 默认扩展名
    
    return f"{hash_str}{ext}"

def is_external_url(url):
    """判断是否为外部URL"""
    if not url:
        return False
    # 检查是否为HTTP/HTTPS链接
    return url.startswith('http://') or url.startswith('https://')

def is_local_path(url):
    """判断是否为本地路径（已本地化）"""
    if not url:
        return False
    # 本地路径通常以 /uploads/ 开头
    return url.startswith('/uploads/') or url.startswith('uploads/')

def download_and_save_image(url):
    """下载外部图片并保存到本地（带重试机制）"""
    for attempt in range(MAX_RETRIES):
        try:
            # 创建HTTP会话
            session = requests.Session()
            session.headers.update({
                'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
                'Accept': 'image/webp,image/apng,image/*,*/*;q=0.8',
                'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8',
                'Referer': url.split('/')[0] + '//' + url.split('/')[2] if '/' in url else url
            })
            
            # 下载图片（禁用SSL验证以避免某些SSL错误）
            response = session.get(url, timeout=TIMEOUT, stream=True, verify=False, allow_redirects=True)
            response.raise_for_status()
            
            # 检查Content-Type
            content_type = response.headers.get('Content-Type', '')
            if not content_type.startswith('image/'):
                print(f"  ⚠️  警告: {url} 不是图片类型 ({content_type})")
                return None
            
            # 检查文件大小
            content_length = response.headers.get('Content-Length')
            if content_length and int(content_length) > MAX_IMAGE_SIZE:
                print(f"  ⚠️  警告: {url} 文件太大 ({content_length} bytes)")
                return None
            
            # 读取图片数据
            image_data = b''
            for chunk in response.iter_content(chunk_size=8192):
                image_data += chunk
                if len(image_data) > MAX_IMAGE_SIZE:
                    print(f"  ⚠️  警告: {url} 文件太大")
                    return None
            
            # 确保目录存在
            upload_path = ensure_upload_dir()
            
            # 生成文件名
            filename = generate_filename(url)
            file_path = upload_path / filename
            
            # 保存文件
            with open(file_path, 'wb') as f:
                f.write(image_data)
            
            # 返回相对路径（用于URL）
            relative_path = file_path.relative_to(Path('.'))
            # 将路径分隔符统一为 /
            local_url = '/' + str(relative_path).replace('\\', '/')
            
            print(f"  ✅ 下载成功: {local_url}")
            return local_url
            
        except requests.exceptions.RequestException as e:
            if attempt < MAX_RETRIES - 1:
                print(f"  ⚠️  下载失败 (尝试 {attempt + 1}/{MAX_RETRIES}): {str(e)[:100]}... 等待 {RETRY_DELAY} 秒后重试")
                time.sleep(RETRY_DELAY)
            else:
                print(f"  ❌ 下载失败 (已重试 {MAX_RETRIES} 次): {str(e)[:100]}")
                return None
        except Exception as e:
            if attempt < MAX_RETRIES - 1:
                print(f"  ⚠️  处理失败 (尝试 {attempt + 1}/{MAX_RETRIES}): {str(e)[:100]}... 等待 {RETRY_DELAY} 秒后重试")
                time.sleep(RETRY_DELAY)
            else:
                print(f"  ❌ 保存失败 (已重试 {MAX_RETRIES} 次): {str(e)[:100]}")
                return None
    
    return None

def migrate_table_images(cursor, table_name, id_field, url_field, where_clause=""):
    """迁移指定表的图片"""
    print(f"\n📋 处理表: {table_name}")
    
    # 查询所有外部URL
    query = f"SELECT {id_field}, {url_field} FROM {table_name} WHERE {url_field} IS NOT NULL AND {url_field} != ''"
    if where_clause:
        query += f" AND {where_clause}"
    
    cursor.execute(query)
    rows = cursor.fetchall()
    
    if not rows:
        print(f"  ℹ️  没有需要处理的记录")
        return 0, 0
    
    print(f"  📊 找到 {len(rows)} 条记录")
    
    success_count = 0
    skip_count = 0
    
    for row in rows:
        record_id = row[0]
        original_url = row[1]
        
        # 跳过已本地化的路径
        if is_local_path(original_url):
            skip_count += 1
            continue
        
        # 只处理外部URL
        if not is_external_url(original_url):
            skip_count += 1
            continue
        
        print(f"  🔄 处理 ID={record_id}: {original_url}")
        
        # 下载并保存图片
        local_url = download_and_save_image(original_url)
        
        if local_url:
            # 更新数据库
            try:
                update_query = f"UPDATE {table_name} SET {url_field} = %s WHERE {id_field} = %s"
                cursor.execute(update_query, (local_url, record_id))
                success_count += 1
            except Exception as e:
                print(f"  ❌ 更新数据库失败: {e}")
        else:
            skip_count += 1
            print(f"  ⏭️  跳过（下载失败）")
    
    return success_count, skip_count

def migrate_tool_images(cursor):
    """迁移工具图片表"""
    return migrate_table_images(
        cursor,
        'tool_images',
        'id',
        'image_url'
    )

def migrate_project_images(cursor):
    """迁移项目图片表"""
    return migrate_table_images(
        cursor,
        'project_images',
        'id',
        'image_url'
    )

def migrate_course_covers(cursor):
    """迁移课程封面"""
    return migrate_table_images(
        cursor,
        'courses',
        'course_id',
        'cover'
    )

def migrate_project_covers(cursor):
    """迁移项目封面"""
    return migrate_table_images(
        cursor,
        'projects',
        'project_id',
        'cover'
    )

def main():
    """主函数"""
    print("=" * 60)
    print("🖼️  图片本地化迁移脚本")
    print("=" * 60)
    
    # 连接数据库
    try:
        if USE_MYSQL_CONNECTOR:
            connection = mysql.connector.connect(**DB_CONFIG)
            cursor = connection.cursor(dictionary=False)
        else:
            connection = pymysql.connect(**DB_CONFIG)
            cursor = connection.cursor()
        print("✅ 数据库连接成功")
    except Exception as e:
        print(f"❌ 数据库连接失败: {e}")
        sys.exit(1)
    
    try:
        total_success = 0
        total_skip = 0
        
        # 迁移工具图片
        success, skip = migrate_tool_images(cursor)
        total_success += success
        total_skip += skip
        
        # 迁移项目图片
        success, skip = migrate_project_images(cursor)
        total_success += success
        total_skip += skip
        
        # 迁移课程封面
        success, skip = migrate_course_covers(cursor)
        total_success += success
        total_skip += skip
        
        # 迁移项目封面
        success, skip = migrate_project_covers(cursor)
        total_success += success
        total_skip += skip
        
        # 提交事务
        connection.commit()
        
        print("\n" + "=" * 60)
        print("📊 迁移完成统计")
        print("=" * 60)
        print(f"✅ 成功本地化: {total_success} 张图片")
        print(f"⏭️  跳过/失败: {total_skip} 张图片")
        print(f"📁 图片保存在: {UPLOAD_DIR}/")
        print("=" * 60)
        
    except Exception as e:
        print(f"\n❌ 迁移过程中出错: {e}")
        connection.rollback()
        sys.exit(1)
    finally:
        cursor.close()
        connection.close()
        print("\n✅ 数据库连接已关闭")

if __name__ == '__main__':
    main()

