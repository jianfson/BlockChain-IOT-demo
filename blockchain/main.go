package main

import (
	"fab-sdk-go-sample/sdkInit"
	"fab-sdk-go-sample/service"
	"fab-sdk-go-sample/web"
	"fab-sdk-go-sample/web/controller"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
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

	clientChannelContext := sdk.ChannelContext(
		initInfo.ChannelID,
		fabsdk.WithUser(initInfo.OrgAdmin),
		fabsdk.WithOrg(initInfo.OrgName),
	)

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
		ChaincodeId:   TeaCC,
		ChannelClient: channelClient,
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
			fmt.Printf("%d 信息保存成功\n交易Id：%v\n", k, txID)
		}
	}
	app := controller.Application{
		Setup: &serviceSetup,
	}
	web.WebStart(&app)
}
