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

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "addTea", Args: [][]byte{b, []byte(eventId)}}
	// the proposal responses from peer(s)

	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventId)
	if err != nil {
		return "", err
	}
	return string(response.TransactionID), nil
}

// 修改 tea 信息
func (t *ServiceSetup) ModifyTea(tea Tea) (string, error) {

	eventID := "eventModifyTea"
	reg, notifier := registerEvent(t.Client, t.ChaincodeId, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	b, err := json.Marshal(tea)
	if err != nil {
		return "", fmt.Errorf("指定的edu对象序列化时发生错误")
	}

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "updateTea", Args: [][]byte{b, []byte(eventID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

// 通过 teaID 查询
func (t *ServiceSetup) FindTeaInfoByID(teaID string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "queryTeaById", Args: [][]byte{[]byte(teaID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 通过 s 查询
func (t *ServiceSetup) QueryTeaByString(s string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "queryTeaByString", Args: [][]byte{[]byte(s)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 通过 id 删除
func (t *ServiceSetup) Delete(teaId string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "delete", Args: [][]byte{[]byte(teaId)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 通过 id 查询历史
func (t *ServiceSetup) GetHistoryForTea(teaId string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "getHistoryForTea", Args: [][]byte{[]byte(teaId)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}