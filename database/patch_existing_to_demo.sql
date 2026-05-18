-- 从旧版 schema 升级到「演示/详情页」所需字段（在已存在的 softeng 库上按需执行）。
-- 若某条报错 Duplicate column name，说明该列已存在，请跳过该条继续。
-- 若已用当前仓库中的 schema.sql 全新建库，通常无需执行本文件。
USE softeng;

ALTER TABLE courses ADD COLUMN description TEXT NULL COMMENT '课程简介（详情页）' AFTER credit;

ALTER TABLE projects
  ADD COLUMN status VARCHAR(50) DEFAULT 'pending' COMMENT '审核状态：pending/approved/rejected' AFTER collections,
  ADD COLUMN audit_time TIMESTAMP NULL COMMENT '审核时间' AFTER status,
  ADD COLUMN reject_reason TEXT NULL COMMENT '拒绝原因' AFTER audit_time,
  ADD INDEX idx_status (status);

UPDATE projects SET status = 'approved' WHERE status IS NULL OR status = '';

ALTER TABLE resource_status_logs ADD COLUMN reject_reason TEXT NULL COMMENT '拒绝原因' AFTER operate_time;
