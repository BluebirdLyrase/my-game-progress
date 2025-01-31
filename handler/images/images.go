package images_handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"

	"my-game-progress/database"
)

func UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()
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
	c.JSON(200, gin.H{
		"message":  "Image uploaded successfully",
		"filename": filename,
		"file_id":  fileID.Hex(),
	})
}

func GetImage(c *gin.Context) {
	fileID := c.Param("file_id")

	objID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	bucket, err := gridfs.NewBucket(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create GridFS bucket"})
		return
	}

	stream, err := bucket.OpenDownloadStream(objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer stream.Close()

	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Disposition", "inline")

	_, err = io.Copy(c.Writer, stream)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send image"})
	}
}
