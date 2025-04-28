package webpage

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"my-game-progress/service"
)

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

func AddInfoPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add-info.html", nil)
}

func AddPlaythroughPage(c *gin.Context) {
	c.HTML(http.StatusOK, "add-playthrough.html", nil)
}
