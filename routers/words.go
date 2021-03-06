package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	
	data "github.com/mazeForGit/WordlistPageExtractor/data"
)

func WordsGET(c *gin.Context) {
	
	session := sessions.Default(c)
	
	var sid int
	v := session.Get("sid")
	if v == nil {
		sid = data.GetNewSessionID()
		session.Set("sid", sid)
		session.Save()
	} else {
		sid = v.(int)
	}
		
	sData := data.GetWordListForSession(sid)
	sData.Session.Count++
	
	c.Header("Content-Type", "application/json")
	c.JSON(200, sData.Words)
}
