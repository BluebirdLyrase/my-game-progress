package playthrough

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"

	"my-game-progress/database"
	model_playthrough "my-game-progress/model/playthrough"
)

const filePath string = "/api/image/"

func Insert(c *gin.Context) {

	var form model_playthrough.PlaythroughInputParam
	var playthrough model_playthrough.Playthrough

	// Bind form values to struct and validate
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	defer file.Close()

	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file"})
		return
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

	playthrough.Difficulty = form.Difficulty
	playthrough.DateFinished = form.DateFinished
	playthrough.Remark = form.Remark
	playthrough.GameID = form.GameID
	playthrough.EnvironmentID = form.EnvironmentID
	playthrough.Screenshots[0] = filePath + fileID.Hex()

	// document := game
	_, err = collection.InsertOne(database.Context, playthrough)
	if err != nil {
		log.Fatal("Error inserting document:", err)
	}

	c.JSON(200, gin.H{"message": "Inserted record successfully"})
}
