drop table `music`;
drop table `creator`;
drop table `tag`;
drop table `creator_music`;

create table `music` (
    id int(11) unsigned not null AUTO_INCREMENT COMMENT '自增id',
    `name` varchar(255) NOT NULL DEFAULT '未知歌曲' COMMENT '歌曲名字',
    `status` int(11) not null DEFAULT 0 COMMENT '歌曲状态',
    `title` varchar(255) not null DEFAULT '' COMMENT '歌曲标题',
    `hot_score` double not null DEFAULT 0.0 COMMENT '热度打分',
    `creator_ids` varchar(255) not null COMMENT '作者id集合',
    `creator_names` varchar(255) not null COMMENT '作者名字集合',
    `music_url` varchar(255) not null DEFAULT '' COMMENT '音乐的地址',
    `play_time` int(11) not null DEFAULT 0 COMMENT '播放时长',
    `tag_ids` varchar(255) not null DEFAULT '标签id集合' COMMENT '标签ID集合',
    `tag_names` varchar(255) not null DEFAULT '标签名字集合' COMMENT '标签名字集合',
    `image_url` varchar(255) not null DEFAULT '' COMMENT '图片地址',
    `publish_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发行时间',
    `update_time`  timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='音乐表';

create table `tag` (
    id int unsigned not null AUTO_INCREMENT COMMENT '自增id',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT '标签名字',
    `status` int(11) not null DEFAULT 0 COMMENT '标签状态',
    `update_time`  timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `key_name` (`name`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签表';

ALTER TABLE `tag` AUTO_INCREMENT =1

create table `creator` (
    id int unsigned not null AUTO_INCREMENT COMMENT '自增id',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT '作者名字',
    `status` int(11) not null DEFAULT 0 COMMENT '标签状态',
    `image_url` varchar(255) not null DEFAULT 0 COMMENT '作者头像',
    `description` text not null COMMENT '自述',
    `similar_creator` varchar(255) not null COMMENT '相似作者集',
    `fans_num` int(11) not null default 0 COMMENT '粉丝数量',
    `hot_score` double not null DEFAULT 0.0 COMMENT '热度打分',
    `type` int not null DEFAULT 0 COMMENT '作者类型',
    `update_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='作者表';

create table `creator_music` (
    `id` int unsigned not null AUTO_INCREMENT COMMENT '自增id',
    `creator_id` int not null COMMENT '作者id',
    `music_id` int not null COMMENT '歌曲id',
    `status` int(11) not null DEFAULT 0 COMMENT '标签状态',
    `update_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `key_music_creator` (`creator_id`,`music_id`),
    KEY `index_creator` (`creator_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='作者和歌曲的索引表';

create table `user` (
    `id` int unsigned not null AUTO_INCREMENT COMMENT '自增id',
    `username` varchar(255) not null COMMENT '用户',
    `password` varchar(255) not null COMMENT '密码',
    `name` varchar(255) not null COMMENT '名字',
    `status` int(11) not null DEFAULT 0 COMMENT '状态',
    `image_url` varchar(255) not null DEFAULT '' COMMENT '头像地址',
    `create_time` timestamp not null DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
    `last_login_time` timestamp not null DEFAULT '1970-01-01 08:00:01' COMMENT '上次登录时间',
    `last_month_login_num` int not null DEFAULT 0 COMMENT '上个月登录次数',
    `last_stat_time` timestamp not null DEFAULT '1970-01-01 08:00:01' COMMENT '上次统计时间',
    `update_time` timestamp not null DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `key_user` (`username`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

create table `follow_creator` (
    `id` int unsigned not null AUTO_INCREMENT COMMENT '自增id',
    `status` int(11) not null DEFAULT 0 COMMENT '状态',
    `creator_id` int not null COMMENT '作者id',
    `username` varchar(255) not null COMMENT '用户名',
    `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `key_user_creator` (`username`,`creator_id`)

)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='关注';

create table `playlist` (
    `id` int unsigned not null AUTO_INCREMENT COMMENT '自增id',
    `status` int(11) not null DEFAULT 0 COMMENT '状态',
    `user_id` int not null COMMENT '用戶id',
    `name` varchar(255) not null COMMENT '歌单名',
	`image_url` varchar(265) not null DEFAULT '' COMMENT '图片地址',
    `hot_score` double not null DEFAULT 0.0 COMMENT '热度打分',
    `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
    PRIMARY KEY (`id`),
	KEY `key_user` (`user_id`),
    UNIQUE KEY `key_user_play_list` (`user_id`,`name`)

)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='歌单';

create table `playlist_music` (
	`id` int unsigned not null AUTO_INCREMENT COMMENT '自增id',
    `status` int(11) not null DEFAULT 0 COMMENT '状态',
    `user_id` int not null COMMENT '用戶id',
	`music_id` int not null COMMENT '音乐id',
	`playlist_id` int not null COMMENT '歌单id',
    `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
    PRIMARY KEY (`id`),
	KEY `key_user` (`user_id`),
	KEY `key_playlist` (`playlist_id`),
    UNIQUE KEY `key_playlist_music` (`playlist_id`,`music_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='歌单歌曲';
