package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"my-game-progress/database"
	"my-game-progress/service"
)

type Package struct {
	PackageName     string   `json:"packageName"`
	Price           string   `json:"price"`
	Time            string   `json:"time"`
	Icon            string   `json:"icon"`
	Cap             string   `json:"cap"`
	BtnMsg          string   `json:"btnMsg"`
	IsSelect        bool     `json:"isSelect"`
	Ussd            string   `json:"ussd"`
	Link            string   `json:"link"`
	Remark          string   `json:"remark"`
	ConditionList   []string `json:"conditionList"`
	ConditionLength int
	PackageIndex    int
}

func main() {
	// var basedPath string = ""
	//var basedPath string = "notify"
	err := database.Init()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
		panic(err)
	}
	r := gin.New()

	r.LoadHTMLGlob("templates/**/*")
	api := r.Group("api")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
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
	})

	api.GET("/game", func(c *gin.Context) {
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
	})

	r.Run(":" + "4200")
}
