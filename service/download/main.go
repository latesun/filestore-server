package main

import (
	"time"

	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"

	"filestore-server/common"
	dbproxy "filestore-server/service/dbproxy/client"
	cfg "filestore-server/service/download/config"
	dlProto "filestore-server/service/download/proto"
	"filestore-server/service/download/route"
	dlRpc "filestore-server/service/download/rpc"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.download"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	if err := dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlRpc.Download)); err != nil {
		panic(err)
	}
	if err := service.Run(); err != nil {
		panic(err)
	}
}

func startAPIService() {
	router := route.Router()
	panic(router.Run(cfg.DownloadServiceHost))
}

func main() {
	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
