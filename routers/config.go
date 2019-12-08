package routers

import (
	//"strconv"
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistPageExtractor/data"
)
func ConfigGET(c *gin.Context) {
	c.JSON(200, data.GlobalConfig)
}
func ConfigPUT(c *gin.Context) {
	var s data.Status
	
	err := c.BindJSON(&data.GlobalConfig)
	if err != nil {
		s = data.Status{Code: 422, Text: "unprocessable entity"}
		c.JSON(422, s)
		return
	}
	
	s = data.Status{Code: 200, Text: "entity added"}
	c.JSON(200, s)
}
func ConfigPOST(c *gin.Context) {
	var s data.Status
	var vars map[string][]string
	vars = c.Request.URL.Query()
	var execution string = ""

	if _, ok := vars["execution"]; ok {
		execution = c.Request.URL.Query().Get("execution")
	} 
	
	//fmt.Println("execution = " + execution)
	if execution == "true" {
		
		if (data.GlobalConfig.WordListUrl == "") {
			s = data.Status{Code: 422, Text: "missing data"}
			c.JSON(200, s)
			return
		}
		s = data.Status{Code: 200, Text: "start execution"}
		c.JSON(200, s)
		
	} else if execution == "false" {
		data.GlobalConfig.RequestExecution = false
		
		s = data.Status{Code: 200, Text: "stop execution"}
		c.JSON(200, s)
	} else {
		s = data.Status{Code: 422, Text: "unknown request"}
		c.JSON(422, s)
	}
}
