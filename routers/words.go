package routers

import (
	"strconv"
	"github.com/gin-gonic/gin"
	model "github.com/mazeForGit/WordlistPageExtractor/model"
)

func WordsGET(c *gin.Context) {
	
	sidString := c.Query("sid")
	if sidString == "" {
		c.String(400, "missing parameter: sid")
		return
	}
	sid, err := strconv.Atoi(sidString)
	if err != nil {
		c.String(400, "wrong parameter: sid")
		return
	}
	
	wl, err := model.GetWordListForSession(sid)
	if err != nil {
		c.String(400, "wrong parameter: sid")
		return
	}
	
	c.Header("Content-Type", "application/json")
	c.JSON(200, wl.Words)
}
