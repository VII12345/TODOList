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
	Done    string `json:"done"`
	Finish_time    string `json:"finish_time"`
}

var todos []TODO

func readTODO(){
	file,_:=os.OpenFile("TODO.txt",os.O_RDWR,0666)
	read,_:=io.ReadAll(file)
	readline:=strings.Split(string(read),"\n")
	lenth:=len(readline)/3
	var todo TODO
	for i:=0;i<lenth;i++{
		todo.Content,todo.Done,todo.Finish_time=readline[i*3],readline[i*3+1],readline[i*3+2]
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
	file.WriteString(text.Finish_time)
	file.WriteString("\n")
	file.Close()
}

func del(ind int){
	file,_:=os.OpenFile("TODO.txt",os.O_RDWR,0666)
	read,_:=io.ReadAll(file)
	readline:=strings.Split(string(read),"\n")
	readline[ind*3],readline[ind*3+1],readline[ind*3+2]="has been deleted","has been deleted","has been deleted"
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
	readline[ind*3],readline[ind*3+1],readline[ind*3+2]=text.Content,text.Done,text.Finish_time
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


	//增加 TODO
	r.POST("/todo",func(c *gin.Context){
		var todo TODO
		c.BindJSON(&todo)
		todos = append(todos, todo)
		fmt.Println(todos) 
		add(todo)
		c.JSON(200,gin.H{"status":"ok"})
	})

	//删除 TODO
	r.DELETE("/todo/:index",func (c *gin.Context)  {
		index,_:=strconv.Atoi(c.Param("index"))
		todos[index].Content=`has been deleted`
		todos[index].Done=`has been deleted`
		todos[index].Finish_time=`has been deleted`
		del(index)
		c.JSON(200,gin.H{"status":"ok"})
	})

	//更新 TODO
	r.POST("/todo/:index",func(c *gin.Context){
		index,_:=strconv.Atoi(c.Param("index"))
		var todo TODO
		c.BindJSON(&todo)
		todos[index]=todo
		change(todo,index)
		c.JSON(200,gin.H{"status":"ok"})
	})

	//列出 TODO
	r.GET("/todo",func(c *gin.Context) {
		c.JSON(200,todos)
	})

	//查询 TODO
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
