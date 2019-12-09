package main

import (
	"os"

	routers "github.com/Wordlist/WordlistPageExtractor/routers"
	data "github.com/Wordlist/WordlistPageExtractor/data"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/static"
	log "github.com/sirupsen/logrus"
)

func getPort() string {
	p := os.Getenv("HTTP_PLATFORM_PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	router := gin.Default()
	router.RedirectTrailingSlash = false

	router.LoadHTMLGlob("public/*.html")
	router.Use(static.Serve("/", static.LocalFile("./public", false)))
	router.GET("/index", routers.Index)
	router.NoRoute(routers.NotFoundError)
	router.GET("/health", routers.HealthGET)
	
	router.GET("/config", routers.ConfigGET)
	router.POST("/config", routers.ConfigPOST)
	router.PUT("/config", routers.ConfigPUT)
	
	log.Info("starting background process")
	go data.ExecuteLongRunningTaskOnRequest()
	
	log.Info("starting application on port " + getPort())
	router.Run(getPort())
}