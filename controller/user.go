package controller

import (
	"crypto/md5"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

const TokenExpiredDuration = 15

type UserLoginResponse struct {
	dao.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	dao.Response
	User dao.APIUser `json:"user"`
}

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	db := service.Connection()
	var count int64

	// Check if the username exists in the database.
	if db.Model(&dao.User{}).Where("name = ?", username).Count(&count); count > 0 {
		log.Printf("[Register] %v is exists in the database", username)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: dao.Response{StatusCode: -1, StatusMsg: "用户已存在"},
		})
		return
	} else {
		tokenString, err := middleware.GenToken(username, TokenExpiredDuration)
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "create user failed"}})
			return
		}
		passwdMd5 := md5.Sum([]byte(password))
		user := dao.User{
			Name:     username,
			Password: string(passwdMd5[:]),
			Token:    *tokenString,
		}
		result := db.Select("Name", "Password", "Token").Create(&user)
		// 判断是否新建用户成功
		if result.Error == nil {
			log.Println("insert rows", result.RowsAffected)
			c.JSON(http.StatusOK, UserLoginResponse{
				UserId:   user.Id,
				Token:    user.Token,
				Response: dao.Response{StatusCode: 0},
			})
			return
		} else {
			log.Printf("create user error, %v", result.Error)
			c.JSON(http.StatusOK, UserLoginResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "create user failed"}})
			return
		}
	}
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	db := service.Connection()
	var user dao.User
	passwdMd5 := md5.Sum([]byte(password))
	result := db.Model(&dao.User{}).Where("name = ? and password = ?", username, string(passwdMd5[:])).First(&user)
	// if user doesn't exist
	if result.Error != nil {
		log.Printf("[Login] 查找用户失败, %v", result.Error)
		c.JSON(http.StatusOK, UserLoginResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "user don't exist"}})
		return
	}
	tokenString, err := middleware.GenToken(username, TokenExpiredDuration)
	if err == nil {
		// update token
		result = db.Model(&dao.User{}).Where("name = ?", username).Update("token", *tokenString)
		if result.Error == nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: dao.Response{
					StatusCode: 0,
					StatusMsg:  "User login successfully",
				},
				UserId: user.Id,
				Token:  *tokenString,
			})
			return
		}
	}
	log.Printf("[Login] update token failed:%v, %v", result.Error, err)
	c.JSON(http.StatusOK, UserLoginResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "update token failed"}})
	return
}

func UserInfo(c *gin.Context) {
	token := middleware.GetToken(c)
	userIdStr := c.DefaultQuery("user_id", "0")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if len(token) == 0 || err != nil {
		c.JSON(http.StatusOK, UserResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "用户未登录"}})
		return
	}

	// Find user by token(i.e. username)
	db := service.Connection()
	var user dao.User

	// 查看当前用户是否登录
	result := db.Model(&dao.User{}).Where("token = ?", token).First(&user)
	if result.Error != nil {
		log.Printf("[UserInfo] find user with token failed, %v", result.Error)
		c.JSON(http.StatusOK, UserResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "登录信息已过期，请重新登录"}})
		return
	}
	if userId == 0 {
		userId = user.Id
	}
	// 查找对应用户信息
	var apiUser dao.APIUser
	result = db.Model(&dao.User{}).Where("id = ?", userId).First(&apiUser)
	if result.Error != nil {
		log.Printf("[UserInfo] find user with user_id failed, %v", result.Error)
		c.JSON(http.StatusOK, UserResponse{Response: dao.Response{StatusCode: -1, StatusMsg: "用户不存在"}})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: dao.Response{
			StatusCode: 0,
			StatusMsg:  "User exist",
		},
		User: apiUser,
	})
	return
}
