# 开发项目

## 上传到1024平台上进行测试
必须要一个文件夹一个文件上传，不然没有权限(1024真不好用)

```shell
go build && ./simple-demo
```
必须是直接运行可执行文件，不能 `go run main.go`!!!

## 需要做的接口
**基础接口**
| API | 名称 | 是否完成 | 完成者 |
| --- | --- | --- | --- |
| /douyin/feed/ | 视频接口流 | 是 | MintXDU |
| /douyin/user/register/ | 用户注册 | 是 | MintXDU |
| /douyin/user/login/ | 用户登录 | 否 |  |
| /douyin/user/ | 用户信息 | 否 |  |
| /douyin/publish/action/ | 投稿接口 | 否 |  |
| /douyin/publish/list/ | 发布列表 | 否 |  |

**互动接口**
| API | 名称 | 是否完成 | 完成者 |
| --- | --- | --- | --- |
| /douyin/favorite/action/ | 赞操作 | 否 |  |
| /douyin/favorite/list/ | 喜欢列表 | 否 |  |
| /douyin/comment/action/ | 评论操作 | 否 |  |
| /douyin/comment/list/ | 评论列表 | 否 |  |

**社交接口**
| API | 名称 | 是否完成 | 完成者 |
| --- | --- | --- | --- |
| /douyin/relation/action/ | 关注操作 | 否 |  |
| /douyin/relation/follow/list/ | 关注列表 | 否 |  |
| /douyin/relation/follower/list/ | 粉丝列表 | 否 |  |
| /douyin/relation/friend/list/ | 好友列表 | 否 |  |
| /douyin/message/action/ | 发送消息 | 否 |  |
| /douyin/message/chat/ | 聊天记录 | 否 |  |

## 需要做的事
>数据目前是以对象的形式存在内存中的，需要把在内存中的数据存在数据库中

为了本地开发与 1024 平台测试方便以及同学们协作开发的实时一致，本项目使用了云数据库，在文件 `./service/mysql.go `中进行配置。

一定记得在更改配置之后不再 commit 该文件，以防止数据库连接数据泄漏。

```
git update-index --assume-unchanged ./service/mysql.go
```

如果你想要开始跟踪该文件的变更，可以使用
```
git update-index --no-assume-unchanged ./service/mysql.go
```
命令将其从忽略列表中移除。

## 项目结构
- /controller 控制层
- /service    业务层
- /dao        数据层
- /public     静态资源

理想的单向调用链：控制层 => 业务层 => 数据层

## 一些修改
1. token 就是 username，token=username+password很奇怪，岂不是如果知道了token就能获取密码了？

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
	Author         User   `gorm:"foreignKey:UserId"`
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
1. 通过`外键`进行约束（很多线上数据库不推荐并且不支持使用外键比如 planetscale ，但是 GORM 中的 BELONG TO 关系推荐使用外键，所以在这里还是使用了外键），例如 `videos.user_id` 依赖于 `users.id`。
2. 在结构体定义处说明本属性通过外键预加载获取，例如`gorm:"foreignKey:UserId,omitempty"`
3. 通过预加载进行查询，例如：`db.Preload("Author").Find(&videos)`
