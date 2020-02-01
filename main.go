package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

func main() {
	db := connectDb()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/user", func(context *gin.Context) {
		var user UserCreateModel
		if context.BindJSON(&user) == nil {
			fmt.Println(user.Password)
			fmt.Println(user.UserName)
			passwordHash,_ := HashPassword(user.Password)
			_,err := db.Exec("insert into users (UserName, PasswordHash) values('"+user.UserName+"',"+"'"+passwordHash+"')");
			if err != nil {
				context.JSON(400, Error {Message: err.Error()})
			} else {
				context.JSON(200,user);
			}
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	defer db.Close()
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


func connectDb() *sql.DB {
	db, err := sql.Open("mysql", "root:@/journey")
	if err != nil {
		fmt.Println("Failed to connect to the database.")
		panic(err.Error())
	}
	return db
}

type UserCreateModel struct {
	UserName string
	Password string
}

type Error struct {
	Message string
}