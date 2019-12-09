package routers

import (
	"strconv"
	
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/Wordlist/data"
)

func WordListGET(c *gin.Context) {
	c.JSON(200, data.GlobalWordList)
}
