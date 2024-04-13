package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const address = "0.0.0.0:8080"

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

var data = TodoPageData{
	PageTitle: "My TODO list",
	Todos: []Todo{
		{Title: "Task 1", Done: false},
		{Title: "Task 2", Done: false},
		{Title: "Task 3", Done: true},
	},
}

func main() {
	log.Println("main.go running")
	StartServer()
}

func StartServer() {
	r := gin.Default()
	r.LoadHTMLFiles("./index.html")

	r.GET("/", homeHandler)
	r.POST("/add", addTodo)
	r.POST("/status", changeStatus)

	log.Printf("listening at port: %v", address)

	r.Run(address)
}

func homeHandler(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", data)
}

func addTodo(ctx *gin.Context) {
	time.Sleep(1 * time.Second)
	var req Todo
	req.Title = ctx.PostForm("email")
	req.Done = false
	data.Todos = append(data.Todos, req)
	ctx.HTML(http.StatusOK, "list", data)
}

func changeStatus(ctx *gin.Context) {
	time.Sleep(1 * time.Second)
	title := strings.ToLower(ctx.Request.FormValue("task"))
	for i := range data.Todos {
		if strings.ToLower(data.Todos[i].Title) == title {
			data.Todos[i].Done = true
			break
		}
	}
	log.Println(title, data)
	ctx.HTML(http.StatusOK, "list", data)
}
