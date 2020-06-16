package main

// tea describes basic details of what makes up a tea
type Tea struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Id      string    `json:"Id"`
	Maker    string `json:"Maker"`
	Owner   string `json:"Owner"`
	Weight  string `json:"Weight"`
	//Histories []HistoryItem `json:"histories"`
}

//type HistoryItem struct {
//	TxId string
//	tea Tea
//}