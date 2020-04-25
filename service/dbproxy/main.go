package main

import (
	"log"
	"time"

	"filestore-server/common"
	"filestore-server/service/dbproxy/config"
	dbConn "filestore-server/service/dbproxy/conn"
	dbProxy "filestore-server/service/dbproxy/proto"
	dbRpc "filestore-server/service/dbproxy/rpc"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
)

func startRpcService() {
	service := micro.NewService(
		micro.Name("go.micro.service.dbproxy"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),      // 声明超时时间, 避免consul不主动删掉已失去心跳的服务节点
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			// 检查是否指定dbhost
			dbhost := c.String("dbhost")
			if len(dbhost) > 0 {
				log.Println("custom db address: " + dbhost)
				config.UpdateDBHost(dbhost)
			}
		}),
	)

	// 初始化db connection
	dbConn.InitDBConn()

	if err := dbProxy.RegisterDBProxyServiceHandler(service.Server(), new(dbRpc.DBProxy)); err != nil {
		log.Println(err)
	}

	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

func main() {
	startRpcService()
}

// res, err := mapper.FuncCall("/user/UserExist", []interface{}{"haha"}...)
// log.Printf("error: %+v\n", err)
// log.Printf("result: %+v\n", res[0].Interface())

// res, err = mapper.FuncCall("/user/UserExist", []interface{}{"admin"}...)
// log.Printf("error: %+v\n", err)
// log.Printf("result: %+v\n", res[0].Interface())
