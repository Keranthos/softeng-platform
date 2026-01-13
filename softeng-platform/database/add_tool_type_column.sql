-- 添加工具类型字段
ALTER TABLE tools ADD COLUMN tool_type VARCHAR(20) DEFAULT 'external' COMMENT '工具类型：internal/external';

