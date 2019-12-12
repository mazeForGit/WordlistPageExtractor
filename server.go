package main

import (
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/static"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	log "github.com/sirupsen/logrus"

	routers "github.com/mazeForGit/WordlistPageExtractor/routers"
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
	store := memstore.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session-id", store))
	
	router.LoadHTMLGlob("public/*.html")
	router.Use(static.Serve("/", static.LocalFile("./public", false)))
	router.NoRoute(routers.NotFoundError)
	router.GET("/content", routers.Index)
	router.GET("/extractor", routers.Extractor)
	router.GET("/voter", routers.Voter)
	router.GET("/health", routers.HealthGET)
	
	// global config
	router.GET("/config", routers.ConfigGET)
	router.POST("/config", routers.ConfigPOST)
	router.PUT("/config", routers.ConfigPUT)
	router.GET("/wordlist", routers.WordListGET)
	
	// session based
	router.GET("/status", routers.StatusGET)
	router.POST("/status", routers.StatusPOST)
	router.GET("/words", routers.WordsGET)
	
	log.Info("starting application on port " + getPort())
	router.Run(getPort())
}