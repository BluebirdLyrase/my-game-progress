package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	r := gin.New()

	r.LoadHTMLGlob("templates/**/*")
	n := r.Group("notify")
	n.Static("/static", "./static")
	n.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	n.GET("/money", func(c *gin.Context) {
		c.HTML(http.StatusOK, "top-up-redirect.html", nil)
	})
	n.GET("/internet", func(c *gin.Context) {

		file, err := os.Open("static/json/package-list.json")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		var data map[string][]Package
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&data); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Extract the package list
		packages := data["packageList"]
		for i := range packages {
			packages[i].ConditionLength = len(packages[i].ConditionList)
			packages[i].PackageIndex = i
		}
		c.HTML(http.StatusOK, "package-offer.html", gin.H{"packages": packages})
	})
	r.Run(":" + "4200")
}
