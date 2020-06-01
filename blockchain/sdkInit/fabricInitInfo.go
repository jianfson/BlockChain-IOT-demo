/* 为了方便管理 Hyperledger Fabric 网络环境，我们将在sdkInit 目录中
创建一个 fabricInitInfo.go 的源代码文件，用于定义一个结构体，包括 Fabric SDK 所需的各项相关信息
 */
package sdkInit

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

type InitInfo struct {
	ChannelID     string   // 通道 id
	ChannelConfig string   // 通道配置文件
	OrgName      string    // 组织名称
	OrgAdmin       string  //已组织管理员身份创建
	OrdererOrgName    string  // order
	OrgResMgmt *resmgmt.Client

	ChaincodeID    string
	ChaincodeGoPath    string
	ChaincodePath    string
	UserName    string
}