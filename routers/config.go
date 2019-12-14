package routers

import (
	//"fmt"

	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistPageExtractor/data"
)
func ConfigGET(c *gin.Context) {
	data.GlobalConfig.CountUsedSID = len(data.GlobalWordListStorage)
	c.JSON(200, data.GlobalConfig)
}
func ConfigPUT(c *gin.Context) {
	var s data.ResponseStatus
	
	err := c.BindJSON(&data.GlobalConfig)
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	
	s = data.ResponseStatus{Code: 200, Text: "entity added"}
	c.JSON(200, s)
}
func ConfigPOST(c *gin.Context) {
	var s data.ResponseStatus
	var err error
	
	err = c.BindJSON(&data.GlobalConfig)
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	
	err = data.ReadGlobalWordlistFromRemote()
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "can't read global wordlist"}
		c.JSON(200, s)
		return
	}
	
	s = data.ResponseStatus{Code: 200, Text: "got global wordlist"}
	c.JSON(200, s)
}
