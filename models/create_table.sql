use bluebell;

## 用户表
CREATE TABLE `communities`
(
    `id`             bigint(20) unsigned                     NOT NULL AUTO_INCREMENT,
    `community_id`   bigint(20)                              NOT NULL COMMENT '社区ID',
    `community_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '社区名称',
    `introduction`   varchar(256) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '详情',
    `created_at`     datetime(3)                                      DEFAULT NULL,
    `updated_at`     datetime(3)                                      DEFAULT NULL,
    `deleted_at`     datetime(3)                                      DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_community_id` (`community_id`),
    UNIQUE KEY `idx_community_name` (`community_name`),
    KEY `idx_communities_deleted_at` (`deleted_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;

## 社区表
CREATE TABLE `user`
(
    `id`          bigint(20)                             NOT NULL AUTO_INCREMENT,
    `user_id`     bigint(20)                             NOT NULL,
    `username`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
    `password`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
    `email`       varchar(64) COLLATE utf8mb4_general_ci,
    `gender`      tinyint(4)                             NOT NULL DEFAULT '0',
    `create_time` timestamp                              NULL     DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp                              NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`) USING BTREE,
    UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;

INSERT INTO `bluebell`.`communities` (`id`, `created_at`, `updated_at`, `deleted_at`, `community_id`, `community_name`, `introduction`) VALUES ('1', '2020-10-03 17:08:06.000', '2020-10-03 17:08:10.000', NULL, '1', 'Go', 'Golang');
INSERT INTO `bluebell`.`communities` (`id`, `created_at`, `updated_at`, `deleted_at`, `community_id`, `community_name`, `introduction`) VALUES ('2', '2020-10-03 17:08:06.000', '2020-10-03 17:08:10.000', NULL, '2', 'LeetCode', '刷题刷题~冲鸭~');
INSERT INTO `bluebell`.`communities` (`id`, `created_at`, `updated_at`, `deleted_at`, `community_id`, `community_name`, `introduction`) VALUES ('3', '2020-10-03 17:08:06.000', '2020-10-03 17:08:10.000', NULL, '3', 'LOL', '各种LOL相关资讯');
INSERT INTO `bluebell`.`communities` (`id`, `created_at`, `updated_at`, `deleted_at`, `community_id`, `community_name`, `introduction`) VALUES ('4', '2020-10-03 17:08:06.000', '2020-10-03 17:08:10.000', NULL, '4', 'PUBG', '没有什么是一把98K解决不了的');
