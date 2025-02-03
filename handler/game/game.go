package game_handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"

	"my-game-progress/database"
	"my-game-progress/model/model_game"
	"my-game-progress/service"
)

const filePath string = "/api/image/"

func GamePage(c *gin.Context) {
	games, err := service.GetGameList()
	if err != nil {
		log.Fatalf("Failed to get game list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"games": games,
	})
}

func GameList(c *gin.Context) {
	games, err := service.GetGameList()
	if err != nil {
		log.Fatalf("Failed to get game list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"games": games,
	})
}

func InsertGame(c *gin.Context) {

	file, header, err := c.Request.FormFile("image")
	defer file.Close()
	detail := c.Request.FormValue("detail")

	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file"})
		return
	}

	var game model_game.Game
	err = json.Unmarshal([]byte(detail), &game)
	if err != nil {
		log.Fatal(err)
	}

	filename := header.Filename

	bucket, err := gridfs.NewBucket(database.DB)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create GridFS bucket"})
		return
	}

	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to open upload stream"})
		return
	}
	defer uploadStream.Close()

	_, err = io.Copy(uploadStream, file)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload image"})
		return
	}
	fileID := uploadStream.FileID.(primitive.ObjectID)

	collection := database.DB.Collection("game")

	game.GameImage.Cover = filePath + fileID.Hex()

	// document := game
	_, err = collection.InsertOne(database.Context, game)
	if err != nil {
		log.Fatal("Error inserting document:", err)
	}

	c.JSON(200, gin.H{"message": "Inserted record successfully"})
}

func InsertMultipleGame(c *gin.Context) {
	var jsonData []model_game.Game
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON data"})
		return
	}

	collection := database.DB.Collection("game")

	var documents []interface{}
	for _, game := range jsonData {
		documents = append(documents, game)
	}

	_, err := collection.InsertMany(database.Context, documents)
	if err != nil {
		log.Println("Error inserting documents:", err)
		c.JSON(500, gin.H{"error": "Failed to insert records"})
		return
	}

	c.JSON(200, gin.H{"message": "Inserted records successfully"})
}
