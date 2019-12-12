package routers

import (
	"fmt"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	
	data "github.com/mazeForGit/WordlistPageExtractor/data"
)

func StatusGET(c *gin.Context) {
	
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
	
	fmt.Println("StatusGET")
	fmt.Println(sData.SessionStatus)
	
	c.JSON(200, sData.SessionStatus)
}
func StatusPOST(c *gin.Context) {

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
	
	fmt.Println("StatusPOST")
	fmt.Println(sData.SessionStatus)
	
	var s data.ResponseStatus
	err := c.BindJSON(&sData.SessionStatus)
	if err != nil {
		s = data.ResponseStatus{Code: 422, Text: "wrong request"}
		c.JSON(200, s)
		return
	}
	
	sData.SessionStatus.RequestExecution = true
	sWords := data.GlobalWordList.Words
	sData.SessionWords = sWords
	go data.ExecuteLongRunningTaskOnRequest(sid)
	s = data.ResponseStatus{Code: 200, Text: "start execution"}
	c.JSON(200, s)
}
