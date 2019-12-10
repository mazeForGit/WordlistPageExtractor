package routers

import (
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/WordlistPageExtractor/data"
)

func WordListGET(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(200, data.GlobalWordList.Words)
}
