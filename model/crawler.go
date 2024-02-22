package model

import (
	"fmt"
	"os"
    "net/url"
//	"net/http"
//	"encoding/json"
//	"io/ioutil"
//	"runtime"
	"strconv"
	"strings"
	"unicode"
//	"sort"
	"bytes"
	"time"
    "path/filepath"
//	"errors"
	
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	
	pdf "github.com/mazeForGit/pdf"
)

func (wl *WordList) Crawler() {
	fmt.Println(". Crawler sid = " + strconv.Itoa(wl.Session.SessionID))
	
	fmt.Println("mainPage = " + wl.Session.PageToScan + ", allowedDomains = " + wl.Session.DomainsAllowed)
	
	// Instantiate default collector
	c := colly.NewCollector(
		// user agend
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/600.7.12 (KHTML, like Gecko) Version/8.0.7 Safari/600.7.12"),
		// visit only domains
		colly.AllowedDomains(wl.Session.DomainsAllowed),
		colly.Async(true),
	)
	//c.IgnoreRobotsTxt = false
	c.Limit(&colly.LimitRule {
		DomainGlob: wl.Session.DomainsAllowed + "/*", 
		Delay: 2 * time.Second,
		Parallelism: 6,
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		wl.Session.NumberLinksFound++
		
		c.Visit(e.Request.AbsoluteURL(link))
		
	})

	// Before making a request ..
	c.OnRequest(func(r *colly.Request) {
		wl.Session.NumberLinksVisited++
//		fmt.Println(". sid = " + strconv.Itoa(wl.Session.SessionID) + " .. visiting", r.URL.String())
	})
	
	// after making a request ..
	// get body from the context of the request
	c.OnResponse(func(r *colly.Response) {
	
fmt.Println(". process: sid = " + strconv.Itoa(wl.Session.SessionID) + ", url = " + r.Request.URL.String())

		if strings.HasPrefix(r.Headers.Get("Content-Type"), "text") {
			var t string
			
			t = string(r.Body)
			//fmt.Println(t[1:30])
			
			// from: https://stackoverflow.com/questions/44441665/how-to-extract-only-text-from-html-in-golang
			p := strings.NewReader(t)
			doc, _ := goquery.NewDocumentFromReader(p)
			doc.Find("script").Each(func(i int, el *goquery.Selection) {
				el.Remove()
			})

			// the text only from the body
			t = doc.Text()
			
			wl.FindWordsFromText(t, false)
			
		} else if r.Headers.Get("Content-Type") == "application/pdf" {
	
			dirDownloadname, err := os.MkdirTemp("", "downloaded")
			//
			// save body to file
		
			filename := filepath.Join(dirDownloadname, "downloaded.pdf")
			err = os.WriteFile(filename, r.Body, 0644)
			if err != nil {
				fmt.Println(err)
			}

			//
			// open file
//fmt.Println(".. open = " + filename)	
			f, err := pdf.Open(filename)
			if err != nil {
				fmt.Println(err)
				f.Close()
			} else {

				totalPage := f.NumPage()
				var buf bytes.Buffer

				for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
				
					p := f.Page(pageIndex)
					if p.V.IsNull() {
						continue
					}
					
					texts := p.Content().Text
					var lastY = 0.0
					line := ""

					for _, text := range texts {
						if lastY != text.Y {
							if lastY > 0 {
								buf.WriteString(line + "\n")
								line = text.S
							} else {
								line += text.S
							}
						} else {
							line += text.S
						}

						lastY = text.Y
					}
					buf.WriteString(line)
				}
				
//fmt.Println(".. extracted text len = ", len(buf.String()))	
				wl.FindWordsFromText(buf.String(), true)
				
				// close the file
				err = f.Close()
				if err != nil {
					fmt.Println(err)
				}
			}
			
			// delete temp
//fmt.Println(".. delete = " + filename)	
			err = os.RemoveAll(dirDownloadname)
			if err != nil {
				fmt.Println(err)
			}
		} 
	
		// update status
		UpdateSession(wl.Session)
	
	})
	
	//fmt.Println("start crawler")
	c.Visit(wl.Session.PageToScan)
	c.Wait()
}

func (wl *WordList) FindWordsFromText(t string, isPdf bool) {

	if isPdf {
		wl.Session.PdfsScanned++
	}
	
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
		if len(ss) > 2 { 
			wl.Session.WordsScanned++
			
			for i := 0; i < len(wl.Words); i++ {
				if wl.Words[i].Name == ss {
					wl.Words[i].Occurance++
					break
				}
			}
		}
	}
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