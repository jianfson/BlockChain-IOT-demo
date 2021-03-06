package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)

// tea describes basic details of what makes up a tea
type Tea struct {
	ObjectType      string `json:"docType"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	Maker           string `json:"maker"`
	Owner           string `json:"owner"`
	Weight          string `json:"weight"`
	Origin          string `json:"origin"`
	Origin_IP       IP `json:"origin_ip"`
	Production_Date string `json:"production_date"`
	Shelf_life      string `json:"shelf_life"`
	TxID            string `json:"txID"`
	Size            string `json:"size"`
	QueryCounter    int    `json:"queryCounter"`
	Boxed           Box    `json:"boxed"`
}

// 成都市 经度:104.07 纬度:30.67
type IP struct {
	Longitude string
	Latitude  string
}
type Box struct {
	Boxed bool `json:"boxed"`
	Num   int  `json:"num"`
}
type Block struct {
	Height    uint64 `json:"height"`
	DataHash  string `json:"datahash"`
	BlcHash   string `json:"blchash"`
	Timestamp string `json:"timestamp"`
}

type ServiceSetup struct {
	ChaincodeId   string
	ChannelClient *channel.Client
	LedgerClient  *ledger.Client
}

//注册链码事件
func registerEvent(client *channel.Client, chaincodeId string, eventId string) (fab.Registration, <-chan *fab.CCEvent) {
	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeId, eventId)
	if err != nil {
		fmt.Printf("注册链码事件发生错误：%s", err)
	}
	return reg, notifier
}

// 执行链码完成后成功了吗?
func eventResult(notifier <-chan *fab.CCEvent, eventId string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件：%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能接受到链码事件：%s\n", eventId)
	}
	return nil
}
