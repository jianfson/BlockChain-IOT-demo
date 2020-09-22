package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type TeaChaincode struct {
}

func main() {
	// Create a new Smart Contract
	err := shim.Start(new(TeaChaincode))
	if err != nil {
		fmt.Printf("Error starting Tea chaincode: %s", err)
	}
}

// 实现 Init 方法, 实例化账本时使用。
func (s *TeaChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	s.initLedger(stub)
	return shim.Success(nil)
}

// 实现 Invoke 方法
func (s *TeaChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// 获取函数名称、参数
	fn, args := stub.GetFunctionAndParameters()

	//调用对应函数
	if fn == "saveTea" {
		return s.saveTea(stub, args)

	}else if fn == "teaExchange" {
		return s.teaExchange(stub, args)

	} else if fn == "queryTeaById" {
		return s.queryTeaById(stub, args)

	} else if fn == "modifyQueryCount" {
		return s.modifyQueryCount(stub, args)

	}else if fn == "queryTeaByMaker" {
		return s.queryTeaByMaker(stub, args)

	}else if fn == "getTeasByRange" {
		return s.getTeasByRange(stub, args)

	} else if fn == "getHistoryForTea" {
		return s.getHistoryForTea(stub, args)


	} else if fn == "deleteTea" {
		return s.deleteTea(stub, args)

	} else if fn == "getTeasByRange" {
		return s.getTeasByRange(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}
