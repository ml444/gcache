package gcache

import (
	"flag"
	"fmt"

	"github.com/ml444/gcache/broker"
	"github.com/ml444/gcache/config"
)

var joinToken string
var configPath string
var address string

func init() {
	flag.StringVar(&joinToken, "join", "", "join token for cluster")
	flag.StringVar(&configPath, "conf", "", "config file path")
	flag.StringVar(&address, "addr", ":6389", "server address")
}

func main() {
	flag.Parse()
	var cfg *config.GroupConfig
	if configPath != "" {
		cfg = config.LoadConfig(configPath)
	} else {
		cfg = config.DefaultConfig()
	}

	if joinToken == "" {
		// create new cluster
		// generate join token
		joinToken = "123456" // TODO: generate join token
		fmt.Println("NewGroup cluster, token: ", joinToken)
		// 启动broker
		b := broker.NewBroker(address)
		// 初始化默认的group
		shardSerialNoList := make([]int, cfg.ShardCount)
		for i := 0; i < cfg.ShardCount; i++ {
			shardSerialNoList[i] = i
		}
		b.AddGroup(broker.NewGroup(cfg, address, shardSerialNoList))

		// 初始化全局metadata
	}
	Run()
}

func Run() {

	// 启动server listen, 等待Join信息

	// 启动broker
	// 获取全局metadata
	// 初始化属于自己的group

	// block server listen
}

/*
POST Join
DELETE Leave
POST AddGroup
DELETE RemoveGroup
POST AddGroupNode
DELETE RemoveGroupNode
Put UpdateGroupMetadata

GET Stats
POST SetCallback
*/
