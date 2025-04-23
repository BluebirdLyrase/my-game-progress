package environment

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"my-game-progress/database"
	model_environment "my-game-progress/model/environment"
	"my-game-progress/service"
)

func Insert(c *gin.Context) {
	var jsonData model_environment.Environment
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON data : " + err.Error()})
		return
	}

	collection := database.DB.Collection("environment")

	_, err := collection.InsertOne(database.Context, jsonData)
	if err != nil {
		log.Println("Error inserting documents:", err)
		c.JSON(500, gin.H{"error": "Failed to insert records"})
		return
	}

	c.JSON(200, gin.H{"message": "Inserted records successfully"})
}

func Get(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Errorf("invalid ObjectID: %w", err)
		return
	}
	collection := database.DB.Collection("environment")

	var environment model_environment.Environment
	err = collection.FindOne(database.Context, bson.M{"_id": objectID}).Decode(&environment)
	if err != nil {
		fmt.Errorf("environment not found: %w", err)
		return
	}

	c.JSON(http.StatusOK, environment)
}

func GetAll(c *gin.Context) {
	environment, err := service.GetEnvironmentsList()

	if err != nil {
		log.Fatalf("Failed to get game list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, environment)
}
