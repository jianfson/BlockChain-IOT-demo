package main

// tea describes basic details of what makes up a tea
type Tea struct {
	ObjectType string `json:"docType"`
	Id      string    `json:"id"`
	Maker    string `json:"make"`
	Owner   string `json:"owner"`
	Weight  string `json:"weight"`
	Histories []HistoryItem `json:"history"`
}

type HistoryItem struct {
	TxId string
	tea Tea
}