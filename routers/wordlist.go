package routers

import (
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistPageExtractor/data"
)

func WordListGET(c *gin.Context) {
	c.JSON(200, data.GlobalWordList)
}
