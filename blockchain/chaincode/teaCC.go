package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"
)

const DOC_TYPE = "teaObj"

// 该文件实现使用链码相关API对账本状态进行具体操作的函数们

// PutTea() 将对象序列化后保存至账本中
func PutTea(stub shim.ChaincodeStubInterface, tea Tea) ([]byte, bool) {

	tea.ObjectType = "teaObj"
	b, err := json.Marshal(tea)
	if err != nil {
		return nil, false
	}

	err = stub.PutState(tea.Id, b)
	if err != nil {
		return nil, false
	}

	return b, true
}

// GetTeaInfo() 根据指定的茶叶 Id 查询对应的状态，反序列化后并返回对象
func GetTeaInfo(stub shim.ChaincodeStubInterface, entityId string) (Tea, bool) {

	var tea Tea

	b, err := stub.GetState(entityId)
	if err != nil || b == nil {
		return tea, false
	} 	// 有错误 或者 Id不存在[id不存在GetState()返回 nil, nil]

	err = json.Unmarshal(b, &tea)
	if err != nil {
		return tea, false
	}

	return tea, true
}

// getTeaByQueryString() 根据指定的字符串进行富查询
func getTeaByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getTeaByQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getTeaByQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil


	// 将查询结果从 resultIterator 提取
	//hasComma := false
	//for resultIterator.HasNext() {
	//	queryResponse, err := resultIterator.Next()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	if hasComma {
	//		buffer.WriteString(",")
	//	}
	//
	//	buffer.WriteString(string(queryResponse.Value))
	//	hasComma = true
	//}
	//
	//return buffer.Bytes(), nil
}

// args[0]: teaObj, args[1]: eventName; eventName 用于区分事件
func (s *TeaChaincode) addTea(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("args not enough")
	}

	var tea Tea
	err := json.Unmarshal([]byte(args[0]), &tea)
	if err != nil {
		return shim.Error("Unmarshal tea failed")
	}

	_, exist := GetTeaInfo(stub, tea.Id)
	if exist {
		return shim.Error("Id specified already exists")
	}

	_, flag := PutTea(stub, tea)
	if !flag {
		return shim.Error("Save data failed")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("add Tea succeed"))
}

// args[0]: teaObj, args[1]: eventName;
func (s *TeaChaincode) updateTea(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("args not enough")
	}

	var tea Tea
	err := json.Unmarshal([]byte(args[0]), &tea)
	if err != nil {
		return shim.Error("Unmarshal tea failed")
	}

	result, flag := GetTeaInfo(stub, tea.Id)
	if !flag {
		return shim.Error("query falied according to Id specified ")
	}

	result.Owner = tea.Owner
	_, flag =PutTea(stub, result)
	if !flag {
		return shim.Error("Save data failed")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("updata succeed"))
}

func (s *TeaChaincode) queryTeaById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("incorrect nums of  args, expect 1" )
	}

	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("query failed according to id")
	}
	if result == nil {
		return shim.Error("get nothing according to id")
	}

	var tea Tea
	err = json.Unmarshal(result, &tea)
	if err != nil {
		return shim.Error("unmarshal tea failed ")
	}
	// 获取历史数据
	iterator, err := stub.GetHistoryForKey(tea.Id)
	if err != nil{
		return shim.Error("get history data failed")
	}
	defer iterator.Close()

	var hisTea Tea
	var histories []HistoryItem
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("err when get history data ")
		}

		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisTea)
		if hisData.Value == nil {
			var empty Tea
			historyItem.tea = empty
		} else {
			historyItem.tea = hisTea
		}

		histories = append(histories, historyItem)
	}
	tea.Histories = histories
	b, err := json.Marshal(tea)
	if err != nil {
		return shim.Error("err when marshal Tea")
	}
	return shim.Success(b)
}

func (s *TeaChaincode) queryTeaByString(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of args, expecting 1")
	}

	owner := strings.ToLower(args[0])

	// 拼接富查询用到的
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%v\",\"owner\":\"%s\"}}", DOC_TYPE,owner)
	results, err := getTeaByQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	if results == nil {
		return shim.Error("get nothing according to weight and maker")
	}
	return shim.Success(results)
}