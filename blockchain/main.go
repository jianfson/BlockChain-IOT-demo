package main

import (
	"encoding/json"
	"fab-sdk-go-sample/service"
	"os"
	"fmt"
	"fab-sdk-go-sample/sdkInit"
)

const (
	configFile = "config.yaml"
	initialized = false
	TeaCC = "teacc"
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID: "kevinkongyixueyuan",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/hyperledger/fabric-sdk-go-sample/fixtures/artifacts/channel.tx",

		OrgAdmin:"Admin",
		OrgName:"Org1",
		OrdererOrgName: "orderer.kevin.kongyixueyuan.com",

		ChaincodeID: TeaCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath: "github.com/hyperledger/fabric-sdk-go-sample/chaincode/",
		UserName:"User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()

	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	serviceSetup := service.ServiceSetup{
		ChaincodeId:TeaCC,
		Client: channelClient,
	}

	tea := service.Tea{
		Id:      "01",
		Maker:    "龙井",
		Owner:   "龙井",
		Weight:  "500g",
	}

	txID, err := serviceSetup.SaveTea(tea)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("信息保存成功\n交易Id：%v\n", txID)
	}

	id := "01"
	b, err := serviceSetup.FindTeaInfoByID(id)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var tea service.Tea
		json.Unmarshal(b, &tea)
		fmt.Println("根据 teaID 查询信息成功：")
		fmt.Println(tea)
	}

	modifiedTea := service.Tea{
		Id:      "01",
		Maker:    "龙井",
		Owner:   "wk",
		Weight:  "500g",
	}

	txID, err = serviceSetup.ModifyTea(modifiedTea)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("修改成功\ntxid：%v\n", txID)
	}

	b, err = serviceSetup.FindTeaInfoByID(id)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var tea service.Tea
		json.Unmarshal(b, &tea)
		fmt.Println("根据 teaID 查询信息成功：")
		fmt.Println(tea)
		for _, v := range tea.Histories {
			fmt.Println(v)
		}
	}
}