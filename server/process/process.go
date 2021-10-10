package process

import (
	"chatroom/public/errs"
	"chatroom/public/message"
	"chatroom/public/processpkg"
	. "chatroom/public/tools"
	"chatroom/server/model"
	"encoding/json"
	"errors"
	"net"
)

type ClientProcess struct {
	Conn          net.Conn
	Model_Id      uint16
	Client_id     string
	Client_name   string
	Client_online int
	Msg           *message.MessageDataR
}

func Connect_DB() (err error) {
	err = model.CrDB.Connect()
	if err != nil {
		MyLOG.Log("连接 DB fail")
		return
	}
	return
}

func FindProcess(conn net.Conn) (err error) {
	//不同消息类型调用不用的 process
	cliProcess := ClientProcess{
		Conn: conn,
	}

	// defer func() {
	// 	err := recover()
	// 	if err != nil {
	// 		MyLOG.ErrLog("服务器错误, 客户端 %#v 只能下线: %#v", cliProcess.Conn, err)
	// 		cliProcess.Conn.Close()
	// 		return
	// 	}
	// }()

	for {
		MyLOG.Log("等待客户端.......")
		var readMsg message.MessageDataR
		err = processpkg.ReadPkg(conn, &readMsg)
		if err != nil {
			MyLOG.ErrLog("客户端失联: %v", cliProcess.Conn)
			cliProcess.ClientOffline()
			return
		}

		MyLOG.Log("收到客户端信息: %#v", readMsg)
		cliProcess.Msg = &readMsg

		switch readMsg.Type {
		case message.LoginRequest:
			cliProcess.loginProc()
			MyLOG.Log("客户端登录: %v; 账号: %s (%s)", cliProcess.Conn, cliProcess.Client_name, cliProcess.Client_id)
		case message.Logout:
			cliProcess.logoutProc()
			MyLOG.Log("客户端下线: %v; 账号: %s (%s)", cliProcess.Conn, cliProcess.Client_name, cliProcess.Client_id)
			return
		case message.RegisterRequest:
			cliProcess.registerProc()
		case message.EnterPCRequest:
			cliProcess.enterPcProc()
		case message.ExitPC:
			cliProcess.exitPcProc()
		case message.DeRegister:
			MyLOG.Log("用户去注册")
			return
		case message.SearchClientResquest:
			cliProcess.SearchClientProc()

		case message.SmsPublicChatroom:
			cliProcess.ReceiveSmsPCMsg()
		case message.SmsP2P:
			cliProcess.ReceiveSmsP2PMsg()

		//Friends
		case message.FriendsAddRequest:
			cliProcess.receiveFriendsAddRequest()
		case message.FriendsAddAccept:
			cliProcess.receiveFriendsAddAccept()
		case message.FriendsAddReject:
			cliProcess.receiveFriendsAddReject()

		//
		case message.SmMessageReadNotify_FriendsAddRequest:
			cliProcess.receive_SmMessageReadNotify_FriendsAddRequest()
		case message.SmMessageReadNotify_FriendsAddAccept:
			cliProcess.receive_SmMessageReadNotify_FriendsAddAccept()
		case message.SmMessageReadNotify_FriendsAddReject:
			cliProcess.receive_SmMessageReadNotify_FriendsAddReject()

		case message.MsgInboxCheckNotify:
			cliProcess.receive_MsgInboxCheckNotify()

		case message.SmsHistoryRequest:
			cliProcess.receive_SmsHistoryRequest()

		default:
			err = errors.New("无效的消息类型")
			MyLOG.Log("Error: %v", err)

		}
	}
}

func (this *ClientProcess) ClientOffline() (causeId message.CauseId) {
	CliMgr.OnlineCliDelete(this.Client_id)
	this.Conn.Close()

	//tell the friends(online): I'am offline
	this.friendsBroadcastOnOffLine(false)

	return
}

