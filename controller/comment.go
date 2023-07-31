package controller

import (
	"errors"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CommentListResponse struct {
	dao.Response
	CommentList []dao.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	dao.Response
	Comment dao.Comment `json:"comment,omitempty"`
}

func CheckUserState(c *gin.Context, db *gorm.DB) (u *dao.User, err error) {
	token := c.Query("token")
	// 无token，用户未登录
	if len(token) == 0 {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "请先登录"},
		})
		return nil, errors.New("用户未登录")
	}
	var user dao.User
	res := db.Where("token = ?", token).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// 数据库根据token未查询到用户
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "用户不存在，请重新登录!"},
		})
		return nil, errors.New("用户不存在，请重新登录!")
	}
	if res.Error != nil {
		// 数据库查询失败
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "数据库查询失败"},
		})
		return nil, errors.New("数据库查询失败")
	}
	return &user, nil
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	db := service.Connection()
	user, err := CheckUserState(c, db)
	if err != nil {
		return
	}
	log.Println("find user with token, user:", user)
	videoId := c.PostForm("video_id")
	commentText := strings.TrimSpace(c.PostForm("comment_text"))
	commentId := c.PostForm("comment_id")

	// 校验 videoId是否合法
	var video dao.Video
	res := db.Where("id = ?", videoId).First(&video)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// 数据库根据token未查询到用户
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "视频找不到了!"},
		})
		return
	}
	if res.Error != nil {
		// 数据库查询失败
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "数据库查询失败"},
		})
		return
	}

	switch actionType := c.PostForm("action_type"); actionType {
	//发布评论
	case "1":
		// 校验context是否合法，不允许发布空评论
		if len(commentText) == 0 {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: dao.Response{StatusCode: -1, StatusMsg: "评论内容不能为空"},
			})
			return
		}
		log.Println("发布评论")
		var comment = dao.Comment{}
		comment.Content = commentText
		comment.CreateDate = "07-05"
		comment.User = *user
		comment.Video = video
		res = db.Create(&comment)
		if res.Error != nil {
			// 插入失败
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: dao.Response{StatusCode: -1, StatusMsg: "评论数据库插入失败，请重试"},
			})
			return
		}
		return
	// 删除评论
	case "2":
		log.Println("删除评论")
		// 校验commentId参数
		commentId, err := strconv.ParseInt(commentId, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: dao.Response{StatusCode: -1, StatusMsg: "非法参数，请重试"},
			})
		}
		res = db.Delete(&dao.Comment{}, commentId)
		if res.Error != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: dao.Response{StatusCode: -1, StatusMsg: "评论数据库删除失败，请重试"},
			})
			return
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: dao.Response{StatusCode: 0, StatusMsg: "删除成功"},
			})
			return
		}
	default:
		log.Println("invalid op")
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "非法参数"},
		})
		return
	}
	//db := service.Connection()
	//var user dao.User
	//if user, exist := usersLoginInfo[token]; exist {
	//	if actionType == "1" {
	//		text := c.Query("comment_text")
	//		c.JSON(http.StatusOK, CommentActionResponse{Response: dao.Response{StatusCode: 0},
	//			Comment: dao.Comment{
	//				Id:         1,
	//				User:       user,
	//				Content:    text,
	//				CreateDate: "05-01",
	//			}})
	//		return
	//	}
	//	c.JSON(http.StatusOK, dao.Response{StatusCode: 0})
	//} else {
	//	c.JSON(http.StatusOK, dao.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoId := c.Query("video_id")
	db := service.Connection()
	// 用户登录校验
	_, err := CheckUserState(c, db)
	if err != nil {
		return
	}
	var comments []dao.Comment
	//查询一级评论列表
	res := db.Where("video_id = ? and parent_id = 0", videoId).Find(&comments)
	if res.Error != nil {
		log.Println("search video error:", res.Error)
		c.JSON(http.StatusOK, CommentListResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "评论查询失败！"},
		})
		return
	} else {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    dao.Response{StatusCode: 0},
			CommentList: comments,
		})
		return
	}
	//c.JSON(http.StatusOK, CommentListResponse{
	//	Response:    dao.Response{StatusCode: 0},
	//	CommentList: DemoComments,
	//})
}
