package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"my-game-progress/database"
	game_handler "my-game-progress/handler/game"
	images_handler "my-game-progress/handler/images"
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
	err := database.Init()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
		panic(err)
	}
	r := gin.New()
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./favicon.ico")

	api := r.Group("api")
	image := api.Group("image")
	game := api.Group("game")

	r.GET("/", game_handler.IndexPage)
	r.GET("/my-progress-list", game_handler.GamePage)
	r.GET("/update-game", game_handler.EditPage)

	game.GET("list", game_handler.GameList)
	game.GET("list-detail", game_handler.GameListDetail)
	game.POST("", game_handler.Insert)
	game.POST("multiple", game_handler.InsertMultiple)
	game.PUT("", game_handler.UpdateGame)

	image.POST("", images_handler.UploadImage)
	image.GET(":file_id", images_handler.GetImage)
	r.Run(":" + "8080")
}

// UploadImage