func (this *ClientProcess) ClientOnline() (causeId message.CauseId) {
	CliMgr.OnlineCliAdd(this)

	this.friendsBroadcastOnOffLine(true)
	return
}

func (this *ClientProcess) friendsBroadcastOnOffLine(online bool) (causeId message.CauseId) {
	my_cud := message.ClientUserData{
		Id:     this.Client_id,
		Name:   this.Client_name,
		Online: online,
	}

	Msg := message.MessageDataW{
		Type: message.FriendOnOffLine,
		Data: my_cud,
	}

	slice_friends, _ := model.CrDB.FriendsSelect(this.Model_Id)
	for _, ucd := range slice_friends {
		if friendsProc, online := CliMgr.ClientIsOnline(ucd.Id); online {
			processpkg.WritePkg(friendsProc.Conn, &Msg)
			MyLOG.Log("通知 %s(%s) 的好友 %s(%s): %s online=%v",
				this.Client_name, this.Client_id, friendsProc.Client_name, friendsProc.Client_id, this.Client_name, my_cud.Online)
		}
	}

	return
}

func (this *ClientProcess) registerProc() (causeId message.CauseId) {

	//Data 反序列化
	var register_req_d message.RegisterRequestData
	err := json.Unmarshal(this.Msg.Data, &register_req_d)
	if err != nil {
		MyLOG.ErrLog("注册信息解析失败")
		this.write_register_rej(message.ID_SERVER_ERROR)
		return
	}
	MyLOG.Log("注册ID: %s", register_req_d.Id)
	exist, causeId := model.CrDB.CountExist(register_req_d.Id)
	if causeId != 0x0 {
		MyLOG.ErrLog("注册失败: 数据库错误")
		this.write_register_rej(message.ID_SERVER_ERROR)
		return
	}
	if exist {
		MyLOG.Log("用户存在")
		this.write_register_rej(message.ID_COUNT_EXIST)
		return
	}

	causeId = model.CrDB.CountRegister(register_req_d.Id, register_req_d.Name, register_req_d.Pwd)
	if causeId != 0x0 {
		MyLOG.ErrLog("注册失败: 数据库错误")
		this.write_register_rej(message.ID_SERVER_ERROR)
		return
	}

	this.write_register_accept()

	return

}

func (this *ClientProcess) loginProc() (causeId message.CauseId) {

	//Data 反序列化
	var login_req_d message.LoginRequestData
	err := json.Unmarshal(this.Msg.Data, &login_req_d)
	if err != nil {
		MyLOG.ErrLog("登录信息解析失败")
		this.write_login_rej(message.ID_SERVER_ERROR)
		err = errs.INVAILD_MSG_PARSE_FAIL
		return
	}

	exist, causeId := model.CrDB.CountExist(login_req_d.Id)
	if !exist {
		MyLOG.ErrLog("用户不存在")
		this.write_login_rej(message.ID_COUNT_NOT_EXIST)
		return
	}

	if causeId != 0x0 {
		MyLOG.ErrLog("用户 %s 登录失败: %s", login_req_d.Id, message.CauseMap[causeId])
		this.write_login_rej(causeId)
		return
	}

	var client_info model.Count
	causeId = model.CrDB.CountLogin(login_req_d.Id, login_req_d.Pwd, &client_info)
	if causeId != 0x0 {
		MyLOG.ErrLog("用户 %s 登录失败: %s", login_req_d.Id, message.CauseMap[causeId])
		this.write_login_rej(causeId)
		return
	}
	this.Client_id = login_req_d.Id
	this.Client_name = client_info.CountName
	this.Model_Id = client_info.Id

	this.ClientOnline()
	this.write_login_accept()
	MyLOG.Log("完成 登录流程 loginProc")
	return
}

