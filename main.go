package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"my-game-progress/database"
	environment_handler "my-game-progress/handler/environment"
	game_handler "my-game-progress/handler/game"
	images_handler "my-game-progress/handler/images"
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

	environment.POST("", environment_handler.Insert)
	environment.GET(":id", environment_handler.Get)
	environment.GET("list", environment_handler.GetAll)

	r.Run(":" + "8080")
}

// UploadImage
