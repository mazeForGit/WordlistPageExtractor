package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	
	model "github.com/mazeForGit/WordlistPageExtractor/model"
)

func StatusGET(c *gin.Context) {
	
	session := sessions.Default(c)
	
	var sid int
	v := session.Get("sid")
	if v == nil {
		sid = model.GetNewSessionID()
		session.Set("sid", sid)
		session.Save()
	} else {
		sid = v.(int)
	}
	
	sData := model.GetWordListForSession(sid)
	sData.Session.Count++
	
	c.Header("Content-Type", "application/json")
	c.JSON(200, sData.Session)
}
func StatusPOST(c *gin.Context) {

	c.Header("Content-Type", "application/json")
	session := sessions.Default(c)
	
	var sid int
	v := session.Get("sid")
	if v == nil {
		sid = model.GetNewSessionID()
		session.Set("sid", sid)
		session.Save()
	} else {
		sid = v.(int)
	}
		
	sData := model.GetWordListForSession(sid)
	sData.Session.Count++
	
	var s model.ResponseStatus
	err := c.BindJSON(&sData.Session)
	
	if err != nil {
		s = model.ResponseStatus{Code: 422, Text: "wrong request"}
		c.JSON(200, s)
		return
	}
	
	sData.Session.RequestExecution = true
	sData.Words = model.CopyWords(model.GlobalWordList.Words)
	//model.SetSessionData(sData)
	
	go model.ExecuteLongRunningTaskOnRequest(sid)

	s = model.ResponseStatus{Code: 200, Text: "start execution"}
	c.JSON(200, s)
}
