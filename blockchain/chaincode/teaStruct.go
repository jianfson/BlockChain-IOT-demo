package main

// tea describes basic details of what makes up a tea
type Tea struct {
	ObjectType      string `json:"docType"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	Maker           string `json:"maker"`
	Owner           string `json:"owner"`
	Weight          string `json:"weight"`
	Origin          string `json:"origin"`
	Production_Date string `json:"production_date"`
	Shelf_life      string `json:"shelf_life"`
	TxID            string `json:"txID"`
}
