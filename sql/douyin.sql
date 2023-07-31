CREATE DATABASE `douyin` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

-- select database douyin
USE douyin;

-- douyin.users definition

CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT NULL,
  `follow_count` int DEFAULT NULL,
  `follower_count` int DEFAULT NULL,
  `is_follow` tinyint(1) DEFAULT NULL,
  `password` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- douyin.videos definition

CREATE TABLE `videos` (
  `id` int NOT NULL,
  `user_id` int DEFAULT NULL,
  `play_url` varchar(100) DEFAULT NULL,
  `cover_url` varchar(100) DEFAULT NULL,
  `favorite_count` int DEFAULT NULL,
  `comment_count` int DEFAULT NULL,
  `is_favorite` tinyint(1) DEFAULT NULL,
  `submission_time` varchar(100) DEFAULT NULL,
  KEY `videos_submission_time_IDX` (`submission_time`) USING BTREE,
  KEY `videos_FK` (`user_id`),
  CONSTRAINT `videos_FK` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- comment table
CREATE TABLE `comment`(
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY ,
    `user_id` INT NOT NULL ,
    `video_id` INT NOT NULL,
    `content` VARCHAR(4096) NOT NULL ,
    `parent_id` INT DEFAULT 0,
    `create_time` DATETIME NOT NULL,
    CONSTRAINT `comment_user_FK` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT,
    CONSTRAINT `comment_video_FK` FOREIGN KEY (`video_id`) REFERENCES `videos`(`id`) ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO douyin.users (name,follow_count,follower_count,is_follow,password) VALUES
	 ('zhanglei',10,5,1,'douyin'),
	 ('1104540868',0,0,0,'123456');
INSERT INTO douyin.videos (id,user_id,play_url,cover_url,favorite_count,comment_count,is_favorite,submission_time) VALUES
	 (1,1,'https://www.w3schools.com/html/movie.mp4','https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg',0,0,0,'1690428549');