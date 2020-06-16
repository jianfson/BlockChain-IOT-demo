package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

const DOC_TYPE = "teaObj"

// 该文件实现使用链码相关API对账本状态进行具体操作的函数们

// PutTea() 将对象序列化后保存至账本中
func PutTea(stub shim.ChaincodeStubInterface, tea Tea) bool {

	//tea.ObjectType = DOC_TYPE
	b, err := json.Marshal(tea)
	if err != nil {
		return false
	}

	err = stub.PutState(tea.Id, b)
	if err != nil {
		return false
	}

	return true
}

// GetTeaInfo() 根据指定的茶叶 Id 查询对应的状态，反序列化后并返回对象
func GetTeaInfo(stub shim.ChaincodeStubInterface, entityId string) (Tea, bool) {

	var tea Tea

	b, err := stub.GetState(entityId)
	if err != nil || b == nil {
		return tea, false
	} // 有错误 或者 Id不存在[id不存在GetState()返回 nil, nil]

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
}

// args[0]: teaObj, args[1]: eventName; eventName 用于区分事件
func (s *TeaChaincode) addTea(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	var tea Tea
	err := json.Unmarshal([]byte(args[0]), &tea)
	if err != nil {
		return shim.Error("Unmarshal tea failed")
	}

	// ==== Input sanitation ====
	if len(tea.Id) <= 0 {
		return shim.Error("Id must be a non-empty string")
	}
	if len(tea.Maker) <= 0 {
		return shim.Error("Maker must be a non-empty string")
	}
	if len(tea.Owner) <= 0 {
		return shim.Error("Owner must be a non-empty string")
	}
	if len(tea.Weight) <= 0 {
		return shim.Error("Weight must be a non-empty string")
	}

	_, exist := GetTeaInfo(stub, tea.Id)
	if exist {
		return shim.Error("Id specified already exists")
	}

	tea.ObjectType = DOC_TYPE
	flag := PutTea(stub, tea)
	if !flag {
		return shim.Error("Add data failed")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("Add Tea succeed"))
}

// args[0]: teaObj, args[1]: eventName;
func (s *TeaChaincode) updateTea(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect numbers of args, expecting 2")
	}

	var tea Tea
	err := json.Unmarshal([]byte(args[0]), &tea)
	if err != nil {
		return shim.Error("Unmarshal tea failed")
	}

	if len(tea.Id) <= 0 {
		return shim.Error("Id must be a non-empty string")
	}
	if len(tea.Maker) <= 0 {
		return shim.Error("Maker must be a non-empty string")
	}
	if len(tea.Owner) <= 0 {
		return shim.Error("Owner must be a non-empty string")
	}
	if len(tea.Weight) <= 0 {
		return shim.Error("Weight must be a non-empty string")
	}

	tea.ObjectType = DOC_TYPE
	flag := PutTea(stub, tea)
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
		return shim.Error("incorrect nums of  args, expecting 1")
	}

	result, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("query failed according to id")
	}
	if result == nil {
		return shim.Error("get nothing according to id")
	}
	return shim.Success(result)
}

func (s *TeaChaincode) queryTeaByString(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of args, expecting 1")
	}

	// 拼接富查询用到的 string
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"Owner\":\"%s\"}}", DOC_TYPE, args[0])
	results, err := getTeaByQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	if results == nil {
		return shim.Error("get nothing according to weight and maker")
	}
	return shim.Success(results)
}

func (s *TeaChaincode) getHistoryForTea(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	teaId := args[0]

	fmt.Printf("- start getHistoryForMarble: %s\n", teaId)

	resultsIterator, err := stub.GetHistoryForKey(teaId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ==================================================
// delete - remove a tea key/value pair from state
// ==================================================
func (s *TeaChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	teaId := args[0]

	valAsbytes, err := stub.GetState(teaId) //get the tea from chaincode state
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + teaId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"Tea does not exist: " + teaId + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(teaId) //remove the tea from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	return shim.Success(nil)
}