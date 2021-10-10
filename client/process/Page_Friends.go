package process

import (
	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func (this *View) FriendsPage() {

	this.OpenFriendViewFlag = true

	layout := VBox{}
	children := []Widget{}

	for _, ucd := range CliMgr.MyFriendsClients {

		online_string := "不在线"
		if ucd.Online {
			online_string = "在线"
		}

		clientId := ucd.Id
		clientName := ucd.Name

		friendView := FriendView{
			ClientId:   clientId,
			ClientName: clientName,
			Online:     ucd.Online,
		}

		rgb := walk.RGB(0, 0, 0)
		if ucd.Online {
			MyLOG.Log("%s在线, 绿色", clientName)
			rgb = walk.RGB(0, 255, 0)
		}

		MyView.FriendsViewList[clientId] = &friendView

		com := Composite{
			Layout: HBox{},
			Children: []Widget{
				PushButton{
					Text:     fmt.Sprintf("%s(%s)", ucd.Name, ucd.Id),
					AssignTo: &friendView.CheckInfo_PushBtn,
					OnClicked: func() {

					},
				},
				Label{
					Text:      fmt.Sprintf("(%s)", online_string),
					TextColor: rgb,
					AssignTo:  &friendView.Online_Lable,
				},
				PushButton{
					Text:     "发送消息",
					AssignTo: &friendView.CheckInfo_PushBtn,
					OnClicked: func() {
						enptySmsList := make([]*message.SmsP2PData, 0, 0)
						MyView.OpenP2PChatRoom(clientId, clientName, enptySmsList)
					},
				},
			},
		}

		children = append(children, com)
	}

	this.mw_FriendsPage.Title = fmt.Sprintf("好友列表 - %s", ClientProc.Name)
	this.mw_FriendsPage.Layout = layout
	this.mw_FriendsPage.Children = children
	this.mw_FriendsPage.Run()

	this.OpenFriendViewFlag = false
}
