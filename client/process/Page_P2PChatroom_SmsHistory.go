package process

import (
	//. "chatroom/public/tools"
	"fmt"

	. "github.com/lxn/walk/declarative"
)

func (this *View) P2PChatRoomSmsHistoryPage(clientId string) {
	chatroom := this.P2PChatRoomList[clientId]

	clientName := chatroom.PeerName
	sms_string := this.GetP2PSmsString(chatroom.SmsHistory_SmsList)

	layout := VBox{}

	children := []Widget{
		Composite{
			Layout: HBox{},
			Children: []Widget{
				TextEdit{
					Text:     sms_string,
					AssignTo: &chatroom.MessageWindow_SmsHistory,
					ReadOnly: true,
					VScroll:  true,
				},
			},
		},
		Composite{
			Layout: HBox{},
			Children: []Widget{
				PushButton{
					Text: "上一页",
					OnClicked: func() {
						this.UpdateP2PChatRoom_SmsHistory_LastNextPage(clientId, "last")
					},
				},
				Label{
					Text:     fmt.Sprintf("%d/%d", chatroom.SmsHistory_CurrentPage, chatroom.SmsHistory_AllPage),
					AssignTo: &chatroom.SmsHistory_Lable_PageOffset,
				},
				PushButton{
					Text: "下一页",
					OnClicked: func() {
						this.UpdateP2PChatRoom_SmsHistory_LastNextPage(clientId, "next")
					},
				},
			},
		},
	}

	chatroom.mw_SmsHistory.Title = fmt.Sprintf("与 %s(%s) 的历史聊天记录", clientName, clientId)
	chatroom.mw_SmsHistory.Size = Size{800, 600}
	chatroom.mw_SmsHistory.Layout = layout
	chatroom.mw_SmsHistory.Children = children
	chatroom.mw_SmsHistory.AssignTo = &chatroom.mwAssign_SmsHistory

	chatroom.mw_SmsHistory.Run()

}
