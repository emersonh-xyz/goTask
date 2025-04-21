package main

import (
	"net/http"

	"encoding/csv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "net/http"
	"strconv"
	// "strings"
)

func main() {
	router := gin.Default()

	// Enable CORS
	router.Use(cors.Default())

	// Routes
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskById)
	router.PUT("/tasks/:id", editTask)
	router.PUT("/tasks/complete/:id", markAsComplete)
	router.POST("/tasks", postTask)
	router.GET("/tasks/export", exportTasksToCSV)

	// Start server
	router.Run("localhost:8080")
}

type task struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Description  string `json:"description"`
	TimeEstimate int    `json:"timeEstimate"`
	DueDate      string `json:"dueDate"`
	IsComplete   bool   `json:"isComplete"`
}

// Mock task data
var tasks = []task{
	{
		ID:           "1",
		Name:         "Task 1",
		Status:       "Pending",
		Description:  "This is task 1",
		TimeEstimate: 5,
		DueDate:      "2023-12-01",
		IsComplete:   false,
	},
	{
		ID:           "2",
		Name:         "Task 2",
		Status:       "In Progress",
		Description:  "This is task 2",
		TimeEstimate: 3,
		DueDate:      "2023-12-05",
		IsComplete:   false,
	},
	{
		ID:           "3",
		Name:         "Task 3",
		Status:       "Completed",
		Description:  "This is task 3",
		TimeEstimate: 2,
		DueDate:      "2023-11-30",
		IsComplete:   true,
	},
}

// Tasks API Routes
func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}
func postTask(c *gin.Context) {
	var newTask task

	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task format"})
		return
	}

	// Assign a new ID (incremented from the last task)
	newTask.ID = generateNextID()

	tasks = append(tasks, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}
func markAsComplete(c *gin.Context) {
	id := c.Param("id")

	for i, t := range tasks {
		if t.ID == id {
			// Toggle the isComplete field
			tasks[i].IsComplete = !tasks[i].IsComplete
			c.IndentedJSON(http.StatusOK, tasks[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

func editTask(c *gin.Context) {
	id := c.Param("id")
	var updatedTask task

	if err := c.BindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task format"})
		return
	}

	for i, t := range tasks {
		if t.ID == id {
			// Keep original ID intact, just update the rest
			tasks[i] = updatedTask
			updatedTask.ID = id
			c.IndentedJSON(http.StatusOK, updatedTask)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

func generateNextID() string {
	if len(tasks) == 0 {
		return "1"
	}
	lastID := tasks[len(tasks)-1].ID
	// Assuming IDs are numeric strings
	idNum, err := strconv.Atoi(lastID)
	if err != nil {
		return "1" // fallback if conversion fails
	}
	return strconv.Itoa(idNum + 1)
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
func exportTasksToCSV(c *gin.Context) {
	// Set headers to indicate it's a CSV file
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=tasks.csv")

	// Create a new CSV writer
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write the header row to the CSV
	writer.Write([]string{"ID", "Name", "Status", "Description", "Time Estimate", "Due Date", "Is Complete"})

	// Write task data to CSV
	for _, t := range tasks {
		writer.Write([]string{
			t.ID,
			t.Name,
			t.Status,
			t.Description,
			strconv.Itoa(t.TimeEstimate),
			t.DueDate,
			strconv.FormatBool(t.IsComplete),
		})
	}

	// No need to return a JSON response since the file is directly written to the response
}
