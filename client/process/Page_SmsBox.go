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
		smsNum = uint16(len(smsList))
		MyLOG.Log("展示收件箱: 来自 %s 的信息数量: %d  地址：%#v", fromClientId, smsNum, smsBox_forClient)

		//fromClientId 下面的方法用这个变量的话，只能拿到最后一次遍历的值
		// 要每次遍历声明变量 来给方法用才行
		// 也就是： 每次遍历都有一个静态方法
		fromClientId_working := smsBox_forClient.FromClientId
		fromClientName_working := smsBox_forClient.FromClientName

		if CliMgr.IsMyFriends(fromClientId) {
			fromClientName = CliMgr.MyFriendsClients[fromClientId].Name
			fromClientName_working = CliMgr.MyFriendsClients[fromClientId].Name
		} else {
			fromClientName = "未知用户"
			fromClientName_working = "未知用户"
		}

		pushBtn1 := PushButton{
			Text:     fmt.Sprintf("(%s)(%s)(%d)", fromClientName, fromClientId, smsNum),
			AssignTo: &smsBox_forClient.PushBtnAssign,
			OnClicked: func() {

				MyLOG.Log("代码错误展示：点击查看收件箱信息: %s(%s)", fromClientName, fromClientId)
				MyLOG.Log("代码正确展示：点击查看收件箱信息: %s(%s)", fromClientName_working, fromClientId_working)

				MySmsBox.ClickSmsInbox(fromClientId_working, fromClientName_working)
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
