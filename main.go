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
	Number  string `json:"number"`
	Content string `json:"content"`
	Done    string `json:"done"`
	Finish_time    string `json:"finish_time"`
}

var todos []TODO

func readTODO(){
	file,_:=os.OpenFile("TODO.txt",os.O_RDWR,0666)
	read,_:=io.ReadAll(file)
	readline:=strings.Split(string(read),"\n")
	lenth:=len(readline)/4
	var todo TODO
	for i:=0;i<lenth;i++{
		todo.Number,todo.Content,todo.Done,todo.Finish_time=readline[i*4],readline[i*4+1],readline[i*4+2],readline[i*4+3]
		todos=append(todos,todo)
	}
	file.Close()
}


func add(text TODO){
	file,_:=os.OpenFile("TODO.txt",os.O_APPEND,0666)
	file.WriteString(text.Number)
	file.WriteString("\n")
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
	readline[ind*4+1],readline[ind*4+2],readline[ind*4+3]=`已被删除`,`已被删除`,`已被删除`
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
	readline[ind*4+1],readline[ind*4+2],readline[ind*4+3]=text.Content,text.Done,text.Finish_time
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
		todo.Number=strconv.Itoa(len(todos)+1)
		todos = append(todos, todo)
		fmt.Println(todos) 
		add(todo)
		c.JSON(200,gin.H{"status":"成功增加TODO"})
	})

	//删除 TODO
	r.DELETE("/todo/:index",func (c *gin.Context)  {
		index,_:=strconv.Atoi(c.Param("index"))
		todos[index-1].Content=`已被删除`
		todos[index-1].Done=`已被删除`
		todos[index-1].Finish_time=`已被删除`
		del(index-1)
		c.JSON(200,gin.H{"status":"成功删除TODO"})
	})

	//更新 TODO
	r.PUT("/todo/:index",func(c *gin.Context){
		index,_:=strconv.Atoi(c.Param("index"))
		var todo TODO
		c.BindJSON(&todo)
		todos[index]=todo
		change(todo,index)
		c.JSON(200,gin.H{"status":"成功更新TODO"})
	})

	//列出 TODO
	r.GET("/todo",func(c *gin.Context) {
		page:=c.Query("page");	
		intpage,_:=strconv.Atoi(page)
		if (intpage-1)*20>len(todos){
			c.JSON(200,gin.H{"status":"页码过大，您没有这么多TODO"})
		}else if len(todos)==0 {
			c.JSON(200,gin.H{"status":"您还未添加TODO"})
		}else {
			c.JSON(200,gin.H{"status":"以下是该页的TODO内容："})
			if intpage*20>len(todos){
				c.JSON(200,todos[20*(intpage-1):])
			}else{
				c.JSON(200,todos[20*(intpage-1):20*intpage])
			}
			
		}
		
	})

	//查询 TODO
	r.GET("/todo/:index",func (c *gin.Context)  {
		index,_:=strconv.Atoi(c.Param("index"))
		if index>len(todos) {
			c.JSON(200,gin.H{"status":"序号过大,请重新查询"})
		}else {
			c.JSON(200,gin.H{"status":"以下是您所查询的的TODO："})
			c.JSON(200,todos[index-1])
		}
		
	})

	//清空 TODO
	r.DELETE("/todo/delete",func(c *gin.Context){
		file,_:=os.OpenFile("TODO.txt",os.O_RDWR,0666)
		file.Truncate(0)
		file.Seek(0,0)
		file.Close()
		c.JSON(200,gin.H{"status":"已成功清空TODO"})
		todos=todos[:0]
	})
	
	r.Run(":8080")


}
