-- 给 ips 表添加 user_id 列
ALTER TABLE `ips` ADD COLUMN `user_id` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id' AFTER `article_id`;
