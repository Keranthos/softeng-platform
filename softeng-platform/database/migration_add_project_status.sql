-- 为projects表添加审核状态相关字段
-- 执行此SQL前请先备份数据库

ALTER TABLE projects 
ADD COLUMN status VARCHAR(50) DEFAULT 'pending' COMMENT '审核状态：pending/approved/rejected' AFTER collections,
ADD COLUMN audit_time TIMESTAMP NULL COMMENT '审核时间' AFTER status,
ADD COLUMN reject_reason TEXT NULL COMMENT '拒绝原因' AFTER audit_time,
ADD INDEX idx_status (status);

-- 将现有项目的状态设置为approved（已审核通过）
UPDATE projects SET status = 'approved' WHERE status IS NULL OR status = '';

