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
	Origin_IP       IP    `json:"origin_ip"`
	Production_Date string `json:"production_date"`
	Shelf_life      string `json:"shelf_life"`
	TxID            string `json:"txID"`
	Size 			string `json:"size"`
	QueryCounter    int    `json:"queryCounter"`
	Boxed           Box    `json:"box"`
}

//成都市 经度:104.07 纬度:30.67
type IP struct {
	Longitude string
	Latitude  string
}

type Box struct {
	Boxed bool `json:"boxed"`
	Num   int  `json:"num"`
}
