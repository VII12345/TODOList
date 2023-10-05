package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)


type TODO struct{
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

var todos []TODO


func main() {
	r:=gin.Default()

	//add TODO
	r.POST("/todo",func(c *gin.Context){
		var todo TODO
		c.BindJSON(&todo)
		todos = append(todos, todo)
		fmt.Println(todos)
		c.JSON(200,gin.H{"status":"ok"})
	})

	//del TODO
	r.DELETE("/todo/:index",func (c *gin.Context)  {
		index,_:=strconv.Atoi(c.Param("index"))
		todos=append(todos[:index],todos[index+1:]...)
		c.JSON(200,gin.H{"status":"ok"})
	})

	//change
	r.PUT("/todo/:index",func(c *gin.Context){
		index,_:=strconv.Atoi(c.Param("index"))
		var todo TODO
		c.BindJSON(&todo)
		todos[index]=todo
		c.JSON(200,gin.H{"status":"ok"})
	})

	//get
	r.GET("/todo",func(c *gin.Context) {
		c.JSON(200,todos)
	})

	//查询
	r.GET("/todo/:index",func (c *gin.Context)  {
		index,_:=strconv.Atoi(c.Param("index"))
		c.JSON(200,todos[index])
	})
	
	r.Run(":8080")


}
