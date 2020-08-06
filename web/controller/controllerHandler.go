package controller

import (
	"blc-iot-demo/web/service"
	"blc-iot-demo/web/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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
	app.Setup.ModifyQueryCount(tea.Id)
	data.Tea = tea
	fmt.Printf("%+v",data.Tea)
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
	if height < curHeight-1 {
		block, _ := app.Setup.LedgerClient.QueryBlock(height + 1)
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
		} else {
			data.Msg = "无权访问"
			ShowView(w, r, "index.html", data)
			return
		}
	} else if !data.IsLogin {
		ShowView(w, r, "AccountRelated/login.html", data)
		return
	}
}

// 添加信息
func (app *Application) AddTea(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsStaff {
		uuid := utils.CreateUUID()
		productionDate := utils.SwitchTimeStampToData(time.Now().Unix())
		tea := service.Tea{
			Id:              uuid,
			Name:            r.FormValue("teaName"),
			Maker:           r.FormValue("teaMaker"),
			Owner:           r.FormValue("teaOwner"),
			Weight:          r.FormValue("teaWeight"),
			Origin:          r.FormValue("teaOrigin"),
			Production_Date: productionDate,
			Shelf_life:      "18个月",
			TxID:            "",
		}
		fmt.Println("---------------------------------------------")
		fmt.Println("写入茶叶数据")
		data.Tea = tea

		txId, _ := app.Setup.SaveTea(tea)
		tea.TxID = txId

		ShowView(w, r, "StaffOption/addSuccess.html", data)
		return

	} else if !data.IsStaff {
		data.Msg = "无权访问"
		ShowView(w, r, "index.html", data)
		return
	}
}

// 根据teaID查询信息
func (app *Application) ModifyQuery(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsStaff || data.IsSuperAdmin || data.IsAdmin {
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
	if data.IsStaff || data.IsSuperAdmin || data.IsAdmin {

		teaId := r.FormValue("teaId")
		teaName := r.FormValue("teaName")
		teaMaker := r.FormValue("teaMaker")
		teaOwner := r.FormValue("teaOwner")
		teaWeight := r.FormValue("teaWeight")
		teaOrigin := r.FormValue("teaOrigin")

		tea := service.Tea{
			Id:     teaId,
			Name:   teaName,
			Maker:  teaMaker,
			Owner:  teaOwner,
			Weight: teaWeight,
			Origin: teaOrigin,
		}

		data.Tea = tea

		_, err := app.Setup.ModifyTea(teaId, teaOwner)
		if err != nil {
			log.Println("modufied teas failed, err:", err)
			return
		}
		ShowView(w, r, "StaffOption/modifySuccess.html", data)
		return
	} else  {
		data.Msg = "无权访问"
		ShowView(w, r, "index.html", data)
		return
	}
}
