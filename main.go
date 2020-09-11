package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

type user struct {
	gorm.Model
	Name string
	Age  int
}

func dbInit() {
	db, err := gorm.Open("sqlite3", "user.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&user{})
}

func dbInsert(name string, age int) {
	db, err := gorm.Open("sqlite3", "user.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Create(&user{Name: name, Age: age})
}

//DB全取得
func dbGetAll() []user {
	db, err := gorm.Open("sqlite3", "user.db")
	if err != nil {
		panic("ERROR!(dbGetAll())")
	}
	var users []user
	db.Order("created_at desc").Find(&users)
	db.Close()
	return users
}

//DB削除
func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "user.db")
	if err != nil {
		panic("ERROR!（dbDelete)")
	}
	var users user
	db.First(&users, id)
	db.Delete(&users)
	db.Close()
}

//DB一つ取得
func dbGetOne(id int) user {
	db, err := gorm.Open("sqlite3", "user.db")
	if err != nil {
		panic("ERROR!(dbGetOne())")
	}
	var users user
	db.First(&users, id)
	db.Close()
	return users
}

//UPDATE
func dbUpdate(id int, name string, age int) {
	db, err := gorm.Open("sqlite3", "user.db")
	if err != nil {
		panic("ERROR")
	}
	var users user
	db.First(&users, id)
	users.Name = name
	users.Age = age
	db.Save(&users)
	db.Close()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("views/*.html")

	dbInit()

	//READ
	r.GET("/", func(c *gin.Context) {
		users := dbGetAll()
		c.HTML(200, "index.html", gin.H{"users": users})
	})

	//CREATE
	r.POST("/new", func(c *gin.Context) {
		name := c.PostForm("name")
		a := c.PostForm("age")
		age, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		dbInsert(name, age)
		c.Redirect(302, "/")
	})

	//削除確認
	r.GET("/check/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		users := dbGetOne(id)
		c.HTML(200, "delete.html", gin.H{"user": users})
	})

	//DELETE
	r.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		dbDelete(id)
		c.Redirect(302, "/")
	})

	r.GET("/edit/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		users := dbGetOne(id)
		c.HTML(200, "edit.html", gin.H{"user": users})
	})

	r.POST("/update/:id", func(c *gin.Context) {
		name := c.PostForm("name")
		a := c.PostForm("age")
		age, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		dbUpdate(id, name, age)
		c.Redirect(302, "/")
	})

	r.Run(":443")
}
