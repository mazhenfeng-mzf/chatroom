package main

import (
	"chatroom/public/tools"
	. "chatroom/server/config"
	serverProcess "chatroom/server/process"
	"fmt"
	"net"
)

var filepath = "./chatroomServer.conf"

func main() {

	tools.MyLOG.Init(tools.LogServer)

	//Read Config
	err := ReadConfig(filepath)
	if err != nil {
		fmt.Println("读取配置文件 \"%s\" 失败, err=", filepath, err)
		return
	}
	tools.MyLOG.Log("读取配置文件 %s 成功: %#v", filepath, MyConfig)

	fmt.Println("server start listening")
	listen, err := net.Listen("tcp", "0.0.0.0:"+MyConfig.LISTEN_PORT)
	if err != nil {
		fmt.Println("listen err=", err)
		return
	}
	//fmt.Printf("listen=%#v", listen)
	defer listen.Close()

	serverProcess.Connect_DB()
	serverProcess.CliMgr.Init()

	//waiting remote client
	for {
		fmt.Println("waitting client1")
		conn1, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept() err=", err)
		} else {
			fmt.Printf("Accept() succ, client ip=%v\n", conn1.RemoteAddr().String())
		}

		go serverProcess.FindProcess(conn1)

	}

}