func (this *ClientProcess) logoutProc() (causeId message.CauseId) {

	//Data 反序列化
	var logout_d message.LogoutData
	err := json.Unmarshal(this.Msg.Data, &logout_d)
	if err != nil {
		MyLOG.ErrLog("下线信息解析失败")
		err = errs.INVAILD_MSG_PARSE_FAIL
		return
	}
	this.ClientOffline()
	return
}

func (this *ClientProcess) write_login_rej(causeId message.CauseId) {
	login_rej_data := message.LoginRejectData{
		Cause: causeId,
	}

	rstMsg := message.MessageDataW{
		Type: message.LoginReject,
		Data: login_rej_data,
	}
	MyLOG.ErrLog("用户登录失败: %v", message.CauseMap[causeId])
	processpkg.WritePkg(this.Conn, &rstMsg)

}

func (this *ClientProcess) write_login_accept() {

	var CliforClient message.ClientUserData
	CliforClient.Id = this.Client_id
	CliforClient.Name = this.Client_name

	// Sms Inbox
	sliceMsgInbox, _ := model.CrDB.ClientInboxSelect(this.Client_id, this.Model_Id)
	//MyLOG.Log("write_login_accept - sliceMsgInbox: %#v", sliceMsgInbox)
	//model.CrDB.ClientInboxDelete(this.Client_id, this.Model_Id)

	// Sm Msg FriendsAddRequest
	slice_FAR, _ := model.CrDB.FriendsAddRequestSelect(this.Model_Id)
	MyLOG.Log("从数据库读取到 Slice_FriendsAddRequest: %#v", slice_FAR)
	// Sm Msg FriendsAddAccept
	slice_FAA, _ := model.CrDB.FriendsAddAcceptSelect(this.Model_Id)
	MyLOG.Log("从数据库读取到 Slice_FriendsAddAccept: %#v", slice_FAA)
	// Sm Msg FriendsAddReject
	slice_FARJ, _ := model.CrDB.FriendsAddRejectSelect(this.Model_Id)
	MyLOG.Log("从数据库读取到 Slice_FriendsAddReject: %#v", slice_FARJ)

	//Friends
	slice_friends, _ := model.CrDB.FriendsSelect(this.Model_Id)
	for _, ucd := range slice_friends {
		_, ucd.Online = CliMgr.ClientIsOnline(ucd.Id)
	}

	login_accept_data := message.LoginAcceptData{
		MyselfClient:            &CliforClient,
		MsgInbox:                sliceMsgInbox,
		Slice_FriendsAddRequest: slice_FAR,
		Slice_FriendsAddAccept:  slice_FAA,
		Slice_FriendsAddReject:  slice_FARJ,
		Friends:                 slice_friends,
	}

	rstMsg := message.MessageDataW{
		Type: message.LoginAccept,
		Data: login_accept_data,
	}

	MyLOG.Log("用户登录成功: %s(%s)", this.Client_name, this.Client_id)
	processpkg.WritePkg(this.Conn, &rstMsg)

}

func (this *ClientProcess) write_register_rej(causeId message.CauseId) {
	register_rej_data := message.RegisterRejectData{
		Cause: causeId,
	}

	rstMsg := message.MessageDataW{
		Type: message.RegisterReject,
		Data: register_rej_data,
	}
	MyLOG.ErrLog("用户注册失败: %s(%d)", message.CauseMap[causeId], causeId)
	processpkg.WritePkg(this.Conn, &rstMsg)

}

func (this *ClientProcess) write_register_accept() {
	register_accept_data := message.RegisterAcceptData{}

	rstMsg := message.MessageDataW{
		Type: message.RegisterAccept,
		Data: register_accept_data,
	}
	MyLOG.Log("用户注册成功")
	processpkg.WritePkg(this.Conn, &rstMsg)

}

