package controller

import (
	"fmt"
	"net/http"
	"strconv"
	//"sync/atomic"
	"time"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
)

// 用户消息时间
type UserMsgInfo struct {
	State_HasNewMsg   bool
	State_NewMsgCount int64
}

var tempChat = map[string]*UserMsgInfo{}
var Chat = map[string][]dao.Message{}
var messageIdSequence = int64(1)

type ChatResponse struct {
	StatusCode  int32         `json:"status_code"`          // 状态码，0-成功，其他值-失败
	StatusMsg   string        `json:"status_msg,omitempty"` // 返回状态描述
	MessageList []dao.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	//token := c.Query("token") // token
	toUserId := c.Query("to_user_id")            //接收者
	content := c.Query("content")                //内容
	actionType := c.Query("action_type")         //参数类型
	loginUserId := c.GetInt64("user_id")         //发送者
	db := dao.Connection()                       //连接数据库
	targetUserId, err1 := strconv.Atoi(toUserId) //接收者int //参数类型int
	targetActionType, err2 := strconv.Atoi(actionType)
	if err1 != nil || err2 != nil {
		fmt.Printf("message action error : %s", err1)
		return
	}
	switch targetActionType {
	// actionType = 1 发送消息
	case 1:
		//err = dao.SendMessage(fromUserId, toUserId, content, actionType)
		//
		nowtime := time.Now().Unix()
		message := dao.Message{
			Content:    content,             // 消息内容
			CreateTime: nowtime,             // 消息发送时间 yyyy-MM-dd HH:MM:ss
			ToUserId:   int64(targetUserId), // 该消息接收者的id
			FromUserId: loginUserId,         // 该消息发送者的id
		}
		res := db.Create(&message)
		key := genChatKey(loginUserId, int64(targetUserId))
		if _, ok := tempChat[key]; !ok {
			tempChat[key] = &UserMsgInfo{State_HasNewMsg: true, State_NewMsgCount: 0}
		}
		tempChat[key].State_NewMsgCount++
		if res.Error != nil {
			c.JSON(http.StatusOK, dao.Response{StatusCode: 1, StatusMsg: "Send Message 接口错误"})
		} else {
			c.JSON(http.StatusOK, dao.Response{
				StatusCode: 0,
				StatusMsg:  "success",
			})
		}
	default:
		fmt.Sprintf("未定义 actionType=%d", actionType)
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	userId := c.GetInt64("user_id")
	toUserId := c.Query("to_user_id")
	preMsgTime := c.Query("pre_msg_time")
	preMsgTimeInt, errAti1 := strconv.Atoi(preMsgTime)
	toUserIdInt, errAti2 := strconv.Atoi(toUserId)
	if errAti1 != nil || errAti2 != nil {
		res := ChatResponse{
			MessageList: nil,
			StatusCode:  1,
			StatusMsg:   fmt.Sprintf("string to int error:%v", errAti1),
		}
		c.JSON(http.StatusOK, res)
		return
	}

	// 从数据库中获取聊天记录
	chats, err := dao.SelectMessageChat(userId, int64(toUserIdInt))
	if err != nil {
		c.JSON(http.StatusOK, ChatResponse{
			MessageList: nil,
			StatusCode:  1,
			StatusMsg:   fmt.Sprintf("select message chat error:%v", err),
		})
		return
	}

	var res = make([]dao.Message, 0, 100)
	for _, chat := range chats {
		if chat.CreateTime > int64(preMsgTimeInt) {
			temp := dao.Message{
				Content:    chat.Content,
				CreateTime: chat.CreateTime,
				ID:         chat.ID,
				ToUserId:   chat.ToUserId,
				FromUserId: chat.FromUserId,
			}
			res = append(res, temp)
		}
	}
	c.JSON(http.StatusOK, ChatResponse{
		MessageList: res,
		StatusCode:  0,
		StatusMsg:   "获取消息成功",
	})
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
