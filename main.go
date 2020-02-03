package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
)
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

func main() {
	db = connectDb()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/user", func(context *gin.Context) {
		PostUser(context)
	})

	r.POST("/session", func(context *gin.Context) {
		Login(context);
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	defer db.Close()
}

func Login(context *gin.Context) {
	var loginModel UserLoginModel
	if context.BindJSON(&loginModel) == nil {
		rows, err := db.Query("select * from users where username='"+loginModel.UserName+"' limit 1")
		if err != nil {
			context.JSON(400, Error{Message: err.Error()})
		} else {
			var user User
			if rows.Next() {
				if err := rows.Scan(&user.Id, &user.UserName, &user.PasswordHash); err != nil {
					log.Println(err.Error())
					context.JSON(500, Error{Message:"Server error"})
				} else {
					if CheckPasswordHash(loginModel.Password,user.PasswordHash   ) {
						log.Println("User "+user.UserName+" logged in")
						context.JSON(200,user )
					} else {
						context.JSON(400, Error{Message:"Incorrect password"})
					}
				}
			} else {
				context.JSON(404, Error{Message: "User name or password incorrect."})
			}
		}
	} else {
		context.JSON(400, Error{Message:"Bad request"})
	}
}

func PostUser(context *gin.Context) {
	var user UserCreateModel
	if context.BindJSON(&user) == nil {
		passwordHash,_ := HashPassword(user.Password)
		_,err := db.Exec("insert into users (UserName, PasswordHash) values('"+user.UserName+"',"+"'"+passwordHash+"')");
		if err != nil {
			context.JSON(400, Error {Message: err.Error()})
		} else {
			fmt.Println("User " + user.UserName +" registered")
			context.JSON(200,user);
		}
	}
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

type UserLoginModel struct {
	UserName string
	Password string
}

type User struct {
	Id int
	UserName string
	PasswordHash string
}

type Error struct {
	Message string
}