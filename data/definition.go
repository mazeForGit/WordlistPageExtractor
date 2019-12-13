package data

import (
	//"errors"
	//"fmt"
)

type ResponseStatus struct {
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
	LastUsedSID  int	`json:"lastusedsid"`
	CountUsedSID  int	`json:"countusedsid"`
    WordListUrl  string	`json:"wordlisturl"`
}
type Status struct {
    Count int `json:"sessioncount"`
    RequestExecution bool `json:"requestexecution"`
	PageToScan string `json:"pagetoscan"`
	DomainsAllowed string `json:"domainsallowed"`
	NumberLinksFound int `json:"numberlinksfound"`
	NumberLinksVisited int `json:"numberlinksvisited"`
	ExecutionStarted bool `json:"executionstarted"`
	ExecutionFinished bool `json:"executionfinished"`
	WordsScanned int `json:"wordsscanned"`
}
type SessionData struct {
    SessionID int `json:"sid"`
    SessionStatus Status `json:"sessionstatus"`
	SessionWords []Word `json:"words"`
}
func DeleteWordsWithOccuranceZero(wl []Word) ([]Word) {
    var wlNew []Word
	for i := 0; i < len(wl); i++ {
		if wl[i].Occurance > 0 {
			wl[i].Tests = nil
			wlNew = append(wlNew, wl[i])
		}
	}
    return wlNew
}
func CopyWords(wl []Word) ([]Word) {
    var wlNew []Word
	for i := 0; i < len(wl); i++ {
		var wNew Word
		wNew.Tests = wl[i].Tests
		wNew.Name = wl[i].Name
		wNew.Occurance = wl[i].Occurance
		wlNew = append(wlNew, wNew)
	}
    return wlNew
}