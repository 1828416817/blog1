package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(32);not null"`
	Telephone string `gorm:"varchar(110);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	username := "root"
	password := "512612lj"
	database := "gin"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset,
	)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database err:" + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}
func isHasTelephone(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	} else {
		return false
	}
}
func main() {
	db := InitDB()
	defer db.Close()
	engine := gin.Default()
	engine.POST("/api/auth/register", func(context *gin.Context) {
		name := context.PostForm("name")
		telephone := context.PostForm("telephone")
		password := context.PostForm("password")
		if len(telephone) != 11 {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "手机号必须为11为",
			})
			return
		}
		if len(password) < 6 {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不能少于6位",
			})
			return
		}
		if len(name) == 0 {
			name = RandomString(10)
		}
		if isHasTelephone(db, telephone) {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"msg": "用户已经存在",
			})
			return
		}
		//创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)
		log.Println(name, telephone, password)
		context.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
		})
	})

	engine.Run(":8081")
}
func RandomString(n int) string {
	letters := []byte("abcdefghijklmnopqrstASDFGHJKLZXXCVBNMQWERTYUIOP")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
