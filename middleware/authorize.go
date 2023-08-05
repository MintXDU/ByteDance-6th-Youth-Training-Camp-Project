package middleware

import (
	"errors"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

// 根据token查询用户登录状态，如果已经登录，返回user指针，否则抛出错误

func CheckUserState(token string, db *gorm.DB) (u *dao.User, err error) {
	if len(token) == 0 {
		return nil, errors.New("用户未登录")
	}
	var user dao.User
	res := db.Where("token = ?", token).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// 数据库根据token未查询到用户
		return nil, errors.New("用户不存在，请重新登录")
	}
	if res.Error != nil {
		// 数据库查询失败
		return nil, errors.New("数据库查询失败")
	}
	return &user, nil
}

func GetToken(c *gin.Context) string {
	token := c.DefaultQuery("token", "")
	var err error
	if len(token) == 0 {
		token, err = c.Cookie("token")
		if err != nil {
			log.Println("get token failed!", err)
		}
	}
	return token
}

// 用户鉴权，从url或cookie获取token后查找用户是否登录

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetToken(c)
		if len(token) == 0 {
			c.Abort()
		} else {
			db := service.Connection()
			_, err := CheckUserState(token, db)
			if err != nil {
				c.Abort()
				log.Println("authorize denied!", err)
			}
			c.Next()
		}
	}
}
