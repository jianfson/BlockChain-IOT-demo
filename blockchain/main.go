package main

import (
	"encoding/json"
	"fab-sdk-go-sample/sdkInit"
	"fab-sdk-go-sample/service"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"os"
)

const (
	configFile  = "config.yaml"
	initialized = false
	TeaCC       = "teacc"
)

func main() {
	sdkInit.SetupLogLevel()

	initInfo := &sdkInit.InitInfo{
		ChannelID:     "teatraceability",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/jianfson/BlockChain-IOT-demo/blockchain/fixtures/artifacts/channel.tx",

		OrgAdmin: "Admin",
		UserName: "User1",
		OrgName:  "Org1",

		OrdererName: "orderer.dragonwell.com",
		Peer:        "peer0.org1.dragonwell.com",

		ChaincodeID:     TeaCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/jianfson/BlockChain-IOT-demo/blockchain/chaincode/",
	}
	//-----------------------------------------
	//----------------实例化 sdk---------------
	//-----------------------------------------
	fmt.Println("----------------实例化 sdk---------------")
	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()
	//-----------------------------------------
	//------------------创建通道-----------------
	//-----------------------------------------
	fmt.Println("----------------创建通道---------------")
	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//-----------------------------------------
	//------------------加入通道-----------------
	//-----------------------------------------
	fmt.Println("----------------加入通道---------------")
	err = sdkInit.JoinChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//-----------------------------------------
	//-----------------查询通道信息--------------
	//-----------------------------------------
	fmt.Println("---------------查询通道信息---------------")
	clientChannelContext := sdk.ChannelContext(
		initInfo.ChannelID,
		fabsdk.WithUser(initInfo.OrgAdmin),
		fabsdk.WithOrg(initInfo.OrgName),
	)

	ledgerClient, err := ledger.New(clientChannelContext)

	if err != nil {
		fmt.Printf("Failed to create channel [%s] client: %#v", initInfo.ChannelID, err)
	}

	sdkInit.QueryChannelInfo(ledgerClient)
	sdkInit.QueryChannelConfig(ledgerClient)
	//----------------------------------------- -------------------------
	//-------------------------------安装链码------------------------------
	//-------------------------------------------------------------------
	fmt.Println("------------------安装链码-----------------")
	err = sdkInit.InstallCC(sdk, initInfo)
	if err != nil {
		fmt.Printf("InstallCC %v failed", initInfo.ChaincodeID)
	}

	sdkInit.QueryInstalledCC(sdk, initInfo)
	//----------------------------------------------------------------
	//-------------------------------实例化链码------------------------
	//-----------------------------------------------------------------
	fmt.Println("----------------实例化链码---------------")
	err = sdkInit.InstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Printf("实例化 %v failed", initInfo.ChaincodeID)
	}

	//-------------------------------------------------------------------
	//--- 创建一个通道客户端, 用来链码的调用、查询以及链码事件的注册和取消注册 ---
	//-------------------------------------------------------------------
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Printf("创建通道客户端失败: %v", err)
	}
	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")

	serviceSetup := service.ServiceSetup{
		ChaincodeId: TeaCC,
		Client:      channelClient,
	}

	teas := []service.Tea{
		{
			Id:     "01",
			Maker:  "dragonwell",
			Owner:  "dragonwell",
			Weight: "500",
		},
		{
			Id:     "02",
			Maker:  "dragonwell",
			Owner:  "wk",
			Weight: "500",
		},
		{
			Id:     "03",
			Maker:  "dragonwell",
			Owner:  "wk",
			Weight: "500",
		},
	}

	fmt.Println("----------------写入茶叶信息---------------")
	for k, tea := range teas {
		txID, err := serviceSetup.SaveTea(tea)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("%d信息保存成功\n交易Id：%v\n", k, txID)
		}
	}
	//
	fmt.Println("----------------查询茶叶信息---------------")
	b, err := serviceSetup.FindTeaInfoByID("01")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var tea service.Tea
		json.Unmarshal(b, &tea)
		fmt.Println("根据 teaID 查询信息成功：")
		fmt.Println(tea)
	}

	modifiedTea := service.Tea{
		Id:     "01",
		Maker:  "dragonwell",
		Owner:  "wk",
		Weight: "300",
	}

	fmt.Println("---------------- 修改茶叶信息 ---------------")
	txID, err := serviceSetup.ModifyTea(modifiedTea)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("修改成功\ntxid：%v\n", txID)
	}

	fmt.Println("----------------查询茶叶 \"01\" 信息---------------")
	teaId := "01"
	b, err = serviceSetup.FindTeaInfoByID(teaId)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var tea service.Tea
		json.Unmarshal(b, &tea)
		fmt.Println("根据 teaID 查询信息成功：")
		fmt.Printf("%v tea : %+v",teaId, tea)
	}

	fmt.Println("---------------- 查询茶叶 \"01\" 历史 ---------------")
	b, err = serviceSetup.GetHistoryForTea("01")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(b))
	}

	fmt.Println("---------------- 富查询茶叶信息 ---------------")
	b, err = serviceSetup.QueryTeaByString("wk")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(b))
	}

	fmt.Println("---------------- 删除茶叶信息 ---------------")
	b, err = serviceSetup.Delete("01")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("删除成功", string(b))
	}

	fmt.Println("----------------查询茶叶信息---------------")
	b, err = serviceSetup.FindTeaInfoByID("01")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var tea service.Tea
		json.Unmarshal(b, &tea)
		fmt.Println("根据 teaID 查询信息成功：")
		fmt.Println(tea)
	}

	fmt.Println("----------------登记、注册用户---------------")
	enrollSecret, err := sdkInit.Register(sdk, initInfo, "User2")
	if err != nil {
		fmt.Println(err)
	} else {
		err = sdkInit.Enroll(sdk, initInfo, enrollSecret)
	}

	sdkInit.GetUserInfo(sdk, "User2", "Org1")
}
