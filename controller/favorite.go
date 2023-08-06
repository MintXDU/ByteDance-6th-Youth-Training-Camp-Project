package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	// 参数校验
	token := c.Query("token")
	videoID, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		c.JSON(http.StatusOK, dao.Response{StatusCode: 1, StatusMsg: "video_id error"})
		return 
	}

	actionType := c.Query("action_type")
	if !(actionType == "1" || actionType == "2") {
		c.JSON(http.StatusOK, dao.Response{StatusCode: 1, StatusMsg: "action_type error"})
		return
	}

	db := service.Connection()
	var user dao.User
	var video dao.Video

	if result := db.Where("name = ?", token).First(&user); result.Error != nil {
		c.JSON(http.StatusOK, dao.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	if result := db.Where("id = ?", videoID).First(&video); result.Error != nil {
		c.JSON(http.StatusOK, dao.Response{StatusCode: 1, StatusMsg: "video doesn't exist"})
		return
	}

	// 点赞操作处理
	if actionType == "1" {
		var count int64
		// 如果数据库中没有记载，就加入点赞数据
		if db.Where("user_id = ? and video_id = ?", user.Id, videoID).Count(&count); count == 0 {
			obj := &dao.Favorite{
				UserId: user.Id,
				VideoId: int64(videoID),
			}
			db.Create(obj)
		}
		c.JSON(http.StatusOK, dao.Response{StatusCode: 0, StatusMsg: "success"})
		return 
	}

	// 取消点赞
	if result := db.Where("user_id = ? and video_id = ?", user.Id, videoID).Delete(&dao.Favorite{}); result.Error != nil {
		c.JSON(http.StatusOK, dao.Response{StatusCode: 1, StatusMsg: "operation error"})
	} else {
		c.JSON(http.StatusOK, dao.Response{StatusCode: 0, StatusMsg: "success"})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	userID, err := strconv.Atoi(c.Query("user_id"))

	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: dao.Response{StatusCode: 1, StatusMsg: "user_id format error"},
		})		
	} 
	
	db := service.Connection()
	var user dao.User
	var videos = []dao.Video{}
	var favorites = []dao.Favorite{}

	if result := db.Where("name = ?", token).First(&user); result.Error != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: dao.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}	
	
	db.Where("user_id = ?", userID).Find(&favorites)
	for i := 0; i < len(favorites); i++ {
		video := dao.Video{}
		result := db.Where("id = ?", favorites[i].VideoId).Preload("Author").Find(&video)
		if result.Error != nil {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: dao.Response{
					StatusCode: 1,
					StatusMsg: "database error",
				},
			})
			return 
		}
		videos = append(videos, video)
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: dao.Response{
			StatusCode: 0,
			StatusMsg: "success",
		},
		VideoList: videos,
	})
}
