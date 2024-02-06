package routers

import (
	"github.com/gin-gonic/gin"
	
	model "github.com/mazeForGit/WordlistPageExtractor/model"
)

func WordListGET(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(200, model.GlobalWordList.Words)
}
