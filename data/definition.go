package data

import (
	//"errors"
	//"fmt"
)

type Status struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
type Test struct {
    Name  string	`json:"name"`
    Category  string	`json:"category"`
}

type Word struct {
	Id int	`json:"id"`
    Name  string	`json:"name"`
	Occurance int	`json:"occurance"`
	New bool	`json:"new"`
	Tests []Test	`json:"tests"`
}

type WordList struct {
    Type  string	`json:"type"`
    LastUsedId  int	`json:"lastusedid"`
    Count  int	`json:"count"`
	Tests []Test	`json:"tests"`
	Words []Word	`json:"words"`
}

type Config struct {
    RequestExecution  bool	`json:"requestexecution"`
    WordListUrl  string	`json:"wordlisturl"`
	PageToScan string `json:"pagetoscan"`
	DomainsAllowed string `json:"domainsallowed"`
	NumberLinksFound int `json:"numberlinksfound"`
	NumberLinksVisited int `json:"numberlinksvisited"`
	ExecutionStarted bool `json:"executionstarted"`
	ExecutionFinished bool `json:"executionfinished"`
	WordsScanned int `json:"wordsscanned"`
}

func DeleteWordsWithOccuranceZero(wl WordList) (WordList) {
    var newwl WordList
	for i := 0; i < len(wl.Words); i++ {
		if wl.Words[i].Occurance > 0 {
			wl.Words[i].Tests = nil
			newwl.Words = append(newwl.Words, wl.Words[i])
			newwl.Count++
		}
	}
    return newwl
}