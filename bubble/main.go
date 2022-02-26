package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Status bool `json:"status"`
}

//声明一个DB的全局变量
var (
	DB *gorm.DB
)

func initMySQL() (err error) {
	//连接数据库，用户名:密码，数据库名
	dsn := "root:123456@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return err
}

func main() {
	//创建数据库
	//连接数据库
	err := initMySQL()
	if err != nil{
		panic(err)
	}
	//模型绑定
	err = DB.AutoMigrate(&Todo{})
	if err != nil {
		panic(err)
	}

	//默认引擎
	r := gin.Default()
	//告诉gin框架模板文件引用的静态文件去哪里找
	r.Static("/static", "static")
	//加载模板文件路径
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	//v1
	v1Group := r.Group("v1")
	{
		//待办事项
		//添加
		v1Group.POST("/todo", func(c *gin.Context) {
			//前端页面填写待办事项，点击提交，会发送请求到这里
			//1、从请求中把数据拿出来
			var todo Todo
			_ = c.BindJSON(&todo)
			//2、存入数据库
			err = DB.Create(&todo).Error
			//3、返回响应
			if err != nil{
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			}else {
				c.JSON(http.StatusOK, todo)
			}
		})
		//查看所有待办事项
		v1Group.GET("/todo", func(c *gin.Context) {
			var todoList []Todo
			err = DB.Find(&todoList).Error
			if err != nil{
				c.JSON(http.StatusOK, gin.H{
					"error": err.Error(),
				})
			}else {
				c.JSON(http.StatusOK, todoList)
			}
		})
		//查看某一待办事项
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})
		//修改某一事项
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, _ := c.Params.Get("id")
			var todo Todo
			DB.Where("id=?", id).First(&todo)
			_ = c.BindJSON(&todo)
			DB.Save(&todo)
		})

		//删除
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, _ := c.Params.Get("id")
			DB.Where("id=?", id).Delete(Todo{})
		})
	}

	err = r.Run()
	if err != nil {
		return
	}
}
