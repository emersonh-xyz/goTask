package main

import (
	"context"
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "net/http"
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var client *mongo.Client

func initMongo() {
	// Initialize MongoDB client
	var err error

	// Probably should set the URI in an env but oh well
	client, err = mongo.Connect(nil, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}

	// Ping the database to ensure the connection is established
	if err := client.Ping(nil, readpref.Primary()); err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}
	fmt.Println("Connected to MongoDB")

	// Ensure the connection is closed when the application exits
}

func main() {

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Establish our connection to mongo
	initMongo()

	// Ensure the connection is closed when the application exits
	defer func() {
		if err := client.Disconnect(nil); err != nil {
			fmt.Println("Error disconnecting from MongoDB:", err)
		}
	}()

	// Load dummy data
	// loadDummyData()

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
	// Connect to the "tasks" collection in the "gotask" database
	collection := client.Database("gotask").Collection("tasks")

	// Fetch all tasks from the collection
	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	defer cursor.Close(c)

	// Parse the tasks into a slice
	var tasks []task
	if err := cursor.All(c, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse tasks"})
		return
	}

	// Return the tasks as JSON
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

func loadDummyData() {
	// Connect to the "tasks" collection in the "gotask" database
	collection := client.Database("gotask").Collection("tasks")

	// Dummy task data
	dummyTasks := []interface{}{
		task{
			ID:           "1",
			Name:         "Task 1",
			Status:       "Pending",
			Description:  "This is task 1",
			TimeEstimate: 5,
			DueDate:      "2023-12-01",
			IsComplete:   false,
		},
		task{
			ID:           "2",
			Name:         "Task 2",
			Status:       "In Progress",
			Description:  "This is task 2",
			TimeEstimate: 3,
			DueDate:      "2023-12-05",
			IsComplete:   false,
		},
		task{
			ID:           "3",
			Name:         "Task 3",
			Status:       "Completed",
			Description:  "This is task 3",
			TimeEstimate: 2,
			DueDate:      "2023-11-30",
			IsComplete:   true,
		},
	}

	// Insert dummy tasks into the collection
	result, err := collection.InsertMany(context.TODO(), dummyTasks)
	if err != nil {
		log.Fatal("Failed to insert dummy tasks:", err)
	}

	fmt.Printf("Inserted %d tasks: %v\n", len(result.InsertedIDs), result.InsertedIDs)
}
