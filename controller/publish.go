package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	dao.Response
	VideoList []dao.Video `json:"video_list"`
}

type PublishResponse struct {
	dao.Response
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := middleware.GetToken(c)
	title := c.PostForm("title")
	// Find user by token(i.e. username)
	db := service.Connection()
	var user *dao.User
	user, err := middleware.CheckUserState(token, db)
	if err != nil {
		c.JSON(http.StatusOK, PublishResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "未登录"}})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, dao.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// save file logic
	filename := filepath.Base(data.Filename)
	timeval := time.Now().Unix()
	finalName := fmt.Sprintf("%d_%d_%s", user.Id, timeval, filename)
	saveFile := filepath.Join("./public/", finalName)
	_, err = os.Stat(saveFile)
	if err == nil {
		// file exists
		c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "请不要重复上传视频"})
		return
	}
	newVideo := dao.Video{
		UserId:         user.Id,
		Title:          title,
		PlayUrl:        "./public" + finalName,
		CoverUrl:       "http://www.baidu.com",
		SubmissionTime: time.Now().Format(time.RFC3339),
	}
	tx := db.Begin()
	result := tx.Model(&dao.Video{}).Select("UserId", "Title", "PlayUrl", "SubmissionTime", "CoverUrl").Create(&newVideo)
	if result.Error != nil {
		log.Println("save video to database error:", result.Error)
		c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "上传失败，请重试"})
		return
	}
	// 保存文件失败，回滚数据库
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dao.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 提交事物
	tx.Commit()
	c.JSON(http.StatusOK, dao.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
	return
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	db := service.Connection()
	userId, err := strconv.ParseInt(c.DefaultQuery("user_id", "0"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: dao.Response{
				StatusCode: -1,
				StatusMsg:  "用户不存在",
			},
		})
		return
	}
	var user dao.User
	result := db.Model(&dao.User{}).Where("id = ?", userId).First(&user)
	if result.Error != nil {
		log.Println("[PublishList] find user with user_id failed:", result.Error)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: dao.Response{
				StatusCode: -1,
				StatusMsg:  "用户不存在",
			},
		})
		return
	}
	var videoLists []dao.Video
	result = db.Model(&dao.Video{}).Where("user_id = ?", userId).Preload("Author").Find(&videoLists)
	if result.Error != nil {
		log.Println(result.Error)
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: dao.Response{
			StatusCode: 0,
		},
		VideoList: videoLists,
	})
	return
}