func (this *ClientProcess) enterPcProc() (causeId message.CauseId) {

	//Data 反序列化
	// var enterPC_d message.EnterPCData
	// err = json.Unmarshal(this.Msg.Data, &enterPC_d)
	// if err != nil {
	// 	MyLOG.Log("登录信息解析失败")
	// 	this.write_login_rej(message.ID_SERVER_ERROR)
	// 	err = errs.INVAILD_MSG_PARSE_FAIL
	// 	return
	// }

	CliMgr.PublicCrCliAdd(this)
	this.write_EnterPCAccept()
	this.broadcastMeInPC()
	return
}

func (this *ClientProcess) exitPcProc() (causeId message.CauseId) {

	CliMgr.PublicCrCliDelete(this)
	this.broadcastMeOutPC()
	return
}

func (this *ClientProcess) write_EnterPCAccept() {

	data := message.EnterPCAcceptData{
		PcClientList: CliMgr.PublicChatroomClientsForClient,
	}

	Msg := message.MessageDataW{
		Type: message.EnterPCAccept,
		Data: data,
	}

	MyLOG.Log("用户 %s (%s) 进入公共聊天室", this.Client_name, this.Client_id)
	processpkg.WritePkg(this.Conn, &Msg)
}

func (this *ClientProcess) broadcastMeInPC() {
	data := message.UpdatePublicChatClientsData{
		Id:   this.Client_id,
		Name: this.Client_name,
		InPC: true,
	}

	Msg := message.MessageDataW{
		Type: message.UpdatePublicChatClients,
		Data: data,
	}

	for _, cliPro := range CliMgr.PublicChatroomClients {
		if cliPro.Client_id == this.Client_id {
			continue
		}
		processpkg.WritePkg(cliPro.Conn, &Msg)
		MyLOG.Log("广播给用户 %s: 用户 %s (%s) 进入公共聊天室", cliPro.Client_name, this.Client_name, this.Client_id)
	}
	MyLOG.Log("广播完成: 用户 %s (%s) 进入公共聊天室", this.Client_name, this.Client_id)
}

func (this *ClientProcess) broadcastMeOutPC() {
	data := message.UpdatePublicChatClientsData{
		Id:   this.Client_id,
		Name: this.Client_name,
		InPC: false,
	}

	Msg := message.MessageDataW{
		Type: message.UpdatePublicChatClients,
		Data: data,
	}

	for _, cliPro := range CliMgr.PublicChatroomClients {
		if cliPro.Client_id == this.Client_id {
			continue
		}
		processpkg.WritePkg(cliPro.Conn, &Msg)
		MyLOG.Log("广播给用户 %s: 用户 %s (%s) 退出公共聊天室", cliPro.Client_name, this.Client_name, this.Client_id)
	}
	MyLOG.Log("广播完成: 用户 %s (%s) 退出公共聊天室", this.Client_name, this.Client_id)
}

