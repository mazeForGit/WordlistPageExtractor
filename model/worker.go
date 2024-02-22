package model

import (
	"fmt"
	"strconv"
	//"strings"
	"sort"
)

// Job represents the job to be run
type Job struct {
	Wordlist WordList
}
var JobChannel chan Job

func Worker(i int, jobChannel <-chan Job) {
	fmt.Println(". start worker .. i = " + strconv.Itoa(i))

	for {
		
		select {
		case job := <-jobChannel:
			
			// we have received a work request.
			fmt.Println(". Worker .. session: ", job.Wordlist.Session)
		
			if job.Wordlist.Session.RequestExecution && !job.Wordlist.Session.ExecutionStarted {
				
				job.Wordlist.Words = GlobalWordList.Words
				
				// just run once
				job.Wordlist.Session.RequestExecution = false		
				job.Wordlist.Session.ExecutionStarted = true
				job.Wordlist.Session.ExecutionFinished = false

				m, d := GetDomainAndMainPageByUrl(job.Wordlist.Session.PageToScan)
				job.Wordlist.Session.PageToScan = m
				job.Wordlist.Session.DomainsAllowed = d
				job.Wordlist.Session.NumberLinksFound = 0
				job.Wordlist.Session.NumberLinksVisited = 0
				job.Wordlist.Session.WordsScanned = 0
				job.Wordlist.Session.PdfsScanned = 0

				UpdateSession(job.Wordlist.Session)
				
				fmt.Println("before: have #Words = " + strconv.Itoa(len(job.Wordlist.Words)))
				
				job.Wordlist.Crawler()
				
				fmt.Println("after: have #SessionWords = " + strconv.Itoa(len(job.Wordlist.Words)))
				job.Wordlist.Words = DeleteWordsWithOccuranceZero(job.Wordlist.Words)
				fmt.Println("after delete: have #Words = " + strconv.Itoa(len(job.Wordlist.Words)))
				sort.Sort(SorterWordByOccurance(job.Wordlist.Words))
				
				job.Wordlist.Session.ExecutionStarted = false
				job.Wordlist.Session.ExecutionFinished = true	
				
				// store local
				UpdateWordlist(job.Wordlist)
				
				// store remote
				if len(job.Wordlist.Words) > 0 && GlobalConfig.WordListStorageUrl != "" {
					StoreWordlistAtRemote(job.Wordlist)
				}
			}
		}
	}
}