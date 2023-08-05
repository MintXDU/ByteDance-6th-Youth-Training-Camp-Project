package controller

import (
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CommentListResponse struct {
	dao.Response
	CommentList []dao.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	dao.Response
	Comment dao.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	db := service.Connection()
	var user *dao.User
	var err error

	token := middleware.GetToken(c)
	user, err = middleware.CheckUserState(token, db)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "用户未登录"},
		})
		return
	}

	actionType, actErr := strconv.ParseInt(c.DefaultPostForm("action_type", "0"), 10, 32)
	videoId, vidErr := strconv.ParseInt(c.PostForm("video_id"), 10, 64)
	commentText := strings.TrimSpace(c.PostForm("comment_text"))
	commentId, cidErr := strconv.ParseInt(c.PostForm("comment_id"), 10, 64)
	// 参数校验
	if actErr != nil || (actionType == 1 && (vidErr != nil || len(commentText) == 0)) ||
		(actionType == 2 && cidErr != nil) ||
		(actionType != 1 && actionType != 2) {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "Invalid param"},
		})
		return
	}
	// 校验 videoId是否合法
	var video dao.Video
	res := db.Where("id = ?", videoId).First(&video)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// 数据库根据token未查询到用户
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "视频找不到了!"},
		})
		return
	} else if res.Error != nil {
		// 数据库查询失败
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "数据库查询失败"},
		})
		return
	}

	switch actionType {
	//发布评论
	case 1:
		// 校验context是否合法，不允许发布空评论
		if len(commentText) == 0 {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: dao.Response{StatusCode: -1, StatusMsg: "评论内容不能为空"},
			})
			return
		}
		log.Println("发布评论")
		var comment = dao.Comment{
			UserId:     user.Id,
			Content:    commentText,
			VideoId:    video.Id,
			CreateTime: time.Now().Format(time.RFC3339),
		}
		// 发布评论需要更新评论表、视频表
		tx := db.Begin()
		createRes := tx.Model(&dao.Comment{}).Select("UserId", "Content", "VideoId", "CreateTime").Create(&comment)
		updateRes := tx.Model(&dao.Video{}).Where("id = ?", video.Id).Update("comment_count", video.CommentCount+1)
		if createRes.Error != nil || updateRes.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, CommentActionResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "请重试"}})
			return
		}
		tx.Commit()

		// return data
		parseTime, err := time.Parse(time.RFC3339, comment.CreateTime)
		if err != nil {
			log.Println("[Comment Action] parseTime failed!", err)
			return
		}
		comment.CreateDate = fmt.Sprintf("%02d-%02d", parseTime.Month(), parseTime.Day())
		comment.User = *user
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: dao.Response{StatusCode: 0},
			Comment:  comment,
		})
		return
	// 删除评论
	case 2:
		log.Println("删除评论")
		tx := db.Begin()
		delRes := tx.Delete(&dao.Comment{}, commentId)
		updateRes := tx.Model(&dao.Video{}).Where("id = ?", video.Id).
			Update("comment_count", video.CommentCount-1)
		if delRes.Error != nil || updateRes.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: dao.Response{StatusCode: -1, StatusMsg: "评论数据库删除失败，请重试"},
			})
			return
		} else {
			tx.Commit()
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
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoId := c.Query("video_id")
	db := service.Connection()
	var comments []dao.Comment
	//查询一级评论列表
	res := db.Where("video_id = ? and parent_id = 0", videoId).Preload("User").Find(&comments)
	if res.Error != nil {
		log.Println("search video error:", res.Error)
		c.JSON(http.StatusOK, CommentListResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "评论查询失败！"},
		})
		return
	} else {
		for idx, v := range comments {
			parseTime, err := time.Parse(time.RFC3339, v.CreateTime)
			if err != nil {
				log.Println("[comment list] parse time error, ", err)
				return
			}
			comments[idx].CreateDate = fmt.Sprintf("%02d-%02d", parseTime.Month(), parseTime.Day())
		}
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    dao.Response{StatusCode: 0},
			CommentList: comments,
		})
		return
	}
}
