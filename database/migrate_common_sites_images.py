#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
迁移"常用"部分的图标到本地
下载主页"常用"部分的10个网站图标并保存到本地
"""

import os
import sys
import requests
import hashlib
from pathlib import Path
from datetime import datetime
import urllib.parse
import urllib3

# 禁用SSL警告
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

# 配置输出编码（Windows兼容）
try:
    sys.stdout.reconfigure(encoding='utf-8')
except AttributeError:
    # Python < 3.7
    import codecs
    sys.stdout = codecs.getwriter('utf-8')(sys.stdout.buffer, 'strict')

# 脚本目录
SCRIPT_DIR = Path(__file__).parent.absolute()
PROJECT_ROOT = SCRIPT_DIR.parent  # database/ 的父目录
UPLOAD_DIR = PROJECT_ROOT / 'uploads' / 'images' / 'common'

MAX_IMAGE_SIZE = 5 * 1024 * 1024  # 5MB
TIMEOUT = 30  # 下载超时时间（秒）
MAX_RETRIES = 3  # 最大重试次数
RETRY_DELAY = 2  # 重试延迟（秒）

# 常用网站数据（从 store/index.js 中提取）
COMMON_SITES = [
    { 'name': '微信文件', 'url': 'https://file.fengfengzhidao.com', 'icon': 'https://file.fengfengzhidao.com/logo/wechat.png' },
    { 'name': '和风天气', 'url': 'https://www.qweather.com', 'icon': 'https://cdn.heweather.com/img/logo.png' },
    { 'name': '小红书', 'url': 'https://www.xiaohongshu.com', 'icon': 'https://ci.xiaohongshu.com/logo_2023.png' },
    { 'name': '哔哩哔哩', 'url': 'https://www.bilibili.com', 'icon': 'https://www.bilibili.com/favicon.ico' },
    { 'name': '知乎', 'url': 'https://www.zhihu.com', 'icon': 'https://static.zhihu.com/static/favicon.ico' },
    { 'name': '百度翻译', 'url': 'https://fanyi.baidu.com', 'icon': 'https://fanyi.bdstatic.com/static/translation/img/favicon.ico' },
    { 'name': '淘宝', 'url': 'https://www.taobao.com', 'icon': 'https://www.taobao.com/favicon.ico' },
    { 'name': '抖音', 'url': 'https://www.douyin.com', 'icon': 'https://lf1-cdn2-tos.bytego.com/obj/ies-fe-bee-prod/cn/fe/bee_prod_cn_bee_home_page_logo.png' },
    { 'name': '京东', 'url': 'https://www.jd.com', 'icon': 'https://www.jd.com/favicon.ico' },
    { 'name': '微博', 'url': 'https://www.weibo.com', 'icon': 'https://weibo.com/favicon.ico' }
]

def ensure_upload_dir():
    """确保上传目录存在"""
    UPLOAD_DIR.mkdir(parents=True, exist_ok=True)
    return UPLOAD_DIR

def generate_filename(original_url, site_name):
    """生成文件名"""
    # 使用URL的MD5哈希 + 网站名称作为文件名
    url_hash = hashlib.md5(original_url.encode('utf-8')).hexdigest()[:12]
    # 从URL中提取扩展名
    parsed = urllib.parse.urlparse(original_url)
    path = parsed.path
    ext = os.path.splitext(path)[1] or '.png'  # 默认使用.png
    # 清理扩展名
    ext = ext.split('?')[0]  # 移除查询参数
    if ext not in ['.png', '.jpg', '.jpeg', '.gif', '.svg', '.ico', '.webp']:
        ext = '.png'
    
    # 使用网站名称（清理特殊字符）作为文件名的一部分
    safe_name = ''.join(c for c in site_name if c.isalnum() or c in ('-', '_'))[:20]
    return f"{safe_name}_{url_hash}{ext}"

def download_and_save_image(url, site_name):
    """下载图片并保存到本地"""
    if not url or not url.startswith('http'):
        print(f"  ⚠️  跳过无效URL: {url}")
        return None
    
    for attempt in range(MAX_RETRIES):
        try:
            headers = {
                'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
            }
            
            response = requests.get(url, headers=headers, timeout=TIMEOUT, verify=False, stream=True)
            response.raise_for_status()
            
            # 检查Content-Type
            content_type = response.headers.get('Content-Type', '')
            if not content_type.startswith('image/'):
                print(f"  ⚠️  警告: {url} 不是图片类型 ({content_type})")
                # 继续尝试保存，可能是favicon.ico等
            
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
            filename = generate_filename(url, site_name)
            file_path = upload_path / filename
            
            # 保存文件
            with open(file_path, 'wb') as f:
                f.write(image_data)
            
            # 返回相对路径（用于URL）
            relative_path = file_path.relative_to(PROJECT_ROOT)
            # 将路径分隔符统一为 /
            local_url = '/' + str(relative_path).replace('\\', '/')
            
            print(f"  ✅ 下载成功: {filename}")
            return local_url
            
        except requests.exceptions.RequestException as e:
            if attempt < MAX_RETRIES - 1:
                print(f"  ⚠️  下载失败 (尝试 {attempt + 1}/{MAX_RETRIES}): {e}")
                import time
                time.sleep(RETRY_DELAY)
            else:
                print(f"  ❌ 下载失败 (已重试 {MAX_RETRIES} 次): {e}")
                return None
        except Exception as e:
            print(f"  ❌ 发生错误: {e}")
            return None
    
    return None

def main():
    """主函数"""
    print("=" * 60)
    print("开始迁移'常用'部分的图标到本地...")
    print("=" * 60)
    
    # 确保目录存在
    ensure_upload_dir()
    
    # 统计信息
    success_count = 0
    fail_count = 0
    updated_sites = []
    
    # 处理每个网站
    for site in COMMON_SITES:
        print(f"\n处理: {site['name']}")
        print(f"  URL: {site['icon']}")
        
        local_url = download_and_save_image(site['icon'], site['name'])
        
        if local_url:
            success_count += 1
            updated_sites.append({
                'name': site['name'],
                'url': site['url'],
                'icon': local_url,
                'desc': site.get('desc', '')
            })
        else:
            fail_count += 1
            # 如果下载失败，保留原始URL
            updated_sites.append({
                'name': site['name'],
                'url': site['url'],
                'icon': site['icon'],
                'desc': site.get('desc', '')
            })
    
    # 输出结果
    print("\n" + "=" * 60)
    print("迁移完成！")
    print("=" * 60)
    print(f"成功: {success_count}/{len(COMMON_SITES)}")
    print(f"失败: {fail_count}/{len(COMMON_SITES)}")
    
    # 输出更新后的数据（用于更新 store/index.js）
    print("\n" + "=" * 60)
    print("更新后的 commonSites 数据（复制到 store/index.js）:")
    print("=" * 60)
    print("commonSites: [")
    for i, site in enumerate(updated_sites):
        desc = site.get('desc', '')
        if not desc:
            # 根据名称生成描述
            desc_map = {
                '微信文件': '快速传输文件到设备',
                '和风天气': '实时天气预报服务',
                '小红书': '生活方式分享社区',
                '哔哩哔哩': '视频弹幕网站',
                '知乎': '高质量问答平台',
                '百度翻译': '多语言翻译工具',
                '淘宝': '在线购物平台',
                '抖音': '短视频分享应用',
                '京东': '电商购物网站',
                '微博': '社交媒体平台'
            }
            desc = desc_map.get(site['name'], '')
        
        comma = ',' if i < len(updated_sites) - 1 else ''
        print(f"  {{ name: '{site['name']}', url: '{site['url']}', icon: '{site['icon']}', desc: '{desc}' }}{comma}")
    print("]")
    
    print("\n" + "=" * 60)
    print("提示: 请手动更新 softeng-platform-frontend/src/store/index.js 中的 commonSites 数据")
    print("=" * 60)

if __name__ == '__main__':
    main()

