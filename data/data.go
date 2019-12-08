package data

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"runtime"
	"time"
)

var GlobalConfig Config
var GlobalWordList WordList

func PrintMemUsage() {
	// from: https://golangcode.com/print-the-current-memory-usage/
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    // for info on each, see: https://golang.org/pkg/runtime/#MemStats
    fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
    fmt.Printf(", TotalAlloc = %v MiB", bToMb(m.TotalAlloc))
    fmt.Printf(", Sys = %v MiB", bToMb(m.Sys))
    fmt.Printf(", NumGC = %v\n", m.NumGC)
}
func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
func ExecuteLongRunningTaskOnRequest() {
    for {
		PrintMemUsage()
		time.Sleep(2 * time.Second)
		if GlobalConfig.RequestExecution {
			readSite()
		}
	}
}

func readSite() {
    var s Status
	var err error
	var resp *http.Response
	var body []byte
	var requestUrl string = ""
	
    // do something with the result
    requestUrl = GlobalConfig.WordListUrl + "/words?testOnly=true&format=json"
    fmt.Println("connect to wordlist and get words with tests = " + requestUrl)
    resp, err = http.Get(requestUrl)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    defer resp.Body.Close()
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    //fmt.Println("body = " + string(body))

    s.Text = ""
    json.Unmarshal(body, &GlobalWordList.Words)

    fmt.Println(GlobalWordList.Words)
    return
}