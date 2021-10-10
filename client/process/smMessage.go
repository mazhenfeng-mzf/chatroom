package process

import (
	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"

	"github.com/lxn/walk"
)

// 存放 SM 消息的结构体
type SmMsgBox struct {
	SmMap    map[uint16]*SmMsg
	SmMsgSum uint16
	Index    uint16
}

var MySmMsgBox SmMsgBox

type SmMsg struct {
	Index         uint16
	Title         string
	Type          uint8
	Data          interface{}
	HasRead       bool //还没用
	AssignPushBtn *walk.PushButton
}

func (this *SmMsgBox) Init() {

	this.SmMap = make(map[uint16]*SmMsg, 4096)
	this.SmMsgSum = 0
	this.Index = 0
}

func (this *SmMsgBox) Add(msgType uint8, sm_data interface{}) (err error) {

	//sm_data: sush as message.FriendsAddRequest, message.FriendsAddAccept, message.FriendsAddReject,

	this.Index += 1
	smMsg := SmMsg{
		Index:   this.Index,
		Type:    msgType,
		Data:    sm_data,
		HasRead: false,
	}

	switch msgType {
	case message.FriendsAddRequest:
		smMsg.Title = "好友添加请求"
	case message.FriendsAddAccept:
		smMsg.Title = "好友添加请求 - 接受"
	case message.FriendsAddReject:
		smMsg.Title = "好友添加请求 - 拒绝"
	}

	this.SmMsgSum += 1
	this.SmMap[this.Index] = &smMsg
	MyLOG.Log("after finish Add, this.SmMap: %#v", this.SmMap)
	return
}

// func (this *SmMsgBox) Add_Slice(msgType uint8, slice_sm_data []interface{}) (err error) {
//  //[]interface, 在调用的时候报错
// 	for _, sm_data := range slice_sm_data {
// 		this.Add(msgType, sm_data)
// 	}
// 	return
// }

func (this *SmMsgBox) Delete(index uint16) (err error) {

	if _, ok := this.SmMap[index]; !ok {
		return
	}

	delete(this.SmMap, index)
	this.SmMsgSum -= 1
	return
}

func (this *SmMsgBox) Display_MainPage() {
	MyLOG.Log("刷新主界面 系统消息数量 this.SmMsgSum: %d", this.SmMsgSum)
	MyView.MainPage_SmMsgBox.SetText(fmt.Sprintf("系统消息(%d)", this.SmMsgSum))
}

// func (this *SmMsgBox) Refresh_SmMsgBoxPage() {
// 	MyLOG.Log("刷新 系统消息列表")
// 	MyView.mwAssign_SmMsgBoxPage.Close()
// 	MyView.SmMsgBoxPage()
// }
