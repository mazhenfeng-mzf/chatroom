package process

import (
	"chatroom/public/errs"
	"chatroom/public/message"
	"chatroom/public/processpkg"
	. "chatroom/public/tools"
	"encoding/json"
	"errors"
	"net"
	"time"
)

type ServerProcess struct {
	Conn net.Conn
	Id   string
	Name string
	Msg  *message.MessageDataR

	IfenterPCAccept bool
}

var serverProcess ServerProcess

func (this *ServerProcess) Start(conn net.Conn) {

	this.Conn = conn
	this.Id = ClientProc.Id
	this.Name = ClientProc.Name

	// defer func() {
	// 	err := recover()
	// 	if err != nil {
	// 		this.Conn.Close()
	// 	}
	// }()

	for {
		var readMsg message.MessageDataR
		err := processpkg.ReadPkg(conn, &readMsg)
		if err != nil {
			MyLOG.ErrLog("Error serverProcessMes ReadPkg: %v", err)
			return

		}
		MyLOG.Log("收到服务器信息: %#v", readMsg)
		this.Msg = &readMsg

		switch readMsg.Type {
		case message.EnterPCAccept:
			this.receiveEnterPCAccept()
		case message.UpdatePublicChatClients:
			this.receiveUpdatePublicChatClients()
		case message.SmsPublicChatroom:
			this.receiveSmsPublicChatroom()
		case message.SmsP2P:
			this.receiveSmsP2P()
		case message.SearchClientResponse:
			this.receiveSearchClientResponse()
		case message.SearchClientFail:
			this.receiveSearchClientFail()

		case message.FriendOnOffLine:
			this.receiveFriendOnOffLine()

		//smMessage
		case message.FriendsAddRequest:
			this.receive_FriendsAddRequest()
		case message.FriendsAddAccept:
			this.receive_FriendsAddAccept()
		case message.FriendsAddReject:
			this.receive_FriendsAddReject()

		//smsHistory
		case message.SmsHistoryResponse:
			this.receive_SmsHistoryResponse()

		default:
			err = errors.New("收到无效的消息类型")
			MyLOG.ErrLog("Error: %v", err)

		}
	}
}

func (this *ServerProcess) receiveEnterPCAccept() (err error) {
	//Data 反序列化

	var data message.EnterPCAcceptData
	err = json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("信息解析失败")
		err = errs.INVAILD_MSG_PARSE_FAIL
		return
	}
	MyLOG.Log("收到 receiveEnterPCAccept -- PcClientList: %#v", data.PcClientList)

	this.IfenterPCAccept = true //主程序执行 EnterPC 后, 检查此标志, 为 true 后进入公共聊天室
	CliMgr.PublicChatroomClientsInit(data.PcClientList)
	MyLOG.Log("主程序即将进入公共聊天室")

	//等待主程序进入聊天室后完成填充工作
	err = this.WaitEnterPC()
	if err != nil {
		MyLOG.ErrLog("WaitEnterPC 失败", err)
		return
	}
	MyLOG.Log("子程序进行 公共聊天室填充工作")
	MyView.DisplayPcClients()
	MyLOG.Log("子程序完成 公共聊天室填充工作")
	return
}

func (this *ServerProcess) WaitEnterPC() (err error) {

	for i := 1; i <= 2000; i++ {
		if ClientProc.HasEnterPC {
			time.Sleep(10 * time.Millisecond) // give more time for walk.mw.Run()
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	err = errs.SERVER_TIMEOUT
	return
}

func (this *ServerProcess) receiveUpdatePublicChatClients() (err error) {
	//Data 反序列化

	var data message.UpdatePublicChatClientsData
	err = json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("信息解析失败")
		err = errs.INVAILD_MSG_PARSE_FAIL
		return
	}
	MyLOG.Log("收到 UpdatePublicChatClientsData : %#v", data)
	CliMgr.PublicChatroomClientsUpdate(data)
	MyView.DisplayPcClients()
	return
}

func (this *ServerProcess) receiveSearchClientResponse() (err error) {
	var scr_data message.SearchClientResponseData
	err = json.Unmarshal(this.Msg.Data, &scr_data)
	if err != nil {
		MyLOG.ErrLog("Unmarshal these data fail: %v (%v)", this.Msg.Data, err)
		return
	}
	MyLOG.Log("收到 SearchClientResponse: %#v", scr_data)
	clientData := scr_data.ClientUserData
	MyView.DisplaySearchClientReponse(&clientData)
	return
}

func (this *ServerProcess) receiveSearchClientFail() (err error) {
	var scf_data message.SearchClientFailData
	err = json.Unmarshal(this.Msg.Data, &scf_data)
	if err != nil {
		MyLOG.ErrLog("Unmarshal these data fail: %v (%v)", this.Msg.Data, err)
		return
	}
	MyView.DisplaySearchClientFail(scf_data.CauseId)
	return
}

func (this *ServerProcess) receiveFriendOnOffLine() (err error) {
	var data message.FriendOnOffLineData
	err = json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("Unmarshal these data fail: %v (%v)", this.Msg.Data, err)
		return
	}
	friendId := data.Id
	friendName := data.Name
	online := data.Online

	CliMgr.FriendsClientsAdd(friendId, friendName, online)

	MyView.UpdateFriendOnOffLine(friendId, friendName, online)

	return
}

func (this *ServerProcess) receive_SmsHistoryResponse() (err error) {
	var data message.SmsHistoryResponseData
	err = json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("Unmarshal these data fail: %v (%v)", this.Msg.Data, err)
		return
	}

	friendClientId := data.ToClientId

	p2pcr := MyView.P2PChatRoomList[friendClientId]
	p2pcr.SmsHistory_SmsList = data.SmsHistorySlice
	p2pcr.SmsHistorySum = data.SmsHistorySum
	p2pcr.receive_SmsHistoryResponse_Flag = true
	return
}
