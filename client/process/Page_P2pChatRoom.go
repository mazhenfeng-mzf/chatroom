package process

import (
	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func (this *View) P2PChatRoomPage(clientId string, clientName string, smsList []*message.SmsP2PData) {
	chatroom := this.P2PChatRoomList[clientId]

	layout := VBox{}

	var allSmsString string
	if len(smsList) > 0 {
		allSmsString = this.setupP2pSmsString_whenOpenChatroom(smsList)
		this.P2pSmsListUpdate_List(clientId, smsList)
	} else {
		allSmsString = ""
	}

	children := []Widget{
		Composite{
			Layout: HBox{},
			Children: []Widget{
				TextEdit{
					Text:     allSmsString,
					AssignTo: &chatroom.MessageWindow,
					ReadOnly: true,
					VScroll:  true,
				},
			},
		},
		Composite{
			Layout: HBox{},
			Children: []Widget{
				TextEdit{
					AssignTo: &chatroom.Input,
					OnKeyDown: func(key walk.Key) { //键盘事件
						if key == walk.KeyReturn { //回车键
							ClientProc.SendP2PSms(clientId, chatroom.Input.Text())
						}
					},
				},
				Composite{
					Layout: VBox{},
					Children: []Widget{
						PushButton{
							Text: "发送",
							OnClicked: func() {
								ClientProc.SendP2PSms(clientId, chatroom.Input.Text())
							},
						},
						PushButton{
							Text: "查看历史消息",
							OnClicked: func() {
								go this.OpenP2PChatRoom_SmsHistory(clientId)
							},
						},
					},
				},
			},
		},
	}

	chatroom.mw.Title = fmt.Sprintf("与 %s(%s) 的聊天室", clientName, clientId)
	chatroom.mw.Size = Size{800, 600}
	chatroom.mw.Layout = layout
	chatroom.mw.Children = children
	chatroom.mw.AssignTo = &chatroom.mwAssign

	MyLOG.Log("主程序已进入与 %s(%s) 的P2P聊天室, 让子程序完成填充内容", clientName, clientId)
	chatroom.mw.Run()

	MyLOG.Log("退出与 %s(%s) 的P2P聊天室", clientName, clientId)
	MyView.CloseP2PChatRoom(clientId)
	//ClientProc.ExitPC()
}

func (this *View) setupP2pSmsString_whenOpenChatroom(smsList []*message.SmsP2PData) (allSms string) {
	allSms = ""
	one_string := ""
	var name string
	for _, smsData := range smsList {

		if ClientProc.Id == smsData.FromClientId {
			name = "我"
		} else {
			name = smsData.FromClientName
		}

		one_string = fmt.Sprintf("%s - %s 说: %s\r\n", smsData.Time, name, smsData.Data)
		allSms += one_string
	}
	return
}
