CREATE DATABASE `douyin` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

INSERT INTO douyin.users (name,follow_count,follower_count,is_follow,password) VALUES
	 ('zhanglei',10,5,1,'douyin'),
	 ('1104540868',0,0,0,'123456');
INSERT INTO douyin.videos (id,user_id,play_url,cover_url,favorite_count,comment_count,is_favorite,submission_time) VALUES
	 (1,1,'https://www.w3schools.com/html/movie.mp4','https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg',0,0,0,'1690428549');
