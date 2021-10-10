package process

import (
	"chatroom/public/errs"
	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type View struct {
	mw_LoginPage       MainWindow
	mwAssign_LoginPage *walk.MainWindow

	mw_RegisterPage       MainWindow
	mwAssign_RegisterPage *walk.MainWindow

	//主界面
	mw_MainPage       MainWindow
	mwAssign_MainPage *walk.MainWindow
	MainPage_SmsBox   *walk.PushButton
	MainPage_SmMsgBox *walk.PushButton

	//收件箱
	mw_SmsBoxPage       MainWindow
	mwAssign_SmsBoxPage *walk.MainWindow
	SmsBoxMap           map[string]*SmsBox_forClient //string is Client Id

	mw_searchPage       MainWindow
	mwAssign_searchPage *walk.MainWindow

	//系统消息
	mw_SmMsgBoxPage       MainWindow
	mwAssign_SmMsgBoxPage *walk.MainWindow

	//好友列表
	mw_FriendsPage       MainWindow
	mwAssign_FriendsPage *walk.MainWindow
	FriendsViewList      map[string]*FriendView //string is Client Id

	LoginCount *walk.LineEdit
	LoginPwd   *walk.LineEdit

	PcCR *PublicChatRoom

	Len_SmsList        int
	Len_SmsHistoryList int

	P2PChatRoomList map[string]*P2PChatRoom //string is Client Id
}

type PublicChatRoom struct {
	mw       MainWindow
	mwAssign *walk.MainWindow

	PcMessageWindow *walk.TextEdit
	PcOnlineClient  *walk.TextEdit
	PcInput         *walk.TextEdit

	PcSmsList []*message.SmsPublicChatroomData
	//Len_SmsList int
}

type P2PChatRoom struct {
	mw       MainWindow
	mwAssign *walk.MainWindow

	PeerId   string
	PeerName string

	MessageWindow *walk.TextEdit
	Input         *walk.TextEdit

	SmsList []*message.SmsP2PData

	mw_SmsHistory                   MainWindow
	mwAssign_SmsHistory             *walk.MainWindow
	MessageWindow_SmsHistory        *walk.TextEdit
	SmsHistory_Lable_PageOffset     *walk.Label
	SmsHistory_CurrentPage          uint16
	SmsHistory_AllPage              uint16
	SmsHistory_SmsList              []*message.SmsP2PData
	SmsHistorySum                   uint32
	receive_SmsHistoryResponse_Flag bool
}

type FriendView struct {
	ClientId          string
	ClientName        string
	SendMsg_PushBtn   *walk.PushButton
	Online            bool
	CheckInfo_PushBtn *walk.PushButton
	Online_Lable      *walk.Label
}

// 存放前端 收件箱内容的 结构体
type SmsBox_forClient struct {
	FromClientId   string
	FromClientName string
	PushBtnAssign  *walk.PushButton
}

var MyView View

func (this *View) Init() {

	this.mw_LoginPage = MainWindow{
		Title:    "聊天室",
		Size:     Size{400, 300},
		AssignTo: &this.mwAssign_LoginPage,
	}

	this.mw_MainPage = MainWindow{
		//Title:    "聊天室",
		Size:     Size{400, 300},
		AssignTo: &this.mwAssign_MainPage,
	}

	this.mw_SmsBoxPage = MainWindow{
		//Title:    "收件箱",
		Size:     Size{400, 300},
		AssignTo: &this.mwAssign_SmsBoxPage,
	}

	this.mw_SmMsgBoxPage = MainWindow{
		//Title:    "系统消息",
		Size:     Size{400, 300},
		AssignTo: &this.mwAssign_SmMsgBoxPage,
	}

	this.mw_FriendsPage = MainWindow{
		//Title:    "好友列表",
		Size:     Size{400, 700},
		AssignTo: &this.mwAssign_FriendsPage,
	}

	this.P2PChatRoomList = make(map[string]*P2PChatRoom, 1024)
	this.Len_SmsList = 100
	this.Len_SmsHistoryList = 30

	this.SmsBoxMap = make(map[string]*SmsBox_forClient, 4096)
	this.FriendsViewList = make(map[string]*FriendView, 65535)
}

func (this *View) OpenPublicChatRoom() {
	pcr := PublicChatRoom{}
	this.PcCR = &pcr

	this.PcCR.PcSmsList = make([]*message.SmsPublicChatroomData, 0, this.Len_SmsList)

	MyView.publicCRPage()
}

func (this *View) OpenP2PChatRoom(clientId string, clientName string, smsList []*message.SmsP2PData) {
	p2pcr := P2PChatRoom{
		PeerId: clientId,
	}
	p2pcr.SmsList = make([]*message.SmsP2PData, 0, this.Len_SmsList)

	this.P2PChatRoomList[clientId] = &p2pcr

	ClientProc.sendMsgInboxCheckNotify(clientId)

	MyView.P2PChatRoomPage(clientId, clientName, smsList)
}

func (this *View) CloseP2PChatRoom(clientId string) {

	delete(this.P2PChatRoomList, clientId)
}

func (this *View) OpenSmsBox() {
	MyView.SmsBoxPage()
}

func (this *View) OpenSmMsgBox() {
	MyView.SmMsgBoxPage()
}

func messageBox(title string, info string) {
	walk.MsgBox(nil, title, info, walk.MsgBoxIconInformation)
}

func (this *View) DisplayPcClients() (err error) {

	all_string := ""
	one_string := ""

	for _, cliData := range CliMgr.PublicChatroomClients {
		one_string = fmt.Sprintf("%s (%s)\r\n", cliData.Name, cliData.Id)
		all_string += one_string
	}
	MyLOG.Log("聊天室在线用户: %s", all_string)
	this.PcCR.PcOnlineClient.SetText(all_string)
	return
}

func (this *View) PcSmsListUpdate(sms *message.SmsPublicChatroomData) {

	//reverse
	if len(this.PcCR.PcSmsList) == 0 {
		this.PcCR.PcSmsList = append(this.PcCR.PcSmsList, sms)
		return
	}

	if len(this.PcCR.PcSmsList) < this.Len_SmsList {
		this.PcCR.PcSmsList = append(this.PcCR.PcSmsList, sms)
		return
	}

	if len(this.PcCR.PcSmsList) == this.Len_SmsList {
		this.PcCR.PcSmsList = append(this.PcCR.PcSmsList[1:], sms)
		return
	}
}

func (this *View) DisplayPcSmsList() {
	all_string := ""
	one_string := ""
	var name string
	for _, smsData := range this.PcCR.PcSmsList {

		if ClientProc.Id == smsData.FromClientId {
			name = "我"
		} else {
			name = smsData.FromClientName
		}

		one_string = fmt.Sprintf("%s - %s 说: %s\r\n", smsData.Time.Format("2006-01-02 15:04:05"), name, smsData.Data)
		all_string += one_string
	}
	this.PcCR.PcMessageWindow.SetText(all_string)
	return
}

func (this *View) P2pSmsListUpdate(clientId string, sms *message.SmsP2PData) {

	//reverse
	if len(this.P2PChatRoomList[clientId].SmsList) == 0 {
		this.P2PChatRoomList[clientId].SmsList = append(this.P2PChatRoomList[clientId].SmsList, sms)
		return
	}

	if len(this.P2PChatRoomList[clientId].SmsList) < this.Len_SmsList {
		this.P2PChatRoomList[clientId].SmsList = append(this.P2PChatRoomList[clientId].SmsList, sms)
		return
	}

	if len(this.P2PChatRoomList[clientId].SmsList) == this.Len_SmsList {
		this.P2PChatRoomList[clientId].SmsList = append(this.P2PChatRoomList[clientId].SmsList[1:], sms)
		return
	}
}

func (this *View) P2pSmsListUpdate_List(clientId string, smsList []*message.SmsP2PData) {
	MyLOG.Log("P2pSmsListUpdate_List - Before, MyView.P2PChatRoomList[%s]: %#v", clientId, smsList)
	for _, smsData := range smsList {
		this.P2pSmsListUpdate(clientId, smsData)
	}
	MyLOG.Log("P2pSmsListUpdate_List - After, MyView.P2PChatRoomList[%s]: %#v", clientId, smsList)
}

func (this *View) DisplayP2pSmsList(clientId string) {
	// all_string := ""
	// one_string := ""
	// var name string
	// for _, smsData := range this.P2PChatRoomList[clientId].SmsList {

	// 	if ClientProc.Id == smsData.FromClientId {
	// 		name = "我"
	// 	} else {
	// 		name = smsData.FromClientName
	// 	}

	// 	one_string = fmt.Sprintf("%s - %s 说: %s\r\n", smsData.Time, name, smsData.Data)
	// 	all_string += one_string
	// }

	all_string := this.GetP2PSmsString(this.P2PChatRoomList[clientId].SmsList)

	this.P2PChatRoomList[clientId].MessageWindow.SetText(all_string)
	return
}

func (this *View) GetP2PSmsString(smsList []*message.SmsP2PData) (all_string string) {
	one_string := ""
	var name string
	for _, smsData := range smsList {

		if ClientProc.Id == smsData.FromClientId {
			name = "我"
		} else {
			name = smsData.FromClientName
		}

		one_string = fmt.Sprintf("%s - %s 说: %s\r\n", smsData.Time, name, smsData.Data)
		all_string += one_string
	}
	return
}

func (this *View) P2pChatroomIsOpen(clientId string) (p2pChatroom *P2PChatRoom, yes bool) {
	p2pChatroom, yes = this.P2PChatRoomList[clientId]
	return
}

func (this *View) SmsBox_forClient_Create(fromClientId string) {

	smsbox := SmsBox_forClient{
		FromClientId: fromClientId,
	}

	this.SmsBoxMap[fromClientId] = &smsbox
}

func (this *View) SmsBox_forClient_Delete(fromClientId string) {
	delete(this.SmsBoxMap, fromClientId)
}

func (this *View) OpenP2PChatRoom_SmsHistory(clientId string) {
	p2pcr := this.P2PChatRoomList[clientId]
	p2pcr.receive_SmsHistoryResponse_Flag = false

	ClientProc.sendSmsHistoryRequestData(clientId, 0, this.Len_SmsHistoryList)
	//主程序发送 SmsHistoryRequest
	// serverProcess 程序 接受 SmsHistoryResponse

	err := this.Wait_receiveSmsHistoryResponse(p2pcr)
	if err != nil {
		MyLOG.Log("发送 SmsHistoryRequest 之后没有收到 SmsHistoryResponse, 打开与 %s(%s) 的聊天历史消息失败", p2pcr.PeerName, clientId)
		return
	}

	p2pcr.SmsHistory_AllPage = uint16(p2pcr.SmsHistorySum / uint32(this.Len_SmsHistoryList))
	remainder := p2pcr.SmsHistorySum % uint32(this.Len_SmsHistoryList)
	if remainder > 0 {
		p2pcr.SmsHistory_AllPage += 1
	}
	p2pcr.SmsHistory_CurrentPage = p2pcr.SmsHistory_AllPage

	this.P2PChatRoomSmsHistoryPage(clientId)
}

func (this *View) Wait_receiveSmsHistoryResponse(p2pcr *P2PChatRoom) (err error) {

	for i := 1; i <= 2000; i++ {
		if p2pcr.receive_SmsHistoryResponse_Flag {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	err = errs.SERVER_TIMEOUT
	return
}

func (this *View) UpdateP2PChatRoom_SmsHistory_LastNextPage(clientId string, lastOrNext string) {
	//lastOrNext: "last" or "next", default is "last"

	p2pcr := this.P2PChatRoomList[clientId]
	p2pcr.receive_SmsHistoryResponse_Flag = false

	var currentPage uint16
	if lastOrNext == "next" {
		if p2pcr.SmsHistory_CurrentPage == p2pcr.SmsHistory_AllPage {
			return
		}
		currentPage = p2pcr.SmsHistory_CurrentPage + 1
	} else {
		if p2pcr.SmsHistory_CurrentPage == 1 {
			return
		}
		currentPage = p2pcr.SmsHistory_CurrentPage - 1
	}

	offset := uint32((p2pcr.SmsHistory_AllPage - currentPage) * uint16(this.Len_SmsHistoryList))

	ClientProc.sendSmsHistoryRequestData(clientId, offset, this.Len_SmsHistoryList)
	//主程序发送 SmsHistoryRequest
	// serverProcess 程序 接受 SmsHistoryResponse

	err := this.Wait_receiveSmsHistoryResponse(p2pcr)
	if err != nil {
		MyLOG.Log("发送 SmsHistoryRequest 之后没有收到 SmsHistoryResponse, 打开与 %s(%s) 的聊天历史消息失败", p2pcr.PeerName, clientId)
		return
	}

	p2pcr.SmsHistory_CurrentPage = currentPage

	p2pcr.SmsHistory_Lable_PageOffset.SetText(fmt.Sprintf("%d/%d", p2pcr.SmsHistory_CurrentPage, p2pcr.SmsHistory_AllPage))
	p2pcr.MessageWindow_SmsHistory.SetText(this.GetP2PSmsString(p2pcr.SmsHistory_SmsList))

}

func (this *View) UpdateFriendOnOffLine(friendId string, friendName string, online bool) {

	online_string := "不在线"
	rgb := walk.RGB(0, 0, 0)
	if online {
		online_string = "在线"
		rgb = walk.RGB(0, 255, 0)
	}

	fv := MyView.FriendsViewList[friendId]

	fv.Online_Lable.SetText(fmt.Sprintf("(%s)", online_string))
	fv.Online_Lable.SetTextColor(rgb)
	fv.CheckInfo_PushBtn.SetText(fmt.Sprintf("%s(%s)", friendName, friendId))
}
