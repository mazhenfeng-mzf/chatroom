package process

import (
	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func (this *View) SmMsgBoxPage() {

	layout := VBox{}
	children := []Widget{}

	for index, smMsg := range MySmMsgBox.SmMap {

		MyLOG.Log("展示系统消息: (index:%d) %#v", index, *smMsg)

		pushBtn := PushButton{
			Text:     fmt.Sprintf("%s(%s)", smMsg.Title, "未读"),
			AssignTo: &smMsg.AssignPushBtn,
			OnClicked: func() {
				MyLOG.Log("点击查看系统消息: index(%d): %#v", index, *smMsg)
				MyView.OpenSmMsgPage(index)
			},
		}

		children = append(children, pushBtn)
	}

	this.mw_SmMsgBoxPage.Title = fmt.Sprintf("系统消息 - %s", ClientProc.Name)
	this.mw_SmMsgBoxPage.Layout = layout
	this.mw_SmMsgBoxPage.Children = children
	this.mw_SmMsgBoxPage.Run()
}

func (this *View) OpenSmMsgPage(SmIndex uint16) {

	smMsg := MySmMsgBox.SmMap[SmIndex]
	//MyLOG.Log("打开系统消息详情: %#v", *smMsg)

	//MySmMsgBox.Delete(smMsg.Index)
	MySmMsgBox.Display_MainPage()

	smMsg.AssignPushBtn.SetText(fmt.Sprintf("%s(%s)", smMsg.Title, "已读"))

	//MySmMsgBox.Delete(smMsg.Index)

	// MySmMsgBox.Refresh_SmMsgBoxPage()
	//MyLOG.Log("smMsg.Data: %#v", smMsg.Data)

	switch smMsg.Type {
	case message.FriendsAddRequest:
		ClientProc.send_SmMessageReadNotify_FriendsAddRequest(smMsg.Data)
		this.FriendAddRequestPage(smMsg.Data.(message.FriendsAddRequestData))
	case message.FriendsAddAccept:
		ClientProc.send_SmMessageReadNotify_FriendsAddAccept(smMsg.Data)
		this.FriendAddAcceptPage(smMsg.Data.(message.FriendsAddAcceptData))
	case message.FriendsAddReject:
		ClientProc.send_SmMessageReadNotify_FriendsAddReject(smMsg.Data)
		this.FriendAddRejectPage(smMsg.Data.(message.FriendsAddRejectData))
	}
}

func (this *View) FriendAddRequestPage(far_d message.FriendsAddRequestData) (causeId message.CauseId) {

	var mw_assign *walk.MainWindow

	fromClientId := far_d.FromClientId
	fromClientName := far_d.FromClientName

	MainWindow{
		AssignTo: &mw_assign,
		Title:    "系统消息 - 添加好友",
		Size:     Size{200, 50},
		Layout:   VBox{},
		Children: []Widget{
			TextLabel{
				Text: fmt.Sprintf("用户 %s(%s) 请求添加你为好友", fromClientName, fromClientId),
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "同意",
						OnClicked: func() {
							ClientProc.sendFriendAddAccept(fromClientId, fromClientName)
							mw_assign.Close()
						},
					},
					PushButton{
						Text: "拒绝",
						OnClicked: func() {
							ClientProc.sendFriendAddReject(fromClientId, fromClientName)
							mw_assign.Close()
						},
					},
					PushButton{
						Text: "忽略",
						OnClicked: func() {
							mw_assign.Close()
						},
					},
				},
			},
		},
	}.Run()
	return
}

func (this *View) FriendAddAcceptPage(acc_d message.FriendsAddAcceptData) (causeId message.CauseId) {

	var mw_assign *walk.MainWindow

	//fromClientId := acc_d.FromClientId
	//fromClientName := acc_d.FromClientName
	toClientId := acc_d.ToClientId
	toClientName := acc_d.ToClientName

	MainWindow{
		AssignTo: &mw_assign,
		Title:    "系统消息 - 添加好友成功",
		Size:     Size{200, 50},
		Layout:   VBox{},
		Children: []Widget{
			TextLabel{
				Text: fmt.Sprintf("您和用户 %s(%s) 已经成为好友", toClientName, toClientId),
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "确认",
						OnClicked: func() {
							mw_assign.Close()
						},
					},
				},
			},
		},
	}.Run()
	return
}

func (this *View) FriendAddRejectPage(acc_d message.FriendsAddRejectData) (causeId message.CauseId) {

	var mw_assign *walk.MainWindow

	//fromClientId := acc_d.FromClientId
	//fromClientName := acc_d.FromClientName
	toClientId := acc_d.ToClientId
	toClientName := acc_d.ToClientName

	MainWindow{
		AssignTo: &mw_assign,
		Title:    "系统消息 - 添加好友失败",
		Size:     Size{200, 50},
		Layout:   VBox{},
		Children: []Widget{
			TextLabel{
				Text: fmt.Sprintf("您发给 %s(%s) 的好友添加请求被拒绝", toClientName, toClientId),
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "确认",
						OnClicked: func() {
							mw_assign.Close()
						},
					},
				},
			},
		},
	}.Run()
	return
}
