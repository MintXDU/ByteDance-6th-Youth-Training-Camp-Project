package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/service"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]dao.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	dao.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	dao.Response
	User dao.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username

	// Check if the username exists in the database.
	db := service.Connection()
	var user dao.User
	var count int64

	if db.Model(&user).Where("name = ?", token).Count(&count); count > 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: dao.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		newUser := dao.User{
			Name:     username,
			Password: password,
		}
		db.Create(&newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: dao.Response{
				StatusCode: 0,
				StatusMsg:  "User registration successful",
			},
			Token: username,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: dao.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: dao.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: dao.Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: dao.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
