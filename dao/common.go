package dao

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id             int64  `json:"id,omitempty"`
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	PlayUrl        string `json:"play_url,omitempty"`
	CoverUrl       string `json:"cover_url"`
	FavoriteCount  int64  `json:"favorite_count,omitempty"`
	CommentCount   int64  `json:"comment_count,omitempty"`
	SubmissionTime string `json:"submission_time,omitempty"`
	UserId         int64  `json:"user_id,omitempty"`
	Author         User   `gorm:"foreignKey:UserId" json:"author"`
	IsFavorite     bool   `json:"is_favorite,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	UserId     int64  `json:"user_id"`
	User       User   `gorm:"foreignKey:UserId" json:"user"`
	VideoId    int64  `json:"video_id"`
	Video      Video  `gorm:"foreignKey:VideoId" json:"video"`
	ParentId   int64  `json:"parent_id"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
	CreateDate string `gorm:"-" json:"create_date,omitempty"`
}

type APIUser struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `gorm:"column:name" json:"name,omitempty"`
	Nickname        string `json:"nickname,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty" gorm:"-"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	TotalFavorited  int    `json:"total_favorited,omitempty"`
	WorkCount       int    `json:"work_count,omitempty"`
	FavoriteCount   int    `json:"favorite_count"`
}

type User struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `gorm:"column:name" json:"name,omitempty"`
	Password        string `json:"-"`
	Nickname        string `json:"nickname,omitempty"`
	Token           string `json:"-"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty" gorm:"-"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	TotalFavorited  int    `json:"total_favorited,omitempty"`
	WorkCount       int    `json:"work_count,omitempty"`
	FavoriteCount   int    `json:"favorite_count"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content,omitempty"`
	CreateTime string `gorm:"column:send_time" json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type FavourVideo struct {
	UserId     int64  `json:"user_id" gorm:"primaryKey"`
	VideoId    int64  `json:"video_id" gorm:"primaryKey"`
	FavourTime string `json:"favour_time"`
	//Videos     []Video `json:"videos" gorm:"foreignKey:UserId"`
}

type Relationship struct {
	UserId       int64  `json:"user_id"`
	FollowedId   int64  `json:"followed_id"`
	User         User   `gorm:"foreignKey:UserId" json:"user"`
	FollowedUser User   `gorm:"foreignKey:FollowedId", json:"followed_user"`
	FollowedTime string `json:"followed_time"`
}
