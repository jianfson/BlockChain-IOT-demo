package cliInit

import (
	"blc-iot-demo/web/service"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
)

const (
	ConfigFile = "./blockchain/config.yaml"

	ChannelID = "teatraceability"

	Org  = "Org1"
	User = "Admin"
)

var SDK *fabsdk.FabricSDK

func CliInit() *service.ServiceSetup {
	// sdk 实例化
	SDK, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		log.Printf("实例化 sdk 失败: %v", err)
		return nil
	}

	//-------------------------Channel Client-------------------------------------
	//	 client context, 用于实例化
	clientChannelContext := SDK.ChannelContext(
		ChannelID,
		fabsdk.WithUser(User),
		fabsdk.WithOrg(Org),
	)

	// channel cli 实例化
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		log.Printf("创建通道客户端失败: %v", err)
		return nil
	}

	//-------------------------Ledger Client-------------------------------------
	clientLedgerContext := SDK.ChannelContext(
		ChannelID,
		fabsdk.WithUser(User),
		fabsdk.WithOrg(Org),
	)

	ledgerClient, err := ledger.New(clientLedgerContext)
	if err != nil {
		log.Printf("创建账本客户端失败: %v", err)
		return nil
	}

	return &service.ServiceSetup{
		ChaincodeId:   "teacc",
		ChannelClient: channelClient,
		LedgerClient:  ledgerClient,
	}
}
