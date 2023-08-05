package controller

import (
	"errors"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FavoriteActionResponse struct {
	dao.Response
}

type FavoriteListResponse struct {
	dao.Response
	videoList []dao.Video `json:"video_list"`
}

func FavoriteAction(c *gin.Context) {
	token := middleware.GetToken(c)
	db := service.Connection()
	var user *dao.User
	var err error
	user, err = middleware.CheckUserState(token, db)
	if err != nil {
		c.JSON(http.StatusOK, FavoriteActionResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "读取用户信息失败，请重试"}})
		return
	}
	videoId, err := strconv.ParseInt(c.DefaultPostForm("video_id", "0"), 10, 64)
	actionType, err := strconv.ParseInt(c.DefaultPostForm("action_type", "0"), 10, 32)
	if err != nil || (actionType != 1 && actionType != 2) {
		log.Println("[Favourite Action] parse param failed:", err)
		c.JSON(http.StatusOK, FavoriteActionResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "Invalid param"}})
		return
	}

	var video dao.Video
	result := db.Model(&dao.Video{}).Where("id = ?", videoId).First(&video)
	if result.Error != nil {
		log.Println("[Favourite Action] video don't exists, ", result.Error)
		c.JSON(http.StatusOK, FavoriteActionResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "video don't exists"}})
		return
	}
	var favourVideo dao.FavourVideo
	result = db.Model(&dao.FavourVideo{}).Where("user_id = ? and video_id = ?", user.Id, videoId).First(&favourVideo)
	// 点赞记录是否存在
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if actionType == 1 {
			// 点赞记录不存在时点赞
			// 需要更新视频表、点赞表、用户表
			tx := db.Begin()
			now := time.Now().Format(time.RFC3339)
			tx.Model(&dao.Video{}).Where("id = ?", videoId).Update("favorite_count", video.FavoriteCount+1)
			tx.Model(&dao.FavourVideo{}).Create(&dao.FavourVideo{
				UserId:     user.Id,
				VideoId:    videoId,
				FavourTime: now,
			})
			tx.Model(&dao.User{}).Where("id = ?", user.Id).Update("favorite_count", user.FavoriteCount+1)
			if err = tx.Commit().Error; err != nil {
				log.Fatal("[Favourite Action] 数据库事务提交失败")
				tx.Rollback()
				c.JSON(http.StatusOK, FavoriteActionResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "请重试"}})
				return
			}
		}
	} else if favourVideo.VideoId == videoId {
		// 当点赞记录存在时取消点赞
		if actionType == 2 {
			// 取消点赞
			// 需要更新视频表、点赞表、用户表
			tx := db.Begin()
			tx.Model(&dao.Video{}).Where("id = ?", videoId).Update("favorite_count", video.FavoriteCount-1)
			tx.Model(&dao.FavourVideo{}).Delete(&dao.FavourVideo{UserId: user.Id, VideoId: videoId})
			tx.Model(&dao.User{}).Where("id = ?", user.Id).Update("favorite_count", user.FavoriteCount-1)
			if err = tx.Commit().Error; err != nil {
				log.Fatal("[Favourite Action] 数据库事务提交失败")
				tx.Rollback()
				c.JSON(http.StatusOK, FavoriteActionResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "请重试"}})
				return
			}
		}
	}
	c.JSON(http.StatusOK, FavoriteActionResponse{Response: dao.Response{StatusCode: 0, StatusMsg: "success"}})
	return
}

func FavoriteList(c *gin.Context) {
	token := middleware.GetToken(c)
	db := service.Connection()
	user, err := middleware.CheckUserState(token, db)
	if err != nil {
		log.Println("[Favorite List] user don't exist, ", err)
		c.JSON(http.StatusOK, FavoriteListResponse{Response: dao.Response{StatusCode: -1, StatusMsg: err.Error()}})
		return
	}

	var favourVideos []dao.FavourVideo
	favourVideoIds := make([]int64, 0, 10)
	var videos []dao.Video
	db.Model(&dao.FavourVideo{}).Where("user_id = ?", user.Id).Find(&favourVideos)
	for idx, v := range favourVideos {
		log.Println(idx, v)
		favourVideoIds = append(favourVideoIds, v.VideoId)
	}
	db.Model(&dao.Video{}).Where("id in ?", favourVideoIds).Preload("Author").Find(&videos)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: dao.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
	return
}
