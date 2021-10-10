package main

import (
	. "chatroom/client/config"
	"chatroom/client/process"
	"chatroom/public/tools"
	"fmt"
)

var filepath = "./chatroomClient.conf"

func main() {
	//Read Config
	err := ReadConfig(filepath)
	if err != nil {
		fmt.Println("ReadConfig \"%s\" err=", filepath, err)
		return
	}
	tools.MyLOG.Log("读取配置文件 %s 成功: %#v", filepath, MyConfig)

	process.ClientProc.MainInit()
	process.MyView.LoginView()

	tools.MyLOG.End()

}
