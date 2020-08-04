/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

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
	return shim.Success(nil)
}

// 实现 Invoke 方法
func (s *TeaChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// 获取函数名称、参数
	fn, args := stub.GetFunctionAndParameters()

	//调用对应函数
	if fn == "saveTea" {
		return s.saveTea(stub, args)

	} else if fn == "initLedger" {
		return s.initLedger(stub, args)

	}else if fn == "teaExchange" {
		return s.teaExchange(stub, args)

	} else if fn == "queryTeaById" {
		return s.queryTeaById(stub, args)

	} else if fn == "queryTeaByMaker" {
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
