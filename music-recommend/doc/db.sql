create table `music` (
    id int(11) unsigned not null AUTO_INCREMENT COMMENT '自增id',
    `name` varchar(255) NOT NULL DEFAULT '未知歌曲' COMMENT '歌曲名字',
    `status` int(11) not null DEFAULT 0 COMMENT '歌曲状态',
    `title` varchar(255) not null DEFAULT '' COMMENT '歌曲标题',
    `hot_score` double not null DEFAULT 0.0 COMMENT '热度打分',
    `creator_id` int not null COMMENT '作者id',
    `music_url` varchar(255) not null DEFAULT '' COMMENT '音乐的连接',
    `tag_ids` varchar(255) not null DEFAULT '标签id集合' COMMENT '标签ID集合',
    `tag_names` varchar(255) not null DEFAULT '标签名字集合' COMMENT '标签名字集合',
    `play_time` int(11) not null DEFAULT 0 COMMENT '播放时长',
    `image_url` varchar(255) not null DEFAULT '' COMMENT '图片地址',
    `publish_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发行时间',
    `update_time`  timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `key_creator` (`creator_id`),
    KEY `index_creator` (`creator_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='音乐表';

create table `tag` (
    id int unsigned not null AUTO_INCREMENT COMMENT '自增id',
    `name` varchar(255) NOT NULL DEFAULT '' COMMENT '标签名字',
    `status` int(11) not null DEFAULT 0 COMMENT '标签状态',
    `update_time`  timestamp NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `key_name` (`name`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签表';