package controller

import (
	"blc-iot-demo/blockchain/service"
	"blc-iot-demo/web/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Application struct {
	Setup *service.ServiceSetup
}

// 进入查询页面
func (app *Application) QueryPage(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	ShowView(w, r, "PublicOption/queryPage.html", data)
}

// 根据teaID查询信息
func (app *Application) FindTeaByID(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	teaID := r.FormValue("id")
	result, err := app.Setup.FindTeaInfoByID(teaID)
	if err != nil {
		log.Println(err)
	}
	var tea = service.Tea{}

	json.Unmarshal(result, &tea)
	data.Tea = tea
	block, err := app.Setup.QueryBlockByTxID(tea.TxID)
	if err != nil {
		log.Println("query block failed, err:", err)
		return
	}

	// 区块头
	BlockHeader := block.Header

	// 区块高度
	height := BlockHeader.Number
	data.Block.Height = height
	fmt.Println("tx 所在区块：", height)

	// 区块 data hash
	hash := BlockHeader.DataHash
	dataHash := hex.EncodeToString(hash)
	data.Block.DataHash = dataHash

	//当前区块的 block hash, ledgerCli 或下一区块获取
	blcInfo, err := app.Setup.LedgerClient.QueryInfo()

	// 从零开始， 类比切片
	curHeight := blcInfo.BCI.Height
	fmt.Println("curHeight:", curHeight)
	if height < curHeight - 1 {
		block, _ := app.Setup.LedgerClient.QueryBlock(height+1)
		bh := hex.EncodeToString(block.GetHeader().PreviousHash)
		data.Block.BlcHash = bh
		fmt.Println("bh=", bh)
	} else {
		bh := hex.EncodeToString(blcInfo.BCI.CurrentBlockHash)
		data.Block.BlcHash = bh
		fmt.Println("bh=", bh)
	}

	ShowView(w, r, "PublicOption/queryResult.html", data)
}

// 显示添加信息页面
func (app *Application) AddTeaPage(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsLogin {
		if data.IsStaff {
			ShowView(w, r, "StaffOption/addTeaPage.html", data)
			return
		}else {
			data.Msg = "无权访问"
			ShowView(w, r, "index.html", data)
			return
		}
	} else if !data.IsLogin{
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

// 添加信息
func (app *Application) AddTea(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsStaff {
		tea := service.Tea{
			Id:     r.FormValue("new_id"),
			Maker:  r.FormValue("new_maker"),
			Owner:  r.FormValue("new_owner"),
			Weight: r.FormValue("new_weight"),
		}
		data.Tea = tea

		app.Setup.SaveTea(tea)

		ShowView(w, r, "StaffOption/addSuccess.html", data)
		return
	} else if !data.IsStaff {
		data.Msg = "无权访问"
		ShowView(w, r, "index.html", data)
		return
	}
}

// 进入修改查询页面
func (app *Application) ModifyQueryPage(w http.ResponseWriter, r *http.Request) {

	data := utils.CheckLogin(r)

	if data.IsStaff {
		ShowView(w, r, "StaffOption/modifyQueryPage.html", data)
		return
	} else {
		data.Msg = "无权访问"
		ShowView(w, r, "index.html", data)
		return
	}
}

// 根据teaID查询信息
func (app *Application) ModifyQuery(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsStaff||data.IsSuperAdmin||data.IsAdmin {

		teaID := r.FormValue("id")
		result, _ := app.Setup.FindTeaInfoByID(teaID)
		var tea = service.Tea{}
		json.Unmarshal(result, &tea)

		data.Tea = tea

		ShowView(w, r, "StaffOption/modifyPage.html", data)
		return
	} else {
		data.Msg = "无权访问"
		ShowView(w, r, "index.html", data)
		return
	}
}

// 修改信息
func (app *Application) ModifyResult(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	teaId := r.FormValue("new_id")
	nextOwner := r.FormValue("new_owner")

	tea := service.Tea{
		Id:     r.FormValue("new_id"),
		Maker:  r.FormValue("new_maker"),
		Owner:  r.FormValue("new_owner"),
		Weight: r.FormValue("new_weight"),
	}

	data.Tea = tea

	fmt.Println(tea.Id)

	_, err := app.Setup.ModifyTea(teaId, nextOwner)
	if err != nil {
		log.Println("modufied teas failed, err:", err)
		return
	}

	if data.IsStaff {
		ShowView(w, r, "StaffOption/modifySuccess.html", data)
		return
	} else if !data.IsStaff {
		data.Msg = "无权访问"
		ShowView(w, r, "index.html", data)
		return
	}
}
func (app *Application) GetHistoryByIDPage(w http.ResponseWriter, r *http.Request) {
	ShowView(w, r, "PublicOption/queryHistoryPage.html", nil)
}
func (app *Application) GetHistoryById(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	id := r.FormValue("id")
	history, err := app.Setup.GetHistoryForTea(id)
	if err != nil{
		log.Println(history)
		log.Println(string(history))
	}
	type MyJsonName struct {
		IsDelete  string `json:"IsDelete"`
		Timestamp string `json:"Timestamp"`
		TxID      string `json:"TxId"`
		Value     struct {
			DocType        string `json:"docType"`
			ID             string `json:"id"`
			Maker          string `json:"maker"`
			Name           string `json:"name"`
			Origin         string `json:"origin"`
			Owner          string `json:"owner"`
			ProductionDate string `json:"production_date"`
			ShelfLife      string `json:"shelf_life"`
			TxID           string `json:"txID"`
			Weight         string `json:"weight"`
		} `json:"Value"`
	}
	teaHistory := []MyJsonName{}
	err = json.Unmarshal(history, &teaHistory)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}
	log.Printf("%+v", teaHistory)
	s := string(history)
	data.History = s
	ShowView(w, r, "PublicOption/queryHistoryResult.html", data)
}
