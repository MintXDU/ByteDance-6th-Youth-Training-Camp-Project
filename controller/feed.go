package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	dao.Response
	VideoList []dao.Video `json:"video_list,omitempty"`
	NextTime  int64       `json:"next_time,omitempty"`
}

const PreloadNum = 10

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	latestTimeStr := c.DefaultQuery("latest_time", "0")
	latestTimeStamp, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		log.Println("[Feed] 请求参数错误", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: dao.Response{
				StatusCode: -1,
				StatusMsg:  "Invalid Parameter",
			},
		})
		return
	}
	if latestTimeStamp == 0 {
		latestTimeStamp = time.Now().Unix()
	}
	latestTime := time.Unix(latestTimeStamp, 0)
	// Get the latest 10 videos.
	var videos = []dao.Video{}
	db := service.Connection()
	result := db.Model(&dao.Video{}).Where("submission_time < ?", latestTime).
		Order("submission_time desc").Limit(PreloadNum).Preload("Author").Find(&videos)
	if result.Error != nil || len(videos) == 0 {
		log.Println("[Feed] search video failed, ", result.Error)
		c.JSON(http.StatusOK, FeedResponse{
			Response: dao.Response{
				StatusCode: -1,
				StatusMsg:  "没有更多视频了",
			},
		})
		return
	}
	nextTime, err := time.Parse(time.RFC3339, videos[len(videos)-1].SubmissionTime)
	if err != nil {
		nextTime = time.Now()
	}
	log.Println("next time:", nextTime)
	c.JSON(http.StatusOK, FeedResponse{
		Response: dao.Response{
			StatusCode: 0,
			StatusMsg:  "Get videos successfully",
		},
		VideoList: videos,
		NextTime:  nextTime.Unix(),
	})
}
