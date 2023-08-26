create table `user` (
    `id` bigint(20) not null auto_increment,
    -- 另外设置user_id的原因：
    -- 1. 防止用户注册时通过id知道当前有多少用户
    -- 2. 当分库和分表时（用户量特别大时），不同库中的user_id可能重复
    `user_id` bigint(20) not null,
    `username` varchar(64) collate utf8mb4_general_ci not null,
    `password` varchar(64) collate utf8mb4_general_ci not null,
    `email` varchar(64) collate utf8mb4_general_ci,
    `gender` tinyint(4) not null default '0',
    `create_time` timestamp null default current_timestamp,
    `update_time` timestamp null default current_timestamp on update
                    current_timestamp,
    primary key (`id`),
    -- 将username和user_id做唯一的索引
    unique key `idx_username` (`username`) using btree,
    unique key `idx_user_id` (`user_id`) using btree
) ENGINE=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;