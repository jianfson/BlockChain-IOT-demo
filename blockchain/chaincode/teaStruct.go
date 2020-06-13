package main

// tea describes basic details of what makes up a tea
type Tea struct {
	ObjectType string `json:"docType"`
	Id      string    `json:"id"`
	Maker    string `json:"maker"`
	Owner   string `json:"owner"`
	Weight  string `json:"weight"`
	Histories []HistoryItem `json:"histories"`
}

type HistoryItem struct {
	TxId string
	tea Tea
}