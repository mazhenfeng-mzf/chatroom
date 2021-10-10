package process

import (
	. "chatroom/public/tools"
	"fmt"

	. "github.com/lxn/walk/declarative"
)

func (this *View) MainPage() {

	layout := VBox{}
	children := []Widget{
		// HSplitter{
		// 	Children: []Widget{
		// 		TextEdit{AssignTo: &inTE},
		// 		TextEdit{AssignTo: &outTE, ReadOnly: true},
		// 	},
		// },
		TextLabel{
			Text: fmt.Sprintf("欢迎使用聊天室 - %s (%s)", ClientProc.Name, ClientProc.Id),
			// MaxSize: Size{
			// 	Width:  20,
			// 	Height: 20,
			// },
		},
		PushButton{
			Text: "查看我的信息",
			OnClicked: func() {

			},
		},
		PushButton{
			Text:     fmt.Sprintf("收件箱(%d)", MySmsBox.SmsSum),
			AssignTo: &MyView.MainPage_SmsBox,
			OnClicked: func() {
				MyView.OpenSmsBox()
			},
		},
		PushButton{
			Text: "进入公共聊天室",
			OnClicked: func() {
				ClientProc.EnterPC()
			},
		},
		PushButton{
			Text: "查看我的好友",
			OnClicked: func() {
				this.FriendsPage()
			},
		},
		PushButton{
			Text: "查找用户",
			OnClicked: func() {
				MyView.searchClientPage()
			},
		},
		PushButton{
			Text:     fmt.Sprintf("系统消息(%d)", MySmMsgBox.SmMsgSum),
			AssignTo: &MyView.MainPage_SmMsgBox,
			OnClicked: func() {
				MyView.OpenSmMsgBox()
			},
		},
		PushButton{
			Text: "退出登录",
			OnClicked: func() {
				ClientProc.Logout()
				this.mwAssign_MainPage.Close()
				MyLOG.End()
				ClientProc.MainInit()
				this.LoginView()
			},
		},
	}

	this.mw_MainPage.Title = fmt.Sprintf("聊天室 - %s", ClientProc.Name)
	this.mw_MainPage.Layout = layout
	this.mw_MainPage.Children = children
	this.mw_MainPage.Run()

	ClientProc.Logout()

}
