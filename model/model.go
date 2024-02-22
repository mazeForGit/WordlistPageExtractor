package model

import (
	"fmt"
	"os"
//    "net/url"
	"net/http"
	"encoding/json"
	"io/ioutil"
//	"runtime"
	"strconv"
//	"strings"
//	"unicode"
//	"sort"
	"bytes"
//	"time"
 //   "path/filepath"
	"errors"
)

var GlobalConfig Config
var GlobalWordList WordList
var GlobalWordListStorage []WordList

type SorterWordByOccurance []Word

func (a SorterWordByOccurance) Len() int           { return len(a) }
func (a SorterWordByOccurance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SorterWordByOccurance) Less(i, j int) bool { return a[i].Occurance > a[j].Occurance }

//
// new session returns existing or new session
//
func NewSession(session Session) (WordList, error) {
	fmt.Println(". NewSession .. session: ", session)
		
	_, d := GetDomainAndMainPageByUrl(session.PageToScan)
	for i := 0; i < len(GlobalWordListStorage); i++ {
		if GlobalWordListStorage[i].Session.DomainsAllowed == d {
			// already existing
			return GlobalWordListStorage[i], errors.New("A session for this domain already exists.")
		}
	}
	
	// new
	session.SessionID = len(GlobalWordListStorage) + 1
	session.DomainsAllowed = d
	session.RequestExecution = true
	var wl = WordList {
		Session: session,
		Words: nil,
		Tests: nil,
	}
	GlobalWordListStorage = append(GlobalWordListStorage, wl)
	
	return GetWordListForSession(session.SessionID)
}
func GetWordListForSession(sid int) (WordList, error) {
	for i := 0; i < len(GlobalWordListStorage); i++ {
		if GlobalWordListStorage[i].Session.SessionID == sid {
			return GlobalWordListStorage[i], nil
		}
	}
	var wl WordList
    return wl, errors.New("no results for sid=" + strconv.Itoa(sid))
}
//
// store session at global
//	this is subject to concurency
//
func UpdateSession(s Session) {

	for i := 0; i < len(GlobalWordListStorage); i++ {
		if GlobalWordListStorage[i].Session.SessionID == s.SessionID {
			GlobalWordListStorage[i].Session = s
		
			break
		}
	}
}
//
// store wordlist at global
//	this is subject to concurency
//
func UpdateWordlist(wl WordList) {

	for i := 0; i < len(GlobalWordListStorage); i++ {
		if GlobalWordListStorage[i].Session.SessionID == wl.Session.SessionID {
			GlobalWordListStorage[i].Session = wl.Session
			GlobalWordListStorage[i].Words = wl.Words
			break
		}
	}
}

//
// read configuration from file
//
func ReadConfigurationFromFile(fileName string) (err error) {
	fmt.Println(". ReadConfigurationFromFile .. fileName = " + fileName)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
    
	defer file.Close() // defer the closing for later parsing
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	
	err = json.Unmarshal([]byte(byteValue), &GlobalConfig)
	return err
}
//
// read wordListStorage from file
//
func ReadWordListStorageFromFile(fileName string) (err error) {
	fmt.Println(". ReadWordListStorageFromFile .. fileName = " + fileName)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	
	defer file.Close() // defer the closing for later parsing
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	
	err = json.Unmarshal([]byte(byteValue), &GlobalWordListStorage)
    return err
}
//
// read wordList from file
//
func ReadWordListFromFile(fileName string) (err error) {
	fmt.Println(". ReadWordListFromFile .. fileName = " + fileName)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	
	defer file.Close() // defer the closing for later parsing
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	
	err = json.Unmarshal([]byte(byteValue), &GlobalWordList)
	fmt.Println("have GlobalWordlist.Words = " + strconv.Itoa(len(GlobalWordList.Words)))
	
    return err
}
func ReadGlobalWordlistFromRemote() error {
	fmt.Println("ReadGlobalWordlist")
	fmt.Println("have GlobalWordlist.Words = " + strconv.Itoa(len(GlobalWordList.Words)))
	
    var err error
	var resp *http.Response
	var body []byte
	var requestUrl string = ""
	
	// would reduce the number of words
    //requestUrl = GlobalConfig.WordListUrl + "/wordlist?testOnly=true"
	// but here we need all to find out which words without test are used as well
	requestUrl = GlobalConfig.WordListUrl + "/words?format=json"
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
    
    //json.Unmarshal(body, &GlobalWordList)
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
    fmt.Println("got status.Code = " + strconv.Itoa(s.Code))
    fmt.Println("got status.Text = " + s.Text)
	
	return err
}