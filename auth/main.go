package main

import (
	"e_commerce/auth/config"
	"e_commerce/auth/db"
	"e_commerce/kitex_gen/auth/authservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

func Init() {
	db.InitDB(config.MySQLDSN)
	db.InitRedis(config.RedisAddr, config.RedisPassword)
}

func main() {
	Init()

	addr, _ := net.ResolveTCPAddr("tcp", config.ServerAddr)

	svr := authservice.NewServer(new(AuthServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "auth"}),
		server.WithServiceAddr(addr),
	)

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
