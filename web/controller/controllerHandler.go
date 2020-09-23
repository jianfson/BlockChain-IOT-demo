package controller

import (
	"blc-iot-demo/web/dao"
	"blc-iot-demo/web/model"
	"blc-iot-demo/web/service"
	"blc-iot-demo/web/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/thinkeridea/go-extend/exnet"
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
	err = json.Unmarshal(result, &tea)

	if err != nil {
		log.Println("unmarshal failed, err:", err)
	}

	data.Tea = tea
fmt.Println("----->",data.Tea)
	userIP := exnet.ClientIP(r)
	fmt.Println(userIP)

	_, _ = app.Setup.ModifyQueryCount(tea.Id)

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

	if data.IsLogin {
		actionTime := utils.SwitchTimeStampToData(time.Now().Unix())
		action := model.Action{
			UserID:           data.Sess.UserID,
			ActionType:       "查询",
			ActionTargetName: data.Tea.Name,
			ActionTargetID:   data.Tea.Id,
			ActionTime:       actionTime,
		}
		_, _ = dao.InsertAction(action)
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

/*
func (app *Application) AddTea(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)
	if data.IsStaff {

		//获取表单输入
		uuid := utils.CreateUUID()
		teaName := r.FormValue("teaNameBulk")
		teaMaker := r.FormValue("teaMakerBulk")
		teaOwner := r.FormValue("teaOwnerBulk")
		teaWeight := r.FormValue("teaWeightBulk")
		teaOrigin := r.FormValue("teaOriginBulk")
		productionDate := utils.SwitchTimeStampToData(time.Now().Unix())	//该批次包装时间
		shelfLife := "18个月"

		tea := service.Tea{
			Id:             uuid,
			Name:           teaName,
			Maker:          teaMaker,
			Owner:          teaOwner,
			Weight:         teaWeight,
			Origin:         teaOrigin,
			Production_Date: productionDate,
			Shelf_life:      shelfLife,
			TxID:           "",
		}
		fmt.Println("------>save tea,", tea)
		txId, err := app.Setup.SaveTea(tea)
		if err !=nil {
			fmt.Println("errr 2:", err)
		}
		data.Tea = tea
		tea.TxID = txId

		actionTime := utils.SwitchTimeStampToData(time.Now().Unix())
		action := model.Action{
			UserID:           data.Sess.UserID,
			ActionType:       "单件上链",
			ActionTargetName: data.Tea.Name,
			ActionTargetID:   data.Tea.Id,
			ActionTime:       actionTime,
		}
		_, _ = dao.InsertAction(action)

		ShowView(w, r, "StaffOption/addSuccess.html", data)
		return
	} else if !data.IsStaff {
		ShowView(w, r, "index.html", data)
		return
	}
}
*/
// 添加信息
//8.7 新代码
func (app *Application) BulkAddTea(w http.ResponseWriter, r *http.Request) {
	data := utils.CheckLogin(r)

	if data.IsStaff {
		data.IsBulkAdd = true

		//获取表单输入
		batchSizeInCHN := r.FormValue("batchSizeInCHN")		//批次大小
		teaName := r.FormValue("teaNameBulk")
		teaMaker := r.FormValue("teaMakerBulk")
		teaOwner := r.FormValue("teaOwnerBulk")
		teaWeight := r.FormValue("teaWeightBulk")
		teaOrigin := r.FormValue("teaOriginBulk")
		lng := "104"
		lat := "23"
		ip := service.IP{
			Longitude: lng,
			Latitude: lat,
		}

		//lnglat := "104,52"

		shelfLife := "18个月"
		codePath := "./QRcode/"

		tea := service.Tea{
			ObjectType: "teaObj",
			Id:              "",
			Name:            teaName,
			Maker:           teaMaker,
			Owner:           teaOwner,
			Weight:          teaWeight,
			Origin:          teaOrigin,
			Production_Date: "",
			Shelf_life:      shelfLife,
			TxID:            "",
			Origin_IP:ip,
			Size: "small",
			QueryCounter: 0,
			Boxed: service.Box{

			},
		}


		if batchSizeInCHN == "十条" {
			batchSize := 10
			for i := 0; i < batchSize; i++{
				productionDate := utils.SwitchTimeStampToData(time.Now().Unix())	//该批次包装时间
				tea.Production_Date = productionDate
				uuid := utils.CreateUUID()			// 生成1000条uuid
				tea.Id = uuid
				fmt.Println("----->", tea,"JINGWEIDU",tea.Origin_IP)

				txID, err := app.Setup.SaveTea(tea)
				tea.TxID = txID
				if err !=nil {
					fmt.Println("save tea err:", err)
				}

				wholePath := codePath + uuid +".png"

				//生成1000张二维码并保存到文件夹，等待上传
				err = qrcode.WriteFile("http://47.108.134.136:9000/findTeaByID?id="+uuid, qrcode.Medium, 256, wholePath)
				fmt.Println(err)

			}
		} else if batchSizeInCHN == "一万条" {
			batchSize := 10000
			for i := 0; i < batchSize; i++{
				productionDate := utils.SwitchTimeStampToData(time.Now().Unix())	//该批次包装时间
				tea.Production_Date = productionDate
				uuid := utils.CreateUUID()			// 生成1000条uuid
				tea.Id = uuid
				_, _ = app.Setup.SaveTea(tea)

				wholePath := codePath + uuid +".png"

				err := qrcode.WriteFile("http://47.108.134.136:9000/findTeaByID?id="+uuid, qrcode.Medium, 256, wholePath)
				fmt.Println(err)

			}
		}else if batchSizeInCHN == "十万条"{
			batchSize := 100000
			for i := 0; i < batchSize; i++{
				productionDate := utils.SwitchTimeStampToData(time.Now().Unix())	//该批次包装时间
				tea.Production_Date = productionDate
				uuid := utils.CreateUUID()			// 生成1000条uuid
				tea.Id = uuid
				_, _ = app.Setup.SaveTea(tea)

				wholePath := codePath + uuid +".png"

				//生成1000张二维码并保存到文件夹，等待上传
				err := qrcode.WriteFile("http://47.108.134.136:9000/findTeaByID?id="+uuid, qrcode.Medium, 256, wholePath)
				fmt.Println(err)

			}
		}

		actionTime := utils.SwitchTimeStampToData(time.Now().Unix())
		action := model.Action{
			UserID:           data.Sess.UserID,
			ActionType:       "批量上链",
			ActionTargetName: data.Tea.Name,
			ActionTargetID:   "111",
			ActionTime:       actionTime,
		}
		_, _ = dao.InsertAction(action)

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
		_ = json.Unmarshal(result, &tea)

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

	if data.IsStaff {

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
	} else if !data.IsStaff {
		data.Msg = "无权访问"
		ShowView(w, r, "index.html", data)
		return
	}
}
