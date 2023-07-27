# 开发项目

## 需要做的
1. 数据目前是以对象的形式存在内存中的，需要把在内存中的数据存在数据库中

## 项目结构
- /controller 控制层
- /service    业务层
- /dao        数据层
- /public     静态资源

理想的单向调用链：控制层 => 业务层 => 数据层

### GORM
#### 安装
```
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```
#### 连接
见文件 /service/mysql 中的函数 Connection()

其他服务想要连接数据库可调用该函数

#### 查询（有属性为自定义类的）对象
比如查询 `Video` 对象：
```
type Video struct {
	Id             int64  `json:"id,omitempty"`
	UserId         int64  `json:"user_id,omitempty"`
	Author         User   `gorm:"foreignKey:UserId,omitempty"`
	PlayUrl        string `json:"play_url,omitempty"`
	CoverUrl       string `json:"cover_url,omitempty"`
	FavoriteCount  int64  `json:"favorite_count,omitempty"`
	CommentCount   int64  `json:"comment_count,omitempty"`
	IsFavorite     bool   `json:"is_favorite,omitempty"`
	SubmissionTime string `json:"submission_time,omitempty"`
}
```
`Video` 对象包含一个属性 `Author` ，该属性是一个 `User` 自定义类：
```
type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Password      string `json:"password,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}
```
查询方法：
1. 通过`外键`进行约束，例如 `videos.user_id` 依赖于 `users.id`。
2. 在结构体定义处说明本属性通过外键预加载获取，例如`gorm:"foreignKey:UserId,omitempty"`
3. 通过预加载进行查询，例如：`db.Preload("Author").Find(&videos)`
