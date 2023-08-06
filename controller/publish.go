package controller

import (
	"fmt"
	"net/http"
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

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	// Find user by token(i.e. username)
	db := service.Connection()
	var user dao.User

	if result := db.Where("name = ?", token).First(&user); result.Error != nil {
		c.JSON(http.StatusOK, dao.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, dao.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	newVideo := dao.Video{
		UserId:         user.Id,
		PlayUrl:        "./public" + finalName,
		SubmissionTime: strconv.FormatInt(time.Now().Unix(), 10),
	}
	db.Create(&newVideo)

	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, dao.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dao.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	token := c.Query("token")
	userID, err := strconv.Atoi(c.Query("user_id"))
	
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: dao.Response{StatusCode: 1, StatusMsg: "user_id format error  " + err.Error()},
		})
		return 
	}

	db := service.Connection()
	var user dao.User

	if result := db.Where("name = ?", token).First(&user); result.Error != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: dao.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}

	var videos = []dao.Video{}
	db.Where("user_id = ?", userID).Order("submission_time desc").Preload("Author").Find(&videos)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: dao.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videos,
	})
}
