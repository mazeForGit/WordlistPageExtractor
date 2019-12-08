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
}

