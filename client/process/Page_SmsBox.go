package process

import (
	. "chatroom/public/tools"
	"fmt"

	. "github.com/lxn/walk/declarative"
)

func (this *View) SmsBoxPage() {

	layout := VBox{}
	children := []Widget{}

	var smsNum uint16

	var fromClientName string

	for fromClientId, smsBox_forClient := range MyView.SmsBoxMap {

		smsList := MySmsBox.SmsMap[fromClientId]

		MyLOG.Log("展示收件箱: 来自 %s 的信息: %v", fromClientId, smsList)
		smsNum = uint16(len(smsList))

		if CliMgr.IsMyFriends(fromClientId) {
			fromClientName = CliMgr.MyFriendsClients[fromClientId].Name
		} else {
			fromClientName = "未知用户"
		}

		pushBtn1 := PushButton{
			Text:     fmt.Sprintf("(%s)(%s)(%d)", fromClientName, fromClientId, smsNum),
			AssignTo: &smsBox_forClient.PushBtnAssign,
			OnClicked: func() {
				MyLOG.Log("点击查看收件箱信息: %s(%s)", fromClientName, fromClientId)
				MySmsBox.ClickSmsInbox(fromClientId, fromClientName)
				// smsList = MySmsBox.Get(fromClientId)
				// MyView.OpenP2PChatRoom(fromClientId, "", smsList)
			},
		}

		// children = append(children, pushBtn)
		children = append(children, pushBtn1)
	}

	this.mw_SmsBoxPage.Title = fmt.Sprintf("收件箱 - %s", ClientProc.Name)
	this.mw_SmsBoxPage.Layout = layout
	this.mw_SmsBoxPage.Children = children
	this.mw_SmsBoxPage.Run()

}