func (this *ClientProcess) SearchClientProc() (causeId message.CauseId) {

	//Data 反序列化
	var req_d message.SearchClientResquestData
	err := json.Unmarshal(this.Msg.Data, &req_d)
	if err != nil {
		MyLOG.Log("SearchClient信息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		this.write_register_rej(causeId)
		return
	}

	model_count, causeId := model.CrDB.CountSearch(req_d.Id)
	if causeId != 0x0 {
		this.write_SearchClientFail(causeId)
		MyLOG.ErrLog("用户 %s(%s) 请求搜索用户 %s 失败: %s", this.Client_name, this.Client_id, req_d.Id, message.CauseMap[causeId])
		return
	}

	clientName := model_count.CountName
	clientId := model_count.CountId
	clientdata := message.ClientUserData{
		Id:   clientId,
		Name: clientName,
	}

	data := message.SearchClientResponseData{
		ClientUserData: clientdata,
	}

	Msg := message.MessageDataW{
		Type: message.SearchClientResponse,
		Data: data,
	}

	MyLOG.Log("用户 %s(%s) 请求搜索用户: %s(%s)", this.Client_name, this.Client_id, clientName, clientId)
	processpkg.WritePkg(this.Conn, &Msg)

	return
}

func (this *ClientProcess) write_SearchClientFail(causeId message.CauseId) {

	data := message.SearchClientFailData{
		CauseId: causeId,
	}

	Msg := message.MessageDataW{
		Type: message.SearchClientFail,
		Data: data,
	}

	processpkg.WritePkg(this.Conn, &Msg)
}

func (this *ClientProcess) receive_SmMessageReadNotify_FriendsAddRequest() (causeId message.CauseId) {

	//Data 反序列化
	var smrn_far_d message.SmMessageReadNotify_FriendsAddRequest_Data
	err := json.Unmarshal(this.Msg.Data, &smrn_far_d)
	if err != nil {
		MyLOG.ErrLog("SmMessageReadNotify_FriendsAddRequest 消息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		return
	}

	model.CrDB.FriendsBoxDelete(smrn_far_d.FromClientId, smrn_far_d.ToClientId, smrn_far_d.Message, message.FriendsAddRequest)

	return
}

func (this *ClientProcess) receive_SmMessageReadNotify_FriendsAddAccept() (causeId message.CauseId) {

	//Data 反序列化
	var smrn_faa_d message.SmMessageReadNotify_FriendsAddAccept_Data
	err := json.Unmarshal(this.Msg.Data, &smrn_faa_d)
	if err != nil {
		MyLOG.ErrLog("SmMessageReadNotify_FriendsAddAccept 消息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		return
	}

	model.CrDB.FriendsBoxDelete(smrn_faa_d.FromClientId, smrn_faa_d.ToClientId, "", message.FriendsAddAccept)

	return
}

func (this *ClientProcess) receive_SmMessageReadNotify_FriendsAddReject() (causeId message.CauseId) {

	//Data 反序列化
	var smrn_farj_d message.SmMessageReadNotify_FriendsAddReject_Data
	err := json.Unmarshal(this.Msg.Data, &smrn_farj_d)
	if err != nil {
		MyLOG.ErrLog("SmMessageReadNotify_FriendsAddReject 消息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		return
	}

	model.CrDB.FriendsBoxDelete(smrn_farj_d.FromClientId, smrn_farj_d.ToClientId, "", message.FriendsAddReject)

	return
}

func (this *ClientProcess) receive_MsgInboxCheckNotify() (causeId message.CauseId) {

	//Data 反序列化
	var data message.MsgInboxCheckNotifyData
	err := json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("MsgInboxCheckNotify 消息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		return
	}
	fromClientId := data.FromClientId
	toClientId := data.ToClientId

	model.CrDB.ClientInboxDelete(fromClientId, toClientId)

	return
}

func (this *ClientProcess) receive_SmsHistoryRequest() (causeId message.CauseId) {

	//Data 反序列化
	var data message.SmsHistoryRequestData
	err := json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("SmsHistoryRequest 消息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		return
	}
	fromClientId := data.FromClientId
	toClientId := data.ToClientId
	offset := data.Offset
	number := data.Number

	smsHistorySum := model.CrDB.SmsHistorySum(fromClientId, toClientId)
	smsHistorySlice := model.CrDB.SmsHistorySelect(fromClientId, toClientId, number, offset)

	this.write_SmsHistoryResponseData(fromClientId, toClientId, smsHistorySum, smsHistorySlice)
	return
}

func (this *ClientProcess) write_SmsHistoryResponseData(fromClientId string, toClientId string, smsSum uint32, smsHistorySlice []*message.SmsP2PData) (causeId message.CauseId) {

	data := message.SmsHistoryResponseData{
		FromClientId:    fromClientId,
		ToClientId:      toClientId,
		SmsHistorySum:   smsSum,
		SmsHistorySlice: smsHistorySlice,
	}

	Msg := message.MessageDataW{
		Type: message.SmsHistoryResponse,
		Data: data,
	}

	processpkg.WritePkg(this.Conn, &Msg)
	return
}
