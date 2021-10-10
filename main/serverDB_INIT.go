package main

import (
	. "chatroom/server/config"
	"chatroom/server/model"
	"fmt"
)

var filepath = "./chatroomServer.conf"

func main() {

	//Read Config
	err := ReadConfig(filepath)
	if err != nil {
		fmt.Println("读取配置文件 \"%s\" 失败, err=", filepath, err)
		return
	}
	fmt.Printf("读取配置文件 %s 成功: %#v \n", filepath, MyConfig)

	err = model.CrDB.Connect()
	if err != nil {
		fmt.Println("连接 DB 失败")
		return
	}
	fmt.Println("连接 DB 成功")

	err = model.CrDB.Init()
	if err != nil {
		fmt.Println("初始化 DB 失败")
		return
	}
	fmt.Println("初始化 DB 成功")
}
