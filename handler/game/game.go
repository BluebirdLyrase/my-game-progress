package game_handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"

	"my-game-progress/database"
	"my-game-progress/model/model_game"
	"my-game-progress/service"
)

const filePath string = "/api/image/"

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func GamePage(c *gin.Context) {
	games, err := service.GetGameList(bson.M{}, bson.M{"year": -1}, 30)
	if err != nil {
		log.Fatalf("Failed to get game list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	c.HTML(http.StatusOK, "my-progress-list.html", gin.H{
		"games": games,
	})
}

func EditPage(c *gin.Context) {
	games, err := service.GetGameFullDetail(bson.M{}, bson.M{"year": -1}, 30)
	if err != nil {
		log.Fatalf("Failed to get game list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	c.HTML(http.StatusOK, "game-update.html", gin.H{
		"games": games,
	})
}

func GameList(c *gin.Context) {
	title := c.Query("title")
	games, err := service.GetGameList(bson.M{
		"title": bson.M{"$regex": title, "$options": "i"},
	}, bson.M{"year": -1}, 0)
	if err != nil {
		log.Fatalf("Failed to get game list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"games": games,
	// })
	c.HTML(http.StatusOK, "game-gallery.html", gin.H{
		"games": games,
	})
}

func GameListDetail(c *gin.Context) {
	title := c.Query("title")
	games, err := service.GetGameFullDetail(bson.M{
		"title": bson.M{"$regex": title, "$options": "i"},
	}, bson.M{"year": -1}, 0)
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
	// c.HTML(http.StatusOK, "game-gallery.html", gin.H{
	// 	"games": games,
	// })
}

func Insert(c *gin.Context) {

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

func InsertMultiple(c *gin.Context) {
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

func UpdateGame(c *gin.Context) {

	file, header, err := c.Request.FormFile("image")
	var gameCover string

	if file != nil {
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
		gameCover = filePath + fileID.Hex()
	}

	gameID := c.Request.FormValue("game_id")
	gameTitle := c.Request.FormValue("game_title")
	gameYear := c.Request.FormValue("game_year")
	slug := c.Request.FormValue("slug")
	gameRemark := c.Request.FormValue("game_remark")
	numYear, err := strconv.Atoi(gameYear)
	if err != nil {
		fmt.Println("Error converting string to int numYear:", err)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"year":   numYear,
			"title":  gameTitle,
			"remark": gameRemark,
			"slug":   slug,
		},
	}

	if gameCover != "" {
		update["$set"].(bson.M)["gameImage.cover"] = gameCover
	}

	hexStr := gameID[10 : len(gameID)-2]
	fmt.Println(hexStr)
	id, err := primitive.ObjectIDFromHex(hexStr)
	fmt.Println(id)

	filter := bson.M{"_id": id}

	collection := database.DB.Collection("game")
	if err != nil {
		log.Fatal(err)
	}

	// Perform the update
	result, err := collection.UpdateOne(database.Context, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	// Return the response
	fmt.Printf("ModifiedCount: %+v\n", result.ModifiedCount)
	if gameCover != "" {
		c.HTML(http.StatusOK, "upload-game-images.html", bson.M{"GameCover": gameCover})
	}
	originalCover := c.Request.FormValue("original_cover")
	c.HTML(http.StatusOK, "upload-game-images.html", bson.M{"GameCover": originalCover})
}
