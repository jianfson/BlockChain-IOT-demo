package main

import (
	"encoding/json"

	"fab-sdk-go-sample/sdkInit"
	"fab-sdk-go-sample/service"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"os"
	"strconv"
	"time"
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
	fmt.Println("------------------查询已安装链码-----------------")
	sdkInit.QueryInstalledCC(sdk, initInfo)

	//----------------------------------------------------------------
	//-------------------------------实例化链码------------------------
	//-----------------------------------------------------------------
	fmt.Println("----------------实例化链码---------------")
	err = sdkInit.InstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Printf("实例化 %v failed, err:%v", initInfo.ChaincodeID, err)
		return
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

	fmt.Println("----------------初始化茶叶信息---------------")
	shelf_life := "18个月"
	t := time.Unix(time.Now().Unix(), 0)
	createtime := t.Format("2006-01-02 15:04:05")

	teas := []service.Tea{
		{
			Id:     "01",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "02",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "03",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "04",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "05",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "06",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "07",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "08",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "09",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "10",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "11",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id:     "12",
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
		{
			Id: strconv.FormatInt(time.Now().Unix(),10),
			Name: "明前龙井",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Production_Date: createtime,
			Shelf_life: shelf_life,
		},
	}

	for _, tea := range teas {
		serviceSetup.SaveTea(tea)
	}
	//
	fmt.Println("----------------查询茶叶信息---------------")
	b, err := serviceSetup.QueryTeaByMaker("杭州龙井茶业集团有限公司")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var tea []service.Tea
		json.Unmarshal(b, &tea)
		fmt.Printf("%+v", tea)
	}
	//enrollmentSecret, err := sdkInit.Register(sdk,initInfo,"user2")
	//log.Println("enrollmentSecret:", enrollmentSecret)
	//if err != nil {
	//	log.Println(err)
	//} else {
	//	err = sdkInit.Enroll(sdk,initInfo,enrollmentSecret)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}
	//sdkInit.GetUserInfo(sdk, "user2", "Org1")
}
