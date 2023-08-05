package controller

import (
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/service"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
)

type ChatResponse struct {
	dao.Response
	MessageList []dao.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	token := middleware.GetToken(c)
	db := service.Connection()
	user, err := middleware.CheckUserState(token, db)
	if err != nil {
		c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "login"})
		return
	}
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, actErr := strconv.ParseInt(c.Query("action_type"), 10, 64)
	content := c.Query("content")
	if err != nil || len(strings.TrimSpace(content)) == 0 || actErr != nil {
		c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "invalid param"})
		return
	}
	if actionType == 1 {
		var msg = dao.Message{
			ToUserId:   toUserId,
			FromUserId: user.Id,
			Content:    content,
			CreateTime: time.Now().Format(time.RFC3339),
		}
		result := db.Create(&msg)
		if result.Error != nil {
			log.Println("[MessageAction")
			c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "重试"})
		}
		c.JSON(http.StatusOK, dao.Response{StatusCode: 0})
		return
	} else {
		c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "invalid param"})
		return
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	token := middleware.GetToken(c)
	db := service.Connection()
	user, err := middleware.CheckUserState(token, db)
	if err != nil {
		c.JSON(http.StatusForbidden, ChatResponse{
			Response: dao.Response{StatusCode: -1},
		})
		return
	}

	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, ChatResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "invalid param"},
		})
		return
	}
	var msgList []dao.Message
	db.Where("(from_user_id = ? and to_user_id = ?) or (to_user_id = ? and from_user_id = ?)",
		user.Id, toUserId, user.Id, toUserId).Find(&msgList)
	c.JSON(http.StatusOK, ChatResponse{
		Response:    dao.Response{StatusCode: 0},
		MessageList: msgList,
	})
	return
}
