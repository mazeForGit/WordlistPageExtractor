package routers

import (
	//"fmt"

	"github.com/gin-gonic/gin"
	model "github.com/mazeForGit/WordlistPageExtractor/model"
)
func ConfigGET(c *gin.Context) {
	c.JSON(200, model.GlobalConfig)
}
func ConfigPOST(c *gin.Context) {
	var s model.ResponseStatus
	var err error
	
	err = c.BindJSON(&model.GlobalConfig)
	if err != nil {
		s = model.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	
	err = model.ReadGlobalWordlistFromRemote()
	if err != nil {
		s = model.ResponseStatus{Code: 422, Text: "can't read global wordlist"}
		c.JSON(200, s)
		return
	}
	
	s = model.ResponseStatus{Code: 200, Text: "got global wordlist"}
	c.JSON(200, s)
}
