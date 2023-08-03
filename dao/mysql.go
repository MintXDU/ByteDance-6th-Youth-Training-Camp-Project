package dao

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
Other methods can be used to connect to the database by calling this method,
which returns a DB object.
*/
var db *gorm.DB

func Connection() *gorm.DB {
	// Connect to the mysql database.
	dsn := "root:1104540868@tcp(localhost:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if db.Error != nil {
		fmt.Println("Mysql database connection failed, error = " + db.Error.Error())
	} else {
		fmt.Println("Mysql database connection successed.")
	}

	return db
}
