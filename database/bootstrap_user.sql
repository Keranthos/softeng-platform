-- 开发环境：创建与 internal/config 默认一致的应用账号（与 DB_USER/DB_PASSWORD 默认 softeng_app / 123456 对齐）
-- 使用方式：用 root 或有权限的账号执行一次（在导入 schema.sql 前后均可，但须保证 softeng 库已存在）
--   mysql -u root -p < database/bootstrap_user.sql

CREATE DATABASE IF NOT EXISTS softeng CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS 'softeng_app'@'localhost' IDENTIFIED BY '123456';
GRANT ALL PRIVILEGES ON softeng.* TO 'softeng_app'@'localhost';

CREATE USER IF NOT EXISTS 'softeng_app'@'127.0.0.1' IDENTIFIED BY '123456';
GRANT ALL PRIVILEGES ON softeng.* TO 'softeng_app'@'127.0.0.1';

FLUSH PRIVILEGES;
