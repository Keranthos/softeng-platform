-- 同步 tools 表的 collections 字段与 collections 表的实际记录数
-- 这个脚本会确保 tools.collections 字段的值与 collections 表的实际记录数一致

UPDATE tools t
SET t.collections = (
    SELECT COUNT(*)
    FROM collections c
    WHERE c.resource_type = 'tool' AND c.resource_id = t.resource_id
);

-- 同样同步 courses 表的 collections 字段
UPDATE courses co
SET co.collections = (
    SELECT COUNT(*)
    FROM collections c
    WHERE c.resource_type = 'course' AND c.resource_id = co.course_id
);

-- 同样同步 projects 表的 collections 字段
UPDATE projects p
SET p.collections = (
    SELECT COUNT(*)
    FROM collections c
    WHERE c.resource_type = 'project' AND c.resource_id = p.project_id
);

