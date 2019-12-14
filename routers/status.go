package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	
	data "github.com/mazeForGit/WordlistPageExtractor/data"
)

func StatusGET(c *gin.Context) {
	
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
	
	c.JSON(200, sData.Session)
}
func StatusPOST(c *gin.Context) {

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
	
	var s data.ResponseStatus
	err := c.BindJSON(&sData.Session)
	
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "wrong request"}
		c.JSON(200, s)
		return
	}
	
	sData.Session.RequestExecution = true
	sData.Words = data.CopyWords(data.GlobalWordList.Words)
	//data.SetSessionData(sData)
	
	go data.ExecuteLongRunningTaskOnRequest(sid)

	s = data.ResponseStatus{Code: 200, Text: "start execution"}
	c.JSON(200, s)
}
