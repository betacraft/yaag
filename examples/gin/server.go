package main

import (
	"net/http"

	yaag_gin "github.com/CoryARamirez/yaag/gin"
	"github.com/CoryARamirez/yaag/yaag"
	"github.com/gin-gonic/gin"
)

func main() {
	yaag.Init(&yaag.Config{On: true, DocTitle: "Gin", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})
	r := gin.Default()
	r.Use(yaag_gin.Document())
	r.GET("/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": "Hello World!"})
	})
	r.GET("/plain", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/xml", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"result": "Hello World!"})
	})
	r.GET("/complex", func(c *gin.Context) {
		value := c.Query("key")
		c.JSON(http.StatusOK, gin.H{"value": value})
	})
	r.Run(":8080")
}
