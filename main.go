package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
)


type TODO struct{
	Content string `json:"content"`
	Done    string   `json:"done"`
}

var todos []TODO

func readTODO(){
	file,_:=os.OpenFile("TODO.txt",os.O_RDWR,0666)
	read,_:=io.ReadAll(file)
	readline:=strings.Split(string(read),"\n")
	lenth:=len(readline)/2
	var todo TODO
	for i:=0;i<lenth;i++{
		todo.Content,todo.Done=readline[i*2],readline[i*2+1]
		todos=append(todos,todo)
	}
	file.Close()
}
func add(text TODO){
	file,_:=os.OpenFile("TODO.txt",os.O_APPEND,0666)
	file.WriteString(text.Content)
	file.WriteString("\n")
	file.WriteString(text.Done)
	file.WriteString("\n")
	file.Close()
}

func del(ind int){
	file,_:=os.OpenFile("TODO.txt",os.O_RDWR,0666)
	read,_:=io.ReadAll(file)
	readline:=strings.Split(string(read),"\n")
	readline[ind*2],readline[ind*2+1]="has been deleted","has been deleted"
	file.Truncate(0)
	file.Seek(0,0)
	for _,i:=range readline{
		file.WriteString(i)
		file.WriteString("\n")
	}
	file.Close()
}

func change(text TODO,ind int){
	file,_:=os.OpenFile("TODO.txt",os.O_RDWR,0666)
	read,_:=io.ReadAll(file)
	readline:=strings.Split(string(read),"\n")
	readline[ind*2],readline[ind*2+1]=text.Content,text.Done
	file.Truncate(0)
	file.Seek(0,0)
	for _,i:=range readline{
		file.WriteString(i)
		file.WriteString("\n")
	}
	file.Close()
}


func main() {
	r:=gin.Default()

	readTODO()


	//add TODO
	r.POST("/todo",func(c *gin.Context){
		var todo TODO
		c.BindJSON(&todo)
		todos = append(todos, todo)
		fmt.Println(todos) 
		add(todo)
		c.JSON(200,gin.H{"status":"ok"})
	})

	//del TODO
	r.DELETE("/todo/:index",func (c *gin.Context)  {
		index,_:=strconv.Atoi(c.Param("index"))
		todos[index].Content=`has been deleted`
		todos[index].Done=`has been deleted`
		del(index)
		c.JSON(200,gin.H{"status":"ok"})
	})

	//更新
	r.POST("/todo/:index",func(c *gin.Context){
		index,_:=strconv.Atoi(c.Param("index"))
		var todo TODO
		c.BindJSON(&todo)
		todos[index]=todo
		change(todo,index)
		c.JSON(200,gin.H{"status":"ok"})
	})

	//列出
	r.GET("/todo",func(c *gin.Context) {
		c.JSON(200,todos)
	})

	//查询
	r.GET("/todo/:index",func (c *gin.Context)  {
		index,_:=strconv.Atoi(c.Param("index"))
		if index>len(todos) {
			c.JSON(200,gin.H{"status":"索引过大,请重新查询"})
		}else {
			c.JSON(200,todos[index])
		}
		
	})
	
	r.Run(":8080")


}
