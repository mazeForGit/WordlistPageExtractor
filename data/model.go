package data

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"runtime"
	"time"
	"strconv"
	"strings"
	"unicode"
	"sort"
	
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var GlobalConfig Config
var GlobalWordList WordList

type SorterWordByOccurance []Word

func (a SorterWordByOccurance) Len() int           { return len(a) }
func (a SorterWordByOccurance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SorterWordByOccurance) Less(i, j int) bool { return a[i].Occurance > a[j].Occurance }

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
		//PrintMemUsage()
		time.Sleep(2 * time.Second)
		if GlobalConfig.RequestExecution && !GlobalConfig.ExecutionStarted {
			
			// just run once
			GlobalConfig.RequestExecution = false
			
			d := GlobalConfig.PageToScan
			last1 := d[len(d)-1:]
			if (last1 == "/") {
				d = d[:len(d)-1]
				GlobalConfig.PageToScan = d
			}
			d = strings.Replace(d, "https://", "", 1)
			d = strings.Replace(d, "http://", "", 1)
			GlobalConfig.DomainsAllowed = d
			GlobalConfig.NumberLinksFound = 0
			GlobalConfig.NumberLinksVisited = 0
			GlobalConfig.ExecutionStarted = true
			GlobalConfig.ExecutionFinished = false
			GlobalConfig.WordsScanned = 0
			
			//fmt.Println(GlobalConfig)
			
			readSite()
			
			GlobalWordList = DeleteWordsWithOccuranceZero(GlobalWordList)
			fmt.Println("have GlobalWordList.Words = " + strconv.Itoa(len(GlobalWordList.Words)))
			sort.Sort(SorterWordByOccurance(GlobalWordList.Words))
			//fmt.Println(GlobalWordList)
			
			GlobalConfig.ExecutionStarted = false
			GlobalConfig.ExecutionFinished = true
		}
	}
}

func readSite() {
	fmt.Println("readSite")
	
    var err error
	var resp *http.Response
	var body []byte
	var requestUrl string = ""
	
    requestUrl = GlobalConfig.WordListUrl + "/words?testOnly=true&format=json"
    fmt.Println("connect to wordlist and get words with tests = " + requestUrl)
    resp, err = http.Get(requestUrl)
    if err != nil {
		fmt.Println("error ..")
        fmt.Println(err.Error())
        return
    }

    defer resp.Body.Close()
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
		fmt.Println("error ..")
        fmt.Println(err.Error())
        return
    }
    
    json.Unmarshal(body, &GlobalWordList.Words)
    fmt.Println("got GlobalWordList.Words = " + strconv.Itoa(len(GlobalWordList.Words)))
        
    // do something with the result
	Crawler(GlobalConfig.PageToScan, GlobalConfig.DomainsAllowed)
	
	return
}

func Crawler(mainPage string, allowedDomains string) {
	fmt.Println("Crawler .. mainPage = " + mainPage + ", allowedDomains = " + allowedDomains)
	
	// Instantiate default collector
	c := colly.NewCollector(
		// visit only domains
		colly.AllowedDomains(allowedDomains),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		GlobalConfig.NumberLinksFound++
		
		if strings.HasSuffix(link, "/") || strings.HasSuffix(link, ".html") || strings.HasSuffix(link, ".htm") {
			// Visit link found on page
			// Only those links are visited which are in AllowedDomains
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Before making a request ..
	c.OnRequest(func(r *colly.Request) {
		GlobalConfig.NumberLinksVisited++
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
			
			FindWordsFromText(t)
		}
	})
	
	fmt.Println("start crawler")
	c.Visit(mainPage)
}

func FindWordsFromText(t string) {
	// replace tabs
	tt := TabToSpace(t);
				
	// split the text
	tt_ := strings.Split(tt, " ")
	//fmt.Println(tt_)
	
	// do something with the result
	// here check if word is in wordlist
	count := 0
	for _, value := range tt_ {
		ss := strings.TrimSpace(value) 
		if len(ss) > 1 { 
			GlobalConfig.WordsScanned++
			for i := 0; i < len(GlobalWordList.Words); i++ {
				if GlobalWordList.Words[i].Name == ss {
					GlobalWordList.Words[i].Occurance++
					//fmt.Println("found word = " + ss)
					count++
					break
				}
			}
		}
	}
	
	//fmt.Println("found words = " + strconv.Itoa(count))
}
// every tab is converted into a space
//
// from: https://www.socketloop.com/tutorials/golang-convert-spaces-to-tabs-and-back-to-spaces-example
//
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