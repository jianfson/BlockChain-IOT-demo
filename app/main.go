package main

import (
	"encoding/json"
	"fmt"

	"context"
	"time"

	"net/http"

	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

func main() {
	router := gin.Default()

	// 定义路由
	{
		router.POST("/teas", saveTea)
		router.GET("/teas/:id", queryTea)
		router.DELETE("/teas/:id", deleteTea)
		router.POST("/teas/exchange", teaExchange)
		//router.GET("/teas/exchange/history", teaExchangeHistory)
	}

	router.Run() // listen and serve on 0.0.0.0:8080
}

type Tea struct {
	Id     string `form:"id" binding:"required"`
	Maker  string `form:"maker" binding:"required"`
	Owner  string `form:"owner" binding:"required"`
	Weight string `form:"weight" binding:"required"`
}

// 茶叶注册
func saveTea(ctx *gin.Context) {
	// 参数处理
	req := new(Tea)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	// 区块链交互
	b, err := json.Marshal(*req)
	resp, err := channelExecute("saveTea", [][]byte{
		b,
		[]byte("event_saveTea"),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 查询茶叶
func queryTea(ctx *gin.Context) {
	teaId := ctx.Param("id")

	resp, err := channelQuery("queryTeaById", [][]byte{
		[]byte(teaId),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	//ctx.JSON(http.StatusOK, resp)
	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}

// 删除茶叶
func deleteTea(ctx *gin.Context) {
	teaId := ctx.Param("id")

	resp, err := channelExecute("deleteTea", [][]byte{
		[]byte(teaId),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

type TeaExchangeRequest struct {
	TeaId        string `form:"teaid" binding:"required"`
	NextOwnerId string `form:"nextownerid" binding:"required"`
}

// 资产转让
func teaExchange(ctx *gin.Context) {
	req := new(TeaExchangeRequest)
	if err := ctx.ShouldBind(req); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	resp, err := channelExecute("teaExchange", [][]byte{
		[]byte(req.TeaId),
		[]byte(req.NextOwnerId),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// 资产历史变更记录
func teaExchangeHistory(ctx *gin.Context) {
	assetId := ctx.Query("assetid")
	queryType := ctx.Query("querytype")

	resp, err := channelQuery("queryAssetHistory", [][]byte{
		[]byte(assetId),
		[]byte(queryType),
	})

	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	//ctx.JSON(http.StatusOK, resp)
	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}

var (
	sdk           *fabsdk.FabricSDK
	channelName   = "teatraceability"
	chaincodeName = "teacc"
	org           = "org1"
	user          = "Admin"
	//configPath    = "$GOPATH/src/github.com/hyperledger/fabric/imocc/application/config.yaml"
	configPath = "../blockchain/config.yaml"
)

func init() {
	var err error
	sdk, err = fabsdk.New(config.FromFile(configPath))
	if err != nil {
		panic(err)
	}
}

// 区块链管理
func manageBlockchain() {
	// 表明身份
	ctx := sdk.Context(fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := resmgmt.New(ctx)
	if err != nil {
		panic(err)
	}

	// 具体操作
	cli.SaveChannel(resmgmt.SaveChannelRequest{}, resmgmt.WithOrdererEndpoint("orderer.dragonwell.com"), resmgmt.WithTargetEndpoints())
}

// 区块链数据查询 账本的查询
func queryBlockchain() {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := ledger.New(ctx)
	if err != nil {
		panic(err)
	}

	resp, err := cli.QueryInfo(ledger.WithTargetEndpoints("peer0.org1.dragonwell.com"))
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	//// 1
	//cli.QueryBlockByHash(resp.BCI.CurrentBlockHash)
	//
	//// 2
	//for i := uint64(0); i <= resp.BCI.Height; i++ {
	//	cli.QueryBlock(i)
	//}
}

// 区块链交互
func channelExecute(fcn string, args [][]byte) (channel.Response, error) {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	// 状态更新，insert/update/delete
	resp, err := cli.Execute(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.org1.dragonwell.com"))
	if err != nil {
		return channel.Response{}, err
	}

	// 链码事件监听
	go func() {
		// channel
		reg, ccevt, err := cli.RegisterChaincodeEvent(chaincodeName, "eventname")
		if err != nil {
			return
		}
		defer cli.UnregisterChaincodeEvent(reg)

		timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		for {
			select {
			case evt := <-ccevt:
				fmt.Printf("received event of tx %s: %+v", resp.TransactionID, evt)
			case <-timeoutctx.Done():
				fmt.Println("event timeout, exit!")
				return
			}
		}

		// event
		//eventcli, err := event.New(ctx)
		//if err != nil {
		//	return
		//}

		//eventcli.RegisterChaincodeEvent(chaincodeName, "eventname")
	}()

	// 交易状态事件监听
	go func() {
		eventcli, err := event.New(ctx)
		if err != nil {
			return
		}

		reg, status, err := eventcli.RegisterTxStatusEvent(string(resp.TransactionID))
		defer eventcli.Unregister(reg) // 注册必有注销

		timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		for {
			select {
			case evt := <-status:
				fmt.Printf("received event of tx %s: %+v", resp.TransactionID, evt)
			case <-timeoutctx.Done():
				fmt.Println("event timeout, exit!")
				return
			}
		}
	}()

	return resp, nil
}

func channelQuery(fcn string, args [][]byte) (channel.Response, error) {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	// 状态的查询，select
	return cli.Query(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         fcn,
		Args:        args,
	}, channel.WithTargetEndpoints("peer0.org1.dragonwell.com"))
}

// 事件监听
func eventHandle() {
	ctx := sdk.ChannelContext(channelName, fabsdk.WithOrg(org), fabsdk.WithUser(user))

	cli, err := event.New(ctx)
	if err != nil {
		panic(err)
	}

	// 交易状态事件
	// 链码事件 业务事件
	// 区块事件
	reg, blkevent, err := cli.RegisterBlockEvent()
	if err != nil {
		panic(err)
	}
	defer cli.Unregister(reg)

	timeoutctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	for {
		select {
		case evt := <-blkevent:
			fmt.Printf("received a block", evt)
		case <-timeoutctx.Done():
			fmt.Println("event timeout, exit!")
			return
		}
	}
}
