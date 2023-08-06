package controller

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"

	"net/http"
)

type UserListResponse struct {
	dao.Response
	UserList []dao.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, dao.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, dao.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: dao.Response{
			StatusCode: 0,
		},
		UserList: []dao.User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	username := c.Query("username")
	token := c.Query("token")

	db := service.Connection()
	var user dao.User
	db.Where("name = ?", username).Find(&user)
	if user.Name == token {
		var userList []dao.User
		var users dao.User
		userId := user.Id
		var follower []dao.Follow
		db.Where("users_id = ?", userId).Find(&follower)
		for i := 0; i < len(follower); i++ {
			followerId := follower[i].FollowerId
			var count int64
			db.Where("users_id = ?", followerId).Find(&users)
			if db.Model(&follower[i]).Where("users_id = ? AND follow_id = ?", followerId, userId).Count(&count); count > 0 {
				users.IsFollow = true
			} else {
				users.IsFollow = false
			}
			userList = append(userList, users)
		}
		c.JSON(http.StatusOK, UserListResponse{
			Response: dao.Response{
				StatusCode: 0,
			},
			UserList: userList,
		})

	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: dao.Response{StatusCode: 1, StatusMsg: "token is error!"},
		})
	}

}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	username := c.Query("username")
	token := c.Query("token")

	db := service.Connection()
	var user dao.User
	db.Where("name = ?", username).Find(&user)
	if user.Name == token {
		var userList []dao.User
		var users dao.User
		userId := user.Id
		var friend []dao.Friend
		var follower dao.Follow
		db.Where("user_id = ?", userId).Find(&friend)
		for i := 0; i < len(friend); i++ {
			friendId := friend[i].FriendID
			var count int64
			db.Where("users_id = ?", friendId).Find(&users)
			if db.Model(&follower).Where("users_id = ? AND follow_id = ?", friendId, userId).Count(&count); count > 0 {
				users.IsFollow = true
			} else {
				users.IsFollow = false
			}
			userList = append(userList, users)
		}
		c.JSON(http.StatusOK, UserListResponse{
			Response: dao.Response{
				StatusCode: 0,
			},
			UserList: userList,
		})

	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: dao.Response{StatusCode: 1, StatusMsg: "token is error!"},
		})
	}
}
