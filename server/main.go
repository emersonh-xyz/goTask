package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Declare our routes
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskById)

	router.POST("/tasks", postTask)

	router.Run("localhost:8080")
}

type task struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

var tasks = []task{
	{ID: "1", Name: "Task 1", Status: "Pending", Description: "This is task 1"},
	{ID: "2", Name: "Task 2", Status: "In Progress", Description: "This is task 2"},
	{ID: "3", Name: "Task 3", Status: "Completed", Description: "This is task 3"},
}

// Tasks API
func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

func postTask(c *gin.Context) {
	var newTask task

	if err := c.BindJSON(&newTask); err != nil {
		return
	}

	tasks = append(tasks, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func getTaskById(c *gin.Context) {
	id := c.Param("id")

	for _, t := range tasks {
		if t.ID == id {
			c.IndentedJSON(http.StatusFound, t)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}
