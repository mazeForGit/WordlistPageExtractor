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
		data.GlobalConfig.LastUsedSID++
		sid = data.GlobalConfig.LastUsedSID
		session.Set("sid", sid)
		session.Save()
	} else {
		sid = v.(int)
	}
		
	sData := data.GetSessionData(sid)
	sData.SessionStatus.Count++
	
	c.Header("Content-Type", "application/json")
	c.JSON(200, sData.SessionWords)
}
