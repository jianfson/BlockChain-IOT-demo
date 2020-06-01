package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// 调用链码向账本添加 tea, 返回一个　TX id
func (t *ServiceSetup) SaveTea(tea Tea) (string, error) {
	fmt.Println("开始写入账本")
	eventId := "event_addTea"
	reg, notifier := registerEvent(t.Client, t.ChaincodeId, eventId)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(tea)
	if err != nil {
		return "", fmt.Errorf("序列化 tea 失败, eventId：%v\n", eventId)
	}

	req := channel.Request{ChaincodeID:t.ChaincodeId, Fcn:"addTea",Args:[][]byte{b, []byte(eventId)}}
	// the proposal responses from peer(s)

	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventId)
	if  err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}