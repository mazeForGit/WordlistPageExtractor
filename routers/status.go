package routers

import (
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
	model "github.com/mazeForGit/WordlistPageExtractor/model"
)

func StatusGET(c *gin.Context) {

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
	
	wl.Session.Count++
	c.Header("Content-Type", "application/json")
	c.JSON(200, wl.Session)
}
//
//
//
func StatusPOST(c *gin.Context) {
fmt.Println(". StatusPOST")
	
	var session model.Session
	err := c.BindJSON(&session)
	if err != nil {
		c.String(422, "wrong request")
		return
	}
	
	wl, err := model.NewSession(session)
	
fmt.Println("new session: ", wl.Session)
	if err != nil {
		c.String(422, err.Error())
		return
	}

	// create job
	work := model.Job{Wordlist: wl}
	
	// push the work onto the queue
	// model.JobChannel <- work
	if !tryEnqueue(work) {
		c.String(503, "Maximum capacity reached. Try later.")
		return
	}
	
	c.Header("Content-Type", "application/json")
	c.JSON(200, wl.Session)
}
//
// tryEnqueue tries to enqueue a job to the given job channel. 
// Returns true if the operation was successful, 
// and false if enqueuing would not have been
// possible without blocking. 
// Job is not enqueued in the latter case.
//
func tryEnqueue(work model.Job) bool {
    select {
    case model.JobChannel <- work:
        return true
    default:
        return false
    }
}
