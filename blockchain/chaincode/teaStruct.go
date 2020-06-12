package main

// tea describes basic details of what makes up a tea
type Tea struct {
	ObjectType string `json:"docType"`
	Id      string    `json:"Id"`
	Maker    string `json:"maker"`
	Owner   string `json:"Owner"`
	Weight  string `json:"Weight"`
	Histories []HistoryItem `json:"Histories"`
}

type HistoryItem struct {
	TxId string
	tea Tea
}