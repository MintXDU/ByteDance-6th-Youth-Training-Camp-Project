DROP DATABASE `douyin`;
CREATE DATABASE `douyin` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

-- select database douyin
USE douyin;

-- douyin.users definition

CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL UNIQUE ,
  `password` binary(16) NOT NULL ,       -- char最大长度255
  `nickname` varchar(64) UNIQUE DEFAULT NULL,
  `token` char(255) UNIQUE DEFAULT NULL,
  `follow_count` int DEFAULT 0,
  `follower_count` int DEFAULT 0,
  `avatar` varchar(1024) DEFAULT 'public/defaultHeader.jpg',
  `signature` varchar(256) DEFAULT '',
  `total_favorited` int default 0,
  `work_count` int default 0,
  `favorite_count` int default  0,
  `background_image` varchar(1024) DEFAULT 'public/defaultBackground.jpg',
  `signup_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_user_token`(`token`)   -- 为token创建索引
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- show create table users;
-- show index from users;

-- douyin.videos definition
CREATE TABLE `videos` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY ,
  `title` varchar(256) NOT NULL ,
  `description` varchar(1024) DEFAULT NULL,
  `play_url` varchar(1024) NOT NULL ,
  `user_id` int NOT NULL ,
  `cover_url` varchar(1024) DEFAULT NULL,
  `favorite_count` int DEFAULT 0,
  `comment_count` int DEFAULT 0,
  `submission_time` datetime DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_video_publishTime`(`submission_time`),        -- 为publish_time创建索引
  CONSTRAINT `videos_user_FK` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE     -- 级联删除
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- favour_video点赞表
CREATE TABLE `favour_videos`(
    `user_id` int NOT NULL ,
    `video_id` int NOT NULL ,
    `favour_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`, `video_id`),
    INDEX `idx_favour_videos_time`(`favour_time`),
    CONSTRAINT `favour_user_FK` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ,
    CONSTRAINT `favour_video_FK` FOREIGN KEY (`video_id`) REFERENCES `videos`(`id`) ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- relationship关系表
CREATE TABLE `relationships`(
    `user_id` int NOT NULL ,
    `followed_id` int NOT NULL comment 'user_id关注了哪些人',
    `followed_time` datetime DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`, `followed_id`),
    CONSTRAINT `relationship_user_FK` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT ,
    CONSTRAINT `relationship_user_followed_FK` FOREIGN KEY (`followed_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT ,
    INDEX `idx_relationships_followedId_followedTime`(`followed_id`, `followed_time`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- comment table
CREATE TABLE `comments`(
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY ,
    `user_id` INT NOT NULL ,
    `video_id` INT NOT NULL,
    `content` VARCHAR(4096) NOT NULL ,
    `parent_id` INT DEFAULT 0,
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `comment_user_FK` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT,
    CONSTRAINT `comment_video_FK` FOREIGN KEY (`video_id`) REFERENCES `videos`(`id`) ON DELETE CASCADE,
    INDEX `idx_comments_videoId_createTime`(`video_id`, `create_time`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- message消息表
CREATE TABLE `messages`(
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `to_user_id` INT NOT NULL ,
    `from_user_id` INT NOT NULL ,
    `content` VARCHAR(1024) NOT NULL DEFAULT '',
    `send_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT `message_user_to_FK` FOREIGN KEY (`to_user_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT ,
    CONSTRAINT `message_user_from_FK` FOREIGN KEY (`from_user_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT ,
    INDEX `idx_message_toUser_time`(`to_user_id`, `send_time`),
    INDEX `idx_message_fromUser_time`(`from_user_id`, `send_time`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO douyin.users (name, password) VALUES
	 ('zhanglei', 'douyin'),
	 ('1104540868', '123456');
INSERT INTO douyin.videos (id, title, user_id,play_url,cover_url,favorite_count,comment_count) VALUES
	 (1,'test',1,'https://www.w3schools.com/html/movie.mp4','https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg',0,0);