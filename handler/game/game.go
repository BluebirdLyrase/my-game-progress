package game_handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"my-game-progress/database"
	model_game "my-game-progress/model/game"
	"my-game-progress/service"
)

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
		return
	}

	c.HTML(http.StatusOK, "game-gallery.html", gin.H{
		"games": games,
	})
}

func GameSelector(c *gin.Context) {
	games, err := service.GetGameList(nil, bson.M{"title": -1}, 0)
	if err != nil {
		log.Fatalf("Failed to get game list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "game-selector.html", gin.H{
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
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"games": games,
	})
}

func Insert(c *gin.Context) {

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file"})
		return
	}

	detail := c.Request.FormValue("detail")
	var game model_game.Game
	err = json.Unmarshal([]byte(detail), &game)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, err)
		return

	}

	fullPath, err := service.UploadImage(file, header)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, err)
		return
	}
	collection := database.DB.Collection("game")

	game.GameImage.Cover = fullPath

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
		fullPath, err := service.UploadImage(file, header)
		if err != nil {
			log.Fatal(err)
			c.JSON(400, err)
			return
		}
		gameCover = fullPath
	}

	gameID := c.Request.FormValue("game_id")
	gameTitle := c.Request.FormValue("game_title")
	gameYear := c.Request.FormValue("game_year")
	slug := c.Request.FormValue("slug")
	gameRemark := c.Request.FormValue("game_remark")
	numYear, err := strconv.Atoi(gameYear)
	if err != nil {
		fmt.Println("Error converting string to int numYear:", err)
		c.JSON(400, err)
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
	id, err := primitive.ObjectIDFromHex(hexStr)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, err)
		return
	}

	filter := bson.M{"_id": id}

	collection := database.DB.Collection("game")

	result, err := collection.UpdateOne(database.Context, filter, update)
	if err != nil {
		log.Fatal(err)
		c.JSON(400, err)
		return
	}

	fmt.Printf("ModifiedCount: %+v\n", result.ModifiedCount)
	if gameCover != "" {
		c.HTML(http.StatusOK, "upload-game-images.html", bson.M{"GameCover": gameCover})
	}
	originalCover := c.Request.FormValue("original_cover")
	c.HTML(http.StatusOK, "upload-game-images.html", bson.M{"GameCover": originalCover})
}
