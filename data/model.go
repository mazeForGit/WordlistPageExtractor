package data

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"sort"
	
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var GlobalConfig Config
var GlobalWordList WordList
var GlobalSessionData []SessionData

type SorterWordByOccurance []Word

func (a SorterWordByOccurance) Len() int           { return len(a) }
func (a SorterWordByOccurance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SorterWordByOccurance) Less(i, j int) bool { return a[i].Occurance > a[j].Occurance }

func GetSessionData(sid int) *SessionData {
	//fmt.Println("GetSessionData sid=" + strconv.Itoa(sid))
	
	for i := 0; i < len(GlobalSessionData); i++ {
		if GlobalSessionData[i].SessionID == sid {
			return &GlobalSessionData[i]
		}
	}
	
	var sData = SessionData {
		SessionID: sid,
		SessionStatus: Status {
			Count: 0,
			RequestExecution: false,
			ExecutionStarted: false,
			ExecutionFinished: false,
			PageToScan: "",
			DomainsAllowed: "",
			NumberLinksFound: 0,
			NumberLinksVisited: 0,
			WordsScanned: 0,
		},
		SessionWords: nil,
	}
	GlobalSessionData = append(GlobalSessionData, sData)
    return GetSessionData(sid)
}

func ReadGlobalWordlist() error {
	fmt.Println("readGlobalWordlist")
	fmt.Println("have GlobalWordlist.Words = " + strconv.Itoa(len(GlobalWordList.Words)))
	
    var err error
	var resp *http.Response
	var body []byte
	var requestUrl string = ""
	
    requestUrl = GlobalConfig.WordListUrl + "/words?testOnly=true&format=json"
    fmt.Println("connect to wordlist and get words with tests = " + requestUrl)
    resp, err = http.Get(requestUrl)
    if err != nil {
        return err
    }

    defer resp.Body.Close()
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    json.Unmarshal(body, &GlobalWordList.Words)
    fmt.Println("got GlobalWordList.Words = " + strconv.Itoa(len(GlobalWordList.Words)))

	return nil
}

func ExecuteLongRunningTaskOnRequest(sid int) {
	//fmt.Println("ExecuteLongRunningTaskOnRequest sid = " + strconv.Itoa(sid))
    sData := GetSessionData(sid)
	
		if sData.SessionStatus.RequestExecution && !sData.SessionStatus.ExecutionStarted {
			
			// just run once
			sData.SessionStatus.RequestExecution = false
			
			d := sData.SessionStatus.PageToScan
			last1 := d[len(d)-1:]
			if (last1 == "/") {
				d = d[:len(d)-1]
				sData.SessionStatus.PageToScan = d
			}
			d = strings.Replace(d, "https://", "", 1)
			d = strings.Replace(d, "http://", "", 1)
			sData.SessionStatus.DomainsAllowed = d
			sData.SessionStatus.NumberLinksFound = 0
			sData.SessionStatus.NumberLinksVisited = 0
			sData.SessionStatus.ExecutionStarted = true
			sData.SessionStatus.ExecutionFinished = false
			sData.SessionStatus.WordsScanned = 0
			
			//fmt.Println(sData)
			
			Crawler(sid)
			
			sData.SessionWords = DeleteWordsWithOccuranceZero(sData.SessionWords)
			fmt.Println("have sData.SessionWords = " + strconv.Itoa(len(sData.SessionWords)))
			sort.Sort(SorterWordByOccurance(sData.SessionWords))
			//fmt.Println(sData.SessionWords)
			
			sData.SessionStatus.ExecutionStarted = false
			sData.SessionStatus.ExecutionFinished = true
		}
}
func Crawler(sid int) {
	//fmt.Println("Crawler sid = " + strconv.Itoa(sid))
	sData := GetSessionData(sid)
	fmt.Println("mainPage = " + sData.SessionStatus.PageToScan + ", allowedDomains = " + sData.SessionStatus.DomainsAllowed)
	
	// Instantiate default collector
	c := colly.NewCollector(
		// visit only domains
		colly.AllowedDomains(sData.SessionStatus.DomainsAllowed),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		sData.SessionStatus.NumberLinksFound++
		
		if strings.HasSuffix(link, "/") || strings.HasSuffix(link, ".html") || strings.HasSuffix(link, ".htm") {
			// Visit link found on page
			// Only those links are visited which are in AllowedDomains
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Before making a request ..
	c.OnRequest(func(r *colly.Request) {
		sData.SessionStatus.NumberLinksVisited++
		//fmt.Println("Visiting", r.URL.String())
	})
	
	// after making a request ..
	// get body from the context of the request
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(".. OnResponse")
		//fmt.Println("Content-Type=", r.Headers.Get("Content-Type"))
		
		if strings.HasPrefix(r.Headers.Get("Content-Type"), "text") {
			var t string
			
			t = string(r.Body)
			
			// from: https://stackoverflow.com/questions/44441665/how-to-extract-only-text-from-html-in-golang
			p := strings.NewReader(t)
			doc, _ := goquery.NewDocumentFromReader(p)
			doc.Find("script").Each(func(i int, el *goquery.Selection) {
				el.Remove()
			})

			// the text only from the body
			t = doc.Text()
			
			FindWordsFromText(t, sid)
		}
	})
	
	fmt.Println("start crawler")
	c.Visit(sData.SessionStatus.PageToScan)
}

func FindWordsFromText(t string, sid int) {
	//fmt.Println("FindWordsFromText sid = " + strconv.Itoa(sid))
	sData := GetSessionData(sid)
	// replace tabs
	tt := TabToSpace(t);
				
	// split the text
	tt_ := strings.Split(tt, " ")
	//fmt.Println(tt_)
	
	// do something with the result
	// here check if word is in wordlist
	//count := 0
	for _, value := range tt_ {
		ss := strings.TrimSpace(value) 
		if len(ss) > 1 { 
			sData.SessionStatus.WordsScanned++
			for i := 0; i < len(sData.SessionWords); i++ {
				if sData.SessionWords[i].Name == ss {
					sData.SessionWords[i].Occurance++
					//fmt.Println("found word = " + ss)
					//count++
					break
				}
			}
		}
	}
	
	//fmt.Println("found words = " + strconv.Itoa(count))
}
// every tab is converted into a space
// from: https://www.socketloop.com/tutorials/golang-convert-spaces-to-tabs-and-back-to-spaces-example
func TabToSpace(input string) string {
         var result []string

         for _, i := range input {
                 switch {
                 // all these considered as space, including tab \t
                 // '\t', '\n', '\v', '\f', '\r',' ', 0x85, 0xA0
                 case unicode.IsSpace(i):
                         result = append(result, " ") // replace tab with space
                 case !unicode.IsSpace(i):
                         result = append(result, string(i))
                 }
         }
         return strings.Join(result, "")
}
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