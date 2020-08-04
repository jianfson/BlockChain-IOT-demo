package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"

)

// 初始化账本
func (t *ServiceSetup) InitLedger() error {

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "initLedger", Args: [][]byte{}}
	_, err := t.ChannelClient.Execute(req)
	if err != nil {
		return err
	}
	return nil
}
// 调用链码向账本添加 tea, 返回一个　TX id
func (t *ServiceSetup) SaveTea(tea Tea) (string, error) {

	b, err := json.Marshal(tea)
	if err != nil {
		return "", fmt.Errorf("序列化 tea 失败")
	}


	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "saveTea", Args: [][]byte{b}}


	// the proposal responses from peer(s)

	response, err := t.ChannelClient.Execute(req)
	if err != nil {
		return "", err
	}


	return string(response.TransactionID), nil
}

// 修改 tea 信息
func (t *ServiceSetup) ModifyTea(teaID, nextOwner string) (string, error) {


	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "teaExchange", Args: [][]byte{[]byte(teaID), []byte(nextOwner)}}
	respone, err := t.ChannelClient.Execute(req)

	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

// 通过 teaID 查询
func (t *ServiceSetup) FindTeaInfoByID(teaID string) ([]byte, error) {

	req := channel.Request{
		ChaincodeID: t.ChaincodeId,
		Fcn: "queryTeaById",
		Args: [][]byte{[]byte(teaID)},
	}

	respone, err := t.ChannelClient.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 通过 s 查询
func (t *ServiceSetup) QueryTeaByMaker(s string) ([]byte, error) {


	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "queryTeaByMaker", Args: [][]byte{[]byte(s)}}

	respone, err := t.ChannelClient.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// query by range
func (t *ServiceSetup) GetTeasByRange(start, end string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "getTeasByRange", Args: [][]byte{[]byte(start), []byte(end)}}
	respone, err := t.ChannelClient.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 通过 id 删除
func (t *ServiceSetup) Delete(teaId string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "delete", Args: [][]byte{[]byte(teaId)}}
	respone, err := t.ChannelClient.Execute(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 通过 id 查询历史
func (t *ServiceSetup) GetHistoryForTea(teaId string) ([]byte, error) {

	req := channel.Request{ChaincodeID: t.ChaincodeId, Fcn: "getHistoryForTea", Args: [][]byte{[]byte(teaId)}}
	respone, err := t.ChannelClient.Query(req)
	if err != nil {
		return []byte{0x00}, err
	}

	return respone.Payload, nil
}

// 通过 txId 查询交易
func (t *ServiceSetup) QueryTransactionByTxID(txId string) ([]byte, error) {
	req, err := t.LedgerClient.QueryTransaction(fab.TransactionID(txId))
	if err != nil {
		return []byte{0x00}, err
	}
	return req.TransactionEnvelope.Payload, nil
}

// 通过 txId 查询区块
func (t *ServiceSetup) QueryBlockByTxID(txId string) (*common.Block, error) {
	req, err := t.LedgerClient.QueryBlockByTxID(fab.TransactionID(txId))
	if err != nil {
		return nil, err
	}
	return req, nil
}

// 通过 区块高度 查询区块
func (t *ServiceSetup) QueryBlockByNum(num uint64) (*common.Block, error) {
	req, err := t.LedgerClient.QueryBlock(num)
	if err != nil {
		return nil, err
	}
	return req, nil
}


