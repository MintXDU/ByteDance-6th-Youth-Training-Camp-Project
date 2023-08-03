package controller

import (
	"net/http"
	"time"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response  dao.Response
	VideoList []dao.Video `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	// Get the latest 10 videos.
	var videos = []dao.Video{}

	db := dao.Connection()
	db.Order("submission_time desc").Limit(10).Preload("Author").Find(&videos)

	c.JSON(http.StatusOK, FeedResponse{
		Response: dao.Response{
			StatusCode: 0,
			StatusMsg:  "Get videos successfully",
		},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
