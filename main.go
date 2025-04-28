package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"my-game-progress/database"
	environment_handler "my-game-progress/handler/environment"
	game_handler "my-game-progress/handler/game"
	images_handler "my-game-progress/handler/images"
	playthrough_handler "my-game-progress/handler/playthrough"
	webpage_handler "my-game-progress/handler/webpage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	err = database.Init()
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
	environment := api.Group("environment")
	playthrough := api.Group("playthrough")

	r.GET("/", webpage_handler.IndexPage)
	r.GET("/my-progress-list", webpage_handler.GamePage)
	r.GET("/update-game", webpage_handler.EditPage)
	r.GET("/add-info", webpage_handler.AddInfoPage)
	r.GET("/add-playthrough", webpage_handler.AddPlaythroughPage)

	game.GET("list", game_handler.GameList)
	game.GET("selector", game_handler.GameSelector)
	game.GET("list-detail", game_handler.GameListDetail)
	game.POST("", game_handler.Insert)
	game.POST("multiple", game_handler.InsertMultiple)
	game.PUT("", game_handler.UpdateGame)

	image.POST("", images_handler.UploadImage)
	image.GET(":file_id", images_handler.GetImage)

	environment.POST("", environment_handler.Insert)
	environment.GET(":id", environment_handler.Get)
	environment.GET("selector", environment_handler.GeSelector)

	playthrough.POST("", playthrough_handler.Insert)

	r.Run(":" + "8080")
}

// UploadImage
