package data

import (
	"fmt"
	//"net"
    "net/url"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"sort"
	"bytes"
	"time"
	
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var GlobalConfig Config
var GlobalWordList WordList
var GlobalWordListStorage []WordList

type SorterWordByOccurance []Word

func (a SorterWordByOccurance) Len() int           { return len(a) }
func (a SorterWordByOccurance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SorterWordByOccurance) Less(i, j int) bool { return a[i].Occurance > a[j].Occurance }

func GetNewSessionID() int {
	GlobalConfig.LastUsedSID++
	return GlobalConfig.LastUsedSID
}
func GetWordListForSession(sid int) *WordList {
	//fmt.Println("GetSessionData sid=" + strconv.Itoa(sid))
	
	for i := 0; i < len(GlobalWordListStorage); i++ {
		if GlobalWordListStorage[i].Session.SessionID == sid {
			return &GlobalWordListStorage[i]
		}
	}
	
	var sData = SessionStatus {
		SessionID: sid,
		Count: 0,
		RequestExecution: false,
		ExecutionStarted: false,
		ExecutionFinished: false,
		PageToScan: "",
		DomainsAllowed: "",
		NumberLinksFound: 0,
		NumberLinksVisited: 0,
		WordsScanned: 0,
	}
	var wl = WordList {
		Session: sData,
		Words: nil,
		Tests: nil,
	}
	GlobalWordListStorage = append(GlobalWordListStorage, wl)
    return GetWordListForSession(sid)
}
func SetSessionData(sData SessionStatus) {
	//fmt.Println("SetSessionData sid=" + strconv.Itoa(sData.SessionID))
	
	for i := 0; i < len(GlobalWordListStorage); i++ {
		if GlobalWordListStorage[i].Session.SessionID == sData.SessionID {
			GlobalWordListStorage[i].Session = sData
		}
	}
}

func ReadGlobalWordlistFromRemote() error {
	fmt.Println("ReadGlobalWordlist")
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
func StoreWordlistAtRemote(wl WordList) error {
	fmt.Println("StoreWordlist")
	
    var err error
	var resp *http.Response
	var body []byte
	var requestUrl string = ""
	//var client http.Client
	
    requestUrl = GlobalConfig.WordListStorageUrl + "/wordlist"
    fmt.Println("connect to wordliststorage = " + requestUrl)
    //fmt.Println(wl)
	
    payload, err := json.Marshal(wl)
	if err != nil {
        return err
    }
	
    req, err := http.NewRequest("PUT", requestUrl, bytes.NewBuffer(payload))
	if err != nil {
        return err
    }
	
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err = client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
	var s ResponseStatus
    json.Unmarshal(body, &s)
    fmt.Println("got s.Code = " + strconv.Itoa(s.Code))
    fmt.Println("got s.Text = " + s.Text)
	
	return err
}
// from: https://gobyexample.com/url-parsing
func GetDomainAndMainPageByUrl(s string) (string, string) {
	//fmt.Println("GetDomainAndMainPageByUrl .. url = " + s)
	
	var domain string = ""
	var mainpage string = ""
	
	u, err := url.Parse(s)
    if err != nil {
        fmt.Println(err)
    }
	
	domain = u.Host
	
	if (u.Scheme == "") {
		mainpage = "https://" + u.Host
	} else {
		mainpage = u.Scheme + "://" + u.Host
	}
	
	return mainpage, domain
}
func ExecuteLongRunningTaskOnRequest(sid int) {
	fmt.Println("ExecuteLongRunningTaskOnRequest sid = " + strconv.Itoa(sid))
    sData := GetWordListForSession(sid)
	//fmt.Println(sData.Session)
	
		if sData.Session.RequestExecution && !sData.Session.ExecutionStarted {
			
			// just run once
			sData.Session.RequestExecution = false
			
			sData.Session.PageToScan = strings.TrimSpace(sData.Session.PageToScan)
			m, d := GetDomainAndMainPageByUrl(sData.Session.PageToScan)
			sData.Session.PageToScan = m
			sData.Session.DomainsAllowed = d
			sData.Session.NumberLinksFound = 0
			sData.Session.NumberLinksVisited = 0
			sData.Session.ExecutionStarted = true
			sData.Session.ExecutionFinished = false
			sData.Session.WordsScanned = 0
			
			//fmt.Println(sData)
			
			fmt.Println("before: have sData.Words = " + strconv.Itoa(len(sData.Words)))
			
			Crawler(sid)
			
			fmt.Println("after: have sData.SessionWords = " + strconv.Itoa(len(sData.Words)))
			sData.Words = DeleteWordsWithOccuranceZero(sData.Words)
			fmt.Println("after delete: have sData.Words = " + strconv.Itoa(len(sData.Words)))
			sort.Sort(SorterWordByOccurance(sData.Words))
			//fmt.Println(sData.Words)
			
			sData.Session.ExecutionStarted = false
			sData.Session.ExecutionFinished = true
			
			if len(sData.Words) > 0 && GlobalConfig.WordListStorageUrl != "" {
				StoreWordlistAtRemote(*sData)
			}
		}
}
func Crawler(sid int) {
	fmt.Println("Crawler sid = " + strconv.Itoa(sid))
	sData := GetWordListForSession(sid)
	fmt.Println("mainPage = " + sData.Session.PageToScan + ", allowedDomains = " + sData.Session.DomainsAllowed)
	
	// Instantiate default collector
	c := colly.NewCollector(
		// user agend
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/600.7.12 (KHTML, like Gecko) Version/8.0.7 Safari/600.7.12"),
		// visit only domains
		colly.AllowedDomains(sData.Session.DomainsAllowed),
		colly.Async(true),
	)
	//c.IgnoreRobotsTxt = false
	c.Limit(&colly.LimitRule {
		DomainGlob: sData.Session.DomainsAllowed + "/*", 
		Delay: 2 * time.Second,
		Parallelism: 6,
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		sData.Session.NumberLinksFound++
		
		// remark:
		// this want work for query or parameters at the end, see e.g. bdbos
		//if strings.HasSuffix(link, "/") || strings.HasSuffix(link, ".html") || strings.HasSuffix(link, ".htm") {
			// Visit link found on page
			// Only those links are visited which are in AllowedDomains
			c.Visit(e.Request.AbsoluteURL(link))
		//}
	})

	// Before making a request ..
	c.OnRequest(func(r *colly.Request) {
		sData.Session.NumberLinksVisited++
		fmt.Println("sid = " + strconv.Itoa(sid) + " .. visiting", r.URL.String())
	})
	
	// after making a request ..
	// get body from the context of the request
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(".. OnResponse")
		//fmt.Println("Content-Type=", r.Headers.Get("Content-Type"))
		
		if strings.HasPrefix(r.Headers.Get("Content-Type"), "text") {
			var t string
			
			t = string(r.Body)
			//fmt.Println(t)
			
			// from: https://stackoverflow.com/questions/44441665/how-to-extract-only-text-from-html-in-golang
			p := strings.NewReader(t)
			doc, _ := goquery.NewDocumentFromReader(p)
			doc.Find("script").Each(func(i int, el *goquery.Selection) {
				el.Remove()
			})

			// the text only from the body
			t = doc.Text()
			//fmt.Println(t)
			FindWordsFromText(t, sid)
		}
	})
	
	//fmt.Println("start crawler")
	c.Visit(sData.Session.PageToScan)
	c.Wait()
}

func FindWordsFromText(t string, sid int) {
	//fmt.Println("FindWordsFromText sid = " + strconv.Itoa(sid))
	sData := GetWordListForSession(sid)
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
			sData.Session.WordsScanned++
			for i := 0; i < len(sData.Words); i++ {
				if sData.Words[i].Name == ss {
					sData.Words[i].Occurance++
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