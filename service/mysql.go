package service

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
Other methods can be used to connect to the database by calling this method,
which returns a DB object.
*/
func Connection() *gorm.DB {
	// Connect to the mysql database.
	dsn := "xxxx:xxx@tcp(mysql.sqlpub.com:3306)/douyin6th?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Mysql database connection failed, error = " + err.Error())
	} else {
		fmt.Println("Mysql database connection successed.")
	}

	return db
}
