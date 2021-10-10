package process

import (
	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var searchInfo *walk.TextEdit
var addFriendBtn *walk.PushButton
var findClientId string
var findClientName string

func (this *View) searchClientPage() {

	var searchID *walk.LineEdit

	layout := VBox{}
	children := []Widget{
		Composite{
			Layout: HBox{},
			Children: []Widget{
				Label{
					Text: "输入用户ID",
				},
				LineEdit{
					AssignTo: &searchID,
				},
				PushButton{
					Text: "搜索",
					OnClicked: func() {
						ClientProc.SearchClient(searchID.Text())
					},
				},
			},
		},
		Composite{
			Layout: HBox{},
			Children: []Widget{
				TextEdit{
					AssignTo: &searchInfo,
					ReadOnly: true,
				},
				Composite{
					Layout: VBox{},
					Children: []Widget{
						PushButton{
							AssignTo: &addFriendBtn,
							Text:     "添加好友",
							OnClicked: func() {
								ClientProc.sendFriendAddRequest(findClientId)
							},
						},
						PushButton{
							AssignTo: &addFriendBtn,
							Text:     "发送消息",
							OnClicked: func() {
								enptySmsList := make([]*message.SmsP2PData, 0, 0)
								MyView.OpenP2PChatRoom(findClientId, findClientName, enptySmsList)
							},
						},
					},
				},
			},
		},
	}

	this.mw_searchPage.Title = fmt.Sprintf("聊天室 - 查找用户")
	this.mw_searchPage.Size = Size{300, 200}
	this.mw_searchPage.Layout = layout
	this.mw_searchPage.Children = children
	this.mw_searchPage.AssignTo = &this.mwAssign_searchPage

	MyLOG.Log("主程序已打开搜索用户界面")

	this.mw_searchPage.Run()
}

func (this *View) DisplaySearchClientReponse(clientdata *message.ClientUserData) {

	findClientId = clientdata.Id
	findClientName = clientdata.Name

	all_string := fmt.Sprintf("搜索用户成功:\r\nId: %s\r\nName: %s", findClientId, findClientName)

	searchInfo.SetText(all_string)
}

func (this *View) DisplaySearchClientFail(causeId message.CauseId) {

	cause := message.CauseMap[causeId]

	searchInfo.SetText(cause)
}
