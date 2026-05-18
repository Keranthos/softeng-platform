-- 删除不属于专必/专选/公必/公选的课程
-- 注意：此脚本会删除所有不符合分类要求的课程及其所有关联数据

USE softeng;

-- 查看需要删除的课程（先检查，不实际删除）
SELECT 
    c.course_id,
    c.name,
    GROUP_CONCAT(cc.category ORDER BY cc.category SEPARATOR ', ') as categories
FROM courses c
LEFT JOIN course_categories cc ON c.course_id = cc.course_id
GROUP BY c.course_id, c.name
HAVING categories IS NULL 
    OR categories NOT LIKE '%专必%' 
    AND categories NOT LIKE '%专选%' 
    AND categories NOT LIKE '%公必%' 
    AND categories NOT LIKE '%公选%';

-- 开始事务
START TRANSACTION;

-- 1. 删除课程评论（resource_type='course'）
DELETE FROM comments 
WHERE resource_type = 'course' 
AND resource_id IN (
    SELECT c.course_id
    FROM courses c
    LEFT JOIN course_categories cc ON c.course_id = cc.course_id
    GROUP BY c.course_id
    HAVING GROUP_CONCAT(cc.category) IS NULL 
        OR (GROUP_CONCAT(cc.category) NOT LIKE '%专必%' 
            AND GROUP_CONCAT(cc.category) NOT LIKE '%专选%' 
            AND GROUP_CONCAT(cc.category) NOT LIKE '%公必%' 
            AND GROUP_CONCAT(cc.category) NOT LIKE '%公选%')
);

-- 2. 删除课程相关的收藏
DELETE FROM collections 
WHERE resource_type = 'course' 
AND resource_id IN (
    SELECT c.course_id
    FROM courses c
    LEFT JOIN course_categories cc ON c.course_id = cc.course_id
    GROUP BY c.course_id
    HAVING GROUP_CONCAT(cc.category) IS NULL 
        OR (GROUP_CONCAT(cc.category) NOT LIKE '%专必%' 
            AND GROUP_CONCAT(cc.category) NOT LIKE '%专选%' 
            AND GROUP_CONCAT(cc.category) NOT LIKE '%公必%' 
            AND GROUP_CONCAT(cc.category) NOT LIKE '%公选%')
);

-- 3. 删除课程相关的点赞
DELETE FROM likes 
WHERE resource_type = 'course' 
AND resource_id IN (
    SELECT c.course_id
    FROM courses c
    LEFT JOIN course_categories cc ON c.course_id = cc.course_id
    GROUP BY c.course_id
    HAVING GROUP_CONCAT(cc.category) IS NULL 
        OR (GROUP_CONCAT(cc.category) NOT LIKE '%专必%' 
            AND GROUP_CONCAT(cc.category) NOT LIKE '%专选%' 
            AND GROUP_CONCAT(cc.category) NOT LIKE '%公必%' 
            AND GROUP_CONCAT(cc.category) NOT LIKE '%公选%')
);

-- 4. 删除课程（由于外键约束 ON DELETE CASCADE，会自动删除相关的教师、分类、资源、贡献者等）
DELETE FROM courses 
WHERE course_id IN (
    SELECT course_id FROM (
        SELECT c.course_id
        FROM courses c
        LEFT JOIN course_categories cc ON c.course_id = cc.course_id
        GROUP BY c.course_id
        HAVING GROUP_CONCAT(cc.category) IS NULL 
            OR (GROUP_CONCAT(cc.category) NOT LIKE '%专必%' 
                AND GROUP_CONCAT(cc.category) NOT LIKE '%专选%' 
                AND GROUP_CONCAT(cc.category) NOT LIKE '%公必%' 
                AND GROUP_CONCAT(cc.category) NOT LIKE '%公选%')
    ) AS subquery
);

-- 提交事务（如果确认无误，取消注释下面的 COMMIT，并注释掉 ROLLBACK）
-- COMMIT;
ROLLBACK;  -- 默认回滚，避免误删。确认无误后改为 COMMIT;

-- 查询删除后的课程数量
SELECT COUNT(*) as remaining_courses FROM courses;

