package main

import (
	"fmt"
	"flag"
	"os"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/static"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	log "github.com/sirupsen/logrus"

	model "github.com/mazeForGit/WordlistPageExtractor/model"
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

	//
	// define and handle flags
	var flagServerName = flag.String("name", "wordListStorage", "server name")
	var flagServerPort = flag.String("port", "6001", "server port")
	var fileConfig = flag.String("frConfig", "./data/config.json", "file containing config")
	var fileWordList = flag.String("frWL", "./data/wordList.json", "file containing wordList")
	flag.Parse()
	
	//
	// handle flags
	if *fileConfig != "" {
		err := model.ReadConfigurationFromFile(*fileConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	if *fileWordList != "" {
		err := model.ReadWordListFromFile(*fileWordList)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	
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
//	router.GET("/voter", routers.Voter)
	router.GET("/health", routers.HealthGET)
	
	// global config
	router.GET("/config", routers.ConfigGET)
	router.POST("/config", routers.ConfigPOST)
	router.GET("/wordlist", routers.WordListGET)
	
	// session based
	router.GET("/status", routers.StatusGET)
	router.POST("/status", routers.StatusPOST)
	router.GET("/words", routers.WordsGET)
	
	log.Info("starting application on port " + getPort())
	router.Run(getPort())
}