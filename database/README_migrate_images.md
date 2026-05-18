# 图片本地化迁移脚本使用说明

## 功能说明

此脚本用于将数据库中已有的外部图片URL下载到本地文件系统，并更新数据库中的URL为本地路径。

## 需要处理的表

1. **tool_images** - 工具图片表 (`image_url` 字段)
2. **project_images** - 项目图片表 (`image_url` 字段)
3. **courses** - 课程表 (`cover` 字段)
4. **projects** - 项目表 (`cover` 字段)

## 使用方法

### 1. 安装依赖

```bash
pip install mysql-connector-python requests
# 或
pip install pymysql requests
```

### 2. 配置环境变量（可选）

```bash
# Windows PowerShell
$env:DB_HOST="localhost"
$env:DB_PORT="3306"
$env:DB_USER="root"
$env:DB_PASSWORD="your_password"
$env:DB_NAME="softeng_platform"

# Linux/Mac
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=softeng_platform
```

### 3. 直接修改脚本中的数据库配置

编辑 `migrate_images_to_local.py`，修改 `DB_CONFIG` 字典：

```python
DB_CONFIG = {
    'host': 'localhost',
    'port': 3306,
    'user': 'root',
    'password': 'your_password',
    'database': 'softeng_platform',
    'charset': 'utf8mb4'
}
```

### 4. 运行脚本

```bash
# 在项目根目录下运行
cd softeng-platform/softeng-platform/database
python migrate_images_to_local.py
```

## 脚本功能

- ✅ 自动识别外部URL（http:// 或 https://）
- ✅ 跳过已本地化的路径（/uploads/ 开头）
- ✅ 下载图片到 `uploads/images/YYYY/MM/` 目录
- ✅ 生成唯一文件名（基于时间戳和URL的MD5）
- ✅ 更新数据库中的URL为本地路径
- ✅ 错误处理和统计报告

## 注意事项

1. **备份数据库**：运行前请先备份数据库
2. **网络连接**：需要能够访问外部图片URL
3. **磁盘空间**：确保有足够的磁盘空间存储图片
4. **执行时间**：根据图片数量，可能需要较长时间
5. **失败处理**：下载失败的图片会保留原URL，不会更新数据库

## 输出示例

```
============================================================
🖼️  图片本地化迁移脚本
============================================================
✅ 数据库连接成功

📋 处理表: tool_images
  📊 找到 50 条记录
  🔄 处理 ID=1: https://example.com/image1.jpg
  ✅ 下载成功: /uploads/images/2024/01/abc123def456.jpg
  ...

📊 迁移完成统计
============================================================
✅ 成功本地化: 45 张图片
⏭️  跳过/失败: 5 张图片
📁 图片保存在: uploads/images/
============================================================
```

