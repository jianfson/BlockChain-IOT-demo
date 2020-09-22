package main

import (
	"blc/sdkInit"
	"fmt"
	"os"
)

const (
	configFile  = "./config.yaml"
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

		OrdererName: "orderer.goldenBud.com",
		Peer:        "peer0.org1.goldenBud.com",

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


	//----------------------------------------- -------------------------
	//-------------------------------安装链码------------------------------
	//-------------------------------------------------------------------
	fmt.Println("------------------安装链码-----------------")
	err = sdkInit.InstallCC(sdk, initInfo)
	if err != nil {
		fmt.Printf("InstallCC %v failed", initInfo.ChaincodeID)
	}
	fmt.Println("------------------查询已安装链码-----------------")
	sdkInit.QueryInstalledCC(sdk, initInfo)

	//----------------------------------------------------------------
	//-------------------------------实例化链码------------------------
	//-----------------------------------------------------------------
	fmt.Println("----------------实例化链码---------------")
	err = sdkInit.InstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
}
