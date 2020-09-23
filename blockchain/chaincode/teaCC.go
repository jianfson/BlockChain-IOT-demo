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
		//buffer.WriteString("{\"Key\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(queryResponse.Key)
		//buffer.WriteString("\"")

		//buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("- getTeaByQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func SwitchTimeStampToData(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

func (s *TeaChaincode) initLedger(stub shim.ChaincodeStubInterface) pb.Response {
	shelf_life := "18个月"
	createtime := SwitchTimeStampToData(time.Now().Unix())
	teas := []Tea{
		{
<<<<<<< HEAD
			Id:     "01",
			Name: "高县红茶1",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井1",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "big",
<<<<<<< HEAD
			Boxed: Box{true,4},
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶2",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井2",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "big",
<<<<<<< HEAD
			Boxed: Box{true,4},
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶3",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井3",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "big",
<<<<<<< HEAD
			Boxed: Box{true,4},
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶4",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井4",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "big",
<<<<<<< HEAD
			Boxed: Box{true,4},
		},
		{
			Id:    strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶5",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:    strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井5",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "big",
<<<<<<< HEAD
			Boxed: Box{true,4},
		},
		{
			Id:    strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶6",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:    strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井6",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "big",
<<<<<<< HEAD
			Boxed: Box{},
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶7",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井7",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "small",
<<<<<<< HEAD
			Boxed: Box{},
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶8",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井8",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "small",
<<<<<<< HEAD
			Boxed: Box{},
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶11",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井11",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "small",
<<<<<<< HEAD
			Boxed: Box{},
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶12",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井12",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "small",
<<<<<<< HEAD
			Boxed: Box{},
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶11",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井11",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "small",
<<<<<<< HEAD
			Boxed: Box{},
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "高县红茶12",
			Maker:  "高县红茶茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "高县",
			Origin_IP: IP{
				"104",
				"23",
=======
		},
		{
			Id:     strconv.FormatInt(time.Now().UnixNano(),10),
			Name: "明前龙井12",
			Maker:  "杭州龙井茶业集团有限公司",
			Owner:  "王坤",
			Weight: "100g",
			Origin: "狮峰",
			Origin_IP: IP{
				"104.07",
				"30.67",
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
			},
			Production_Date: createtime,
			Shelf_life: shelf_life,
			Size: "small",
<<<<<<< HEAD
			Boxed: Box{},
=======
>>>>>>> 92729cb2df8df410654ee33e2f2af5c4041514c0
		},
	}
	for k, tea := range teas {
		tea.ObjectType = DOC_TYPE
		tea.TxID = stub.GetTxID()
		flag := PutTea(stub, tea)
		if flag != true {
			fmt.Println("写入第 %d 条信息，失败", k)
		}
		//time.After(time.Second)
	}

	return shim.Success(nil)
}

// args[0]: teaObj, args[1]: eventName; eventName 用于区分事件
func (s *TeaChaincode) saveTea(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
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

	tea.ObjectType = DOC_TYPE
	tea.TxID = stub.GetTxID()
	flag := PutTea(stub, tea)
	if !flag {
		return shim.Error("Add data failed")
	}

	return shim.Success([]byte("Add Tea succeed"))
}


func (s *TeaChaincode) teaExchange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect numbers of args, expecting 2")
	}

	teaID := args[0]
	nextOwner := args[1]
	tea,_ := GetTeaInfo(stub, teaID)
	tea.Owner = nextOwner
	tea.TxID = stub.GetTxID()
	flag := PutTea(stub, tea)
	if !flag {
		return shim.Error("Save data failed")
	}

	return shim.Success([]byte("exchange succeed"))
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

func (s *TeaChaincode) modifyQueryCount(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect numbers of args, expecting 1")
	}

	teaID := args[0]

	tea,_ := GetTeaInfo(stub, teaID)
	tea.QueryCounter++

	flag := PutTea(stub, tea)
	if !flag {
		return shim.Error("Save data failed")
	}

	return shim.Success([]byte("modify succeed"))
}
func (s *TeaChaincode) queryTeaByMaker(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of args, expecting 1")
	}

	// 拼接富查询用到的 string
	queryString := fmt.Sprintf("{\"selector\":{\"maker\":\"%s\"}}", args[0])
	results, err := getTeaByQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	if results == nil {
		return shim.Error("get nothing according to weight and maker")
	}
	return shim.Success(results)
}

func (s *TeaChaincode) getTeasByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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

	fmt.Printf("- getTeasByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
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
func (s *TeaChaincode) deleteTea(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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
