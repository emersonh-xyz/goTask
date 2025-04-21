package main

import (
	"context"
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"fmt"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	router.DELETE("/tasks/:id", deleteTask) // Add this line

	// Start server
	router.Run("localhost:8080")
}

type Task struct {
	ID           bson.ObjectID `bson:"_id,omitempty"    json:"id"`
	Name         string        `bson:"name"             json:"name"`
	Status       string        `bson:"status"           json:"status"`
	Description  string        `bson:"description"      json:"description"`
	TimeEstimate int           `bson:"timeEstimate"     json:"timeEstimate"`
	DueDate      string        `bson:"dueDate"          json:"dueDate"`
	IsComplete   bool          `bson:"isComplete"       json:"isComplete"`
}

func tasksColl() *mongo.Collection {
	return client.Database("gotask").Collection("tasks")
}

// Tasks API Routes
func getTasks(c *gin.Context) {
	ctx := context.TODO()
	cur, err := tasksColl().Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	var all []Task
	if err := cur.All(ctx, &all); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, all)
}

func postTask(c *gin.Context) {
	var t Task
	if err := c.BindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t.ID = bson.NewObjectID()
	_, err := tasksColl().InsertOne(context.TODO(), t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func editTask(c *gin.Context) {
	hexID := c.Param("id")
	oid, err := bson.ObjectIDFromHex(hexID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var upd Task
	if err := c.BindJSON(&upd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// build a $set document (you can omit fields you don't want to change)
	set := bson.M{
		"name":         upd.Name,
		"status":       upd.Status,
		"description":  upd.Description,
		"timeEstimate": upd.TimeEstimate,
		"dueDate":      upd.DueDate,
		"isComplete":   upd.IsComplete,
	}

	_, err = tasksColl().
		UpdateByID(context.TODO(), oid, bson.M{"$set": set})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return the new document
	upd.ID = oid
	c.JSON(http.StatusOK, upd)
}

func getTaskById(c *gin.Context) {
	hexID := c.Param("id")
	oid, err := bson.ObjectIDFromHex(hexID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var t Task
	err = tasksColl().
		FindOne(context.TODO(), bson.M{"_id": oid}).
		Decode(&t)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func markAsComplete(c *gin.Context) {
	hexID := c.Param("id")
	oid, err := bson.ObjectIDFromHex(hexID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// fetch current value
	var t Task
	if err := tasksColl().
		FindOne(context.TODO(), bson.M{"_id": oid}).
		Decode(&t); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// flip it
	newVal := !t.IsComplete
	_, err = tasksColl().
		UpdateByID(context.TODO(), oid, bson.M{"$set": bson.M{"isComplete": newVal}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	t.IsComplete = newVal
	c.JSON(http.StatusOK, t)
}

func deleteTask(c *gin.Context) {
	hexID := c.Param("id")
	oid, err := bson.ObjectIDFromHex(hexID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	ctx := context.TODO() // Consider using c.Request.Context() with timeout
	result, err := tasksColl().DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"}) // Or use http.StatusNoContent
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

	// Query all tasks from the collection
	cur, err := tasksColl().Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(context.TODO())

	// Iterate over the cursor and write each task to the CSV
	for cur.Next(context.TODO()) {
		var t Task
		if err := cur.Decode(&t); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		writer.Write([]string{
			t.ID.Hex(),
			t.Name,
			t.Status,
			t.Description,
			strconv.Itoa(t.TimeEstimate),
			t.DueDate,
			strconv.FormatBool(t.IsComplete),
		})
	}

	if err := cur.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
