package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"my-game-progress/database"
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
		log.Print("Hello")
		c.HTML(http.StatusOK, "index.html", nil)
	})

	api.GET("/game", func(c *gin.Context) {
		gameCollection := database.DB.Collection("Game")

		games, err := database.Fetch(gameCollection)
		if err != nil {
			log.Fatalf("Failed to fetch games: %v", err)
		}

		res := gin.H{
			"games": games,
		}

		c.JSON(http.StatusOK, res)
	})

	// r.GET("/money", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "top-up-redirect.html", nil)
	// })
	// n.GET("/internet", func(c *gin.Context) {

	// 	file, err := os.Open("static/json/package-list.json")
	// 	if err != nil {
	// 		fmt.Println("Error opening file:", err)
	// 		return
	// 	}
	// 	defer file.Close()

	// 	var data map[string][]Package
	// 	decoder := json.NewDecoder(file)
	// 	if err := decoder.Decode(&data); err != nil {
	// 		fmt.Println("Error decoding JSON:", err)
	// 		return
	// 	}

	// 	// Extract the package list
	// 	packages := data["packageList"]
	// 	for i := range packages {
	// 		packages[i].ConditionLength = len(packages[i].ConditionList)
	// 		packages[i].PackageIndex = i
	// 	}
	// 	c.HTML(http.StatusOK, "package-offer.html", gin.H{"packages": packages})
	// })
	r.Run(":" + "4200")
}
