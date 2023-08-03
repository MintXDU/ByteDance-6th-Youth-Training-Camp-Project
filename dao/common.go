package dao

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id             int64  `json:"id,omitempty"`              // 视频唯一标识
	UserId         int64  `json:"user_id,omitempty"`         // 视频作者ID
	Author         User   `gorm:"foreignKey:UserId"`         // 视频作者信息
	PlayUrl        string `json:"play_url,omitempty"`        // 视频播放地址
	CoverUrl       string `json:"cover_url,omitempty"`       // 视频封面地址
	FavoriteCount  int64  `json:"favorite_count,omitempty"`  // 视频的点赞总数
	CommentCount   int64  `json:"comment_count,omitempty"`   // 视频的评论总数
	IsFavorite     bool   `json:"is_favorite,omitempty"`     // true-已点赞，false-未点赞
	SubmissionTime string `json:"submission_time,omitempty"` //视频递交时间
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`          // 评论id
	User       User   `json:"user"`                  // 评论用户信息
	Content    string `json:"content,omitempty"`     // 评论内容
	CreateDate string `json:"create_date,omitempty"` // 评论发布日期，格式 mm-dd
}

type User struct {
	Id            int64  `json:"id,omitempty"`             //用户ID
	Name          string `json:"name,omitempty"`           // 用户名称
	Password      string `json:"password,omitempty"`       //用户密码
	FollowCount   int64  `json:"follow_count,omitempty"`   // 关注总数
	FollowerCount int64  `json:"follower_count,omitempty"` //粉丝总数
	IsFollow      bool   `json:"is_follow,omitempty"`      //true-已关注，false-未关注
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
type Message struct {
	Content    string `json:"content,omitempty"`      // 消息内容
	CreateTime int64  `json:"create_time,omitempty"`  // 消息发送时间 yyyy-MM-dd HH:MM:ss
	ID         int64  `json:"id,omitempty"`           // 消息id
	ToUserId   int64  `json:"to_user_id,omitempty"`   // 该消息接收者的id
	FromUserId int64  `json:"from_user_id,omitempty"` // 该消息发送者的id
}

type Follow struct {
	ID          int64
	FollowingID int64 `gorm:"column:following_user_id"`
	FollowedID  int64 `gorm:"column:followed_user_id"`
	Relation    int64 `gorm:"relationship"`
}
