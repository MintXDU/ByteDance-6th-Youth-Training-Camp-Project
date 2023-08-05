package controller

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"

	"net/http"
)

type UserListResponse struct {
	dao.Response
	UserList []dao.User `json:"user_list"`
}

// RelationAction action_type=1关注，action_type=2取消关注
func RelationAction(c *gin.Context) {
	// 校验
	token := middleware.GetToken(c)
	db := service.Connection()
	user, err := middleware.CheckUserState(token, db)
	if err != nil {
		log.Println("[RelationAction] token error", err)
		c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "请重新登录"})
		return
	}
	toUserId, err := strconv.ParseInt(c.PostForm("to_user_id"), 10, 64)
	if err != nil || toUserId == user.Id {
		log.Println("[RelationAction] to_user_id invalid", err)
		c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "to_user_id invalid"})
		return
	}
	actionType, err := strconv.ParseInt(c.PostForm("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		log.Println("[RelationAction] action_type invalid")
		c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "action_type invalid"})
		return
	}
	// 事务更新
	tx := db.Begin()
	var toUser dao.User
	toUserRes := tx.Model(&dao.User{}).Where("id = ?", toUserId).First(&toUser)
	userRes := tx.Model(&dao.User{}).Where("id = ?", user.Id).First(&user)
	// 需要更新：users,relationship
	var relation dao.Relationship //= dao.Relationship{UserId: toUserId, FollowedId: user.Id}
	relationRes := tx.Model(&dao.Relationship{}).Where("user_id = ? and followed_id = ?",
		toUserId, user.Id).First(&relation)
	if toUserRes.RowsAffected == 0 || userRes.RowsAffected == 0 {
		tx.Rollback()
		log.Println("[RelationAction] get user failed")
		c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "retry"})
		return
	}
	if actionType == 1 && relationRes.RowsAffected == 0 {
		// 点赞
		toUserRes = tx.Model(&dao.User{}).Where("id = ?", toUserId).
			Update("follower_count", toUser.FollowerCount+1)
		userRes = tx.Model(&dao.User{}).Where("id = ?", user.Id).
			Update("follow_count", toUser.FollowCount+1)
		relation.UserId = toUser.Id
		relation.FollowedId = user.Id
		relation.FollowedTime = time.Now().Format(time.RFC3339)
		insertRes := tx.Model(&dao.Relationship{}).Select("UserId", "FollowedId", "FollowedTime").Create(&relation)
		if insertRes.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "失败，请重试"})
			return
		}
		tx.Commit()
	} else if actionType == 2 && relationRes.RowsAffected == 1 {
		// 取消
		toUserRes = tx.Model(&dao.User{}).Where("id = ?", toUserId).
			Update("follower_count", toUser.FollowerCount-1)
		userRes = tx.Model(&dao.User{}).Where("id = ?", user.Id).
			Update("follow_count", user.FollowCount-1)
		delRes := tx.Model(&dao.Relationship{}).
			Where("user_id = ? and followed_id = ?", toUser.Id, user.Id).Delete(&relation)
		if delRes.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "失败，请重试"})
			return
		}
		tx.Commit()
	} else {
		tx.Rollback()
		//log.Println("[RelationAction] update database failed")
		//c.JSON(http.StatusOK, dao.Response{StatusCode: -1, StatusMsg: "失败，请重试"})
		//return
	}
	c.JSON(http.StatusOK, dao.Response{StatusCode: 0, StatusMsg: "success"})
	return
}

// FollowList 获取关注列表
func FollowList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "请登录"}})
		return
	}
	token := middleware.GetToken(c)
	db := service.Connection()
	user, err := middleware.CheckUserState(token, db)
	if err != nil || userId != user.Id {
		c.JSON(http.StatusOK, UserListResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "请登录"}})
		return
	}
	var relation = make([]dao.Relationship, 0)
	var userList = make([]dao.User, 0)
	db.Model(&dao.Relationship{}).Where("followed_id = ?", user.Id).Preload("User").Find(&relation)
	for _, val := range relation {
		userList = append(userList, val.User)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: dao.Response{StatusCode: 0},
		UserList: userList,
	})
	return
}

// FollowerList 获取粉丝列表
func FollowerList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "请登录"}})
		return
	}
	token := middleware.GetToken(c)
	db := service.Connection()
	user, err := middleware.CheckUserState(token, db)
	if err != nil || userId != user.Id {
		c.JSON(http.StatusOK, UserListResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "请登录"}})
		return
	}
	var relation = make([]dao.Relationship, 0)
	var userList = make([]dao.User, 0)
	// 查找关注本user的用户
	db.Model(&dao.Relationship{}).Where("user_id = ?", user.Id).Preload("FollowedUser").Find(&relation)
	for _, val := range relation {
		userList = append(userList, val.FollowedUser)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: dao.Response{StatusCode: 0},
		UserList: userList,
	})
	return
}

// FriendList 获取好友列表（互相关注即为好友）
func FriendList(c *gin.Context) {
	token := middleware.GetToken(c)
	db := service.Connection()
	user, err := middleware.CheckUserState(token, db)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{Response: dao.Response{StatusCode: -1}})
		return
	}
	var userIds = make([]int64, 10)
	db.Raw("select r1.followed_id from relationships r1 INNER JOIN relationships r2 ON "+
		"r1.user_id = r2.followed_id and r1.followed_id = r2.user_id where r1.user_id = ?", user.Id).Scan(&userIds)
	var userList = make([]dao.User, 10)
	db.Where("id in ?", userIds).Find(&userList)
	c.JSON(http.StatusOK, UserListResponse{
		Response: dao.Response{StatusCode: 0},
		UserList: userList,
	})
	return
}
