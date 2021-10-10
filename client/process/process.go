package process

import (
	. "chatroom/client/config"
	"chatroom/public/errs"
	"chatroom/public/message"
	"chatroom/public/processpkg"
	. "chatroom/public/tools"
	"encoding/json"
	"errors"
	"net"
	"time"
)

type ClientProcess struct {
	Conn net.Conn
	Id   string
	Name string
	Msg  *message.MessageDataR

	//Enter Public Chatroom
	HasEnterPC bool

	//Message Inbox store temp
	MsgInbox []*message.SmsP2PData
}

var ClientProc ClientProcess

func (this *ClientProcess) MainInit() {
	MyLOG.Init(LogClient)
	MySmsBox.Init()
	MySmMsgBox.Init()
	CliMgr.Init()
	MyView.Init()
}

func (this *ClientProcess) connect_socket() (err error) {
	//连接socket
	conn, err := net.Dial("tcp", MyConfig.SERVER_IP+":"+MyConfig.SERVER_PORT)
	if err != nil {
		MyLOG.ErrLog("client dail err=%v", err)
		return
	}
	this.Conn = conn
	return
}

func (this *ClientProcess) Login(id string, pwd string) (err error) {

	err = this.connect_socket()
	if err != nil {
		MyLOG.ErrLog("client dail err=%v", err)
		return
	}
	//defer this.Conn.Close()

	//构建登录消息
	logindata := message.LoginRequestData{
		Id:  id,
		Pwd: pwd,
	}

	loginMsg := message.MessageDataW{
		Type: message.LoginRequest,
		Data: logindata,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &loginMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", loginMsg)

	//查看 从服务器返回的登录结果
	var msg_data message.MessageDataR
	err = processpkg.ReadPkg(this.Conn, &msg_data)
	if err != nil {
		MyLOG.ErrLog("ReadPkg err=%v", err)
		return
	}

	switch msg_data.Type {
	case message.LoginAccept:
		err = this.HandleLoginAccept(&msg_data)
		if err != nil {
			MyLOG.ErrLog("HandleLoginAccept fail %v", err)
			return
		}
		go serverProcess.Start(this.Conn)

	case message.LoginReject:

		//获取失败原因
		rej_data_byte := msg_data.Data
		var rej_data message.LoginRejectData
		err = json.Unmarshal(rej_data_byte, &rej_data)
		if err != nil {
			MyLOG.ErrLog("Unmarshal these data fail: %v (%v)", rej_data_byte, err)
			return
		}
		causeString := message.CauseMap[rej_data.Cause]
		MyLOG.ErrLog("登录失败: %s", causeString)
		err = errors.New(causeString)
		this.Conn.Close()
	}

	return
}

func (this *ClientProcess) Logout() (err error) {

	logoutdata := message.LogoutData{
		Id: this.Id,
	}

	logoutMsg := message.MessageDataW{
		Type: message.Logout,
		Data: logoutdata,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &logoutMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", logoutMsg)
	this.Conn.Close()
	return
}

func (this *ClientProcess) HandleLoginAccept(msg *message.MessageDataR) (err error) {
	//After login is accept
	//1. Get the info of clients online
	//2. Open a go task to accept data from server
	//3. Enter User-Page
	acc_data_byte := msg.Data
	var acc_data message.LoginAcceptData
	err = json.Unmarshal(acc_data_byte, &acc_data)
	if err != nil {
		MyLOG.ErrLog("Unmarshal these data fail: %v", acc_data_byte)
		return
	}
	//CliMgr.OnlineCliInit(acc_data.ClientOnline)

	this.Id = acc_data.MyselfClient.Id
	this.Name = acc_data.MyselfClient.Name

	//handle the MsgInbox and store in smsBox
	MyLOG.Log("HandleLoginAccept = acc_data.MsgInbox: %#v", acc_data.MsgInbox)
	MySmsBox.Add_Slice(acc_data.MsgInbox)

	//handle the SmMsg and store in MySmMsgBox
	slice_FAR := acc_data.Slice_FriendsAddRequest
	MyLOG.Log("LoginAccept 里面有 系统消息 Slice_FriendsAddRequest : %#v", slice_FAR)
	for _, sm_data := range slice_FAR {
		MySmMsgBox.Add(message.FriendsAddRequest, *sm_data)
	}
	slice_FAA := acc_data.Slice_FriendsAddAccept
	MyLOG.Log("LoginAccept 里面有 系统消息 Slice_FriendsAddAccept : %#v", slice_FAR)
	for _, sm_data := range slice_FAA {
		MySmMsgBox.Add(message.FriendsAddAccept, *sm_data)
	}
	slice_FARJ := acc_data.Slice_FriendsAddReject
	MyLOG.Log("LoginAccept 里面有 系统消息 Slice_FriendsAddReject : %#v", slice_FARJ)
	for _, sm_data := range slice_FARJ {
		MySmMsgBox.Add(message.FriendsAddReject, *sm_data)
	}

	CliMgr.FriendsClientsAdd_Slice(acc_data.Friends)

	return
}

func (this *ClientProcess) EnterPC() (err error) {
	if ClientProc.HasEnterPC {
		MyView.PcCR.mwAssign.Fullscreen()
		return
	}

	//构建登录消息
	enterPCData := message.EnterPCData{}

	enterPCMsg := message.MessageDataW{
		Type: message.EnterPCRequest,
		Data: enterPCData,
	}
	serverProcess.IfenterPCAccept = false
	//发送消息
	err = processpkg.WritePkg(this.Conn, &enterPCMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", enterPCMsg)

	//等待子程序接收 EnterPCAccept 后主程序继续处理, 进入公共聊天室, 2秒超时
	err = this.WaitEnterPCAccept()
	if err != nil {
		MyLOG.ErrLog("WaitEnterPCAccept 失败", err)
		return
	}
	MyLOG.Log("主程序进入公共聊天室")
	ClientProc.HasEnterPC = true //子程序 receiveEnterPCAccept 检测主程序 进行公共聊天室的标记, 检测到 true 后,再等待 0.01s 后进行填充工作
	MyView.OpenPublicChatRoom()

	return
}

func (this *ClientProcess) ExitPC() (err error) {

	//构建登录消息
	exitPCData := message.ExitPCData{}

	exitPCMsg := message.MessageDataW{
		Type: message.ExitPC,
		Data: exitPCData,
	}
	ClientProc.HasEnterPC = false
	MyLOG.Log("ExitPC - HasEnterPC: %v", ClientProc.HasEnterPC)
	//发送消息
	err = processpkg.WritePkg(this.Conn, &exitPCMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", exitPCMsg)
	return
}

func (this *ClientProcess) WaitEnterPCAccept() (err error) {

	for i := 1; i <= 2000; i++ {
		if serverProcess.IfenterPCAccept {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	err = errs.SERVER_TIMEOUT
	return
}

func (this *ClientProcess) Register(id string, name string, pwd string) (err error) {

	err = this.connect_socket()
	if err != nil {
		MyLOG.ErrLog("client dail err=%v", err)
		return
	}
	defer this.Conn.Close()

	//构建注册消息
	registerdata := message.RegisterRequestData{
		Id:   id,
		Name: name,
		Pwd:  pwd,
	}

	registerMsg := message.MessageDataW{
		Type: message.RegisterRequest,
		Data: registerdata,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &registerMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	//MyLOG.Log("发送消息成功: %#v", loginMsg)

	//查看 从服务器返回的登录结果
	var msg_data message.MessageDataR
	err = processpkg.ReadPkg(this.Conn, &msg_data)
	if err != nil {
		MyLOG.ErrLog("ReadPkg err=%v", err)
		return
	}

	switch msg_data.Type {
	case message.RegisterAccept:
		return

	case message.RegisterReject:
		//Register is rejected
		rej_data_byte := msg_data.Data
		var rej_data message.RegisterRejectData
		json.Unmarshal(rej_data_byte, &rej_data)
		causeId := rej_data.Cause
		causeString := message.CauseMap[causeId]
		err = errors.New(causeString)
		MyLOG.ErrLog("注册请求被拒绝: %v", causeString)
	}

	return
}

func (this *ClientProcess) SearchClient(clientId string) (err error) {

	//构建消息
	scdata := message.SearchClientResquestData{
		Id: clientId,
	}

	scMsg := message.MessageDataW{
		Type: message.SearchClientResquest,
		Data: scdata,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &scMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", scMsg)

	return
}

func (this *ClientProcess) sendFriendAddRequest(toClientId string) (err error) {

	//构建消息
	farData := message.FriendsAddRequestData{
		FromClientId:   this.Id,
		FromClientName: this.Name,
		ToClientId:     toClientId,
	}

	farMsg := message.MessageDataW{
		Type: message.FriendsAddRequest,
		Data: farData,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &farMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", farMsg)

	return
}

func (this *ClientProcess) sendFriendAddAccept(fromClientId string, fromClientName string) (err error) {

	//构建消息
	faaData := message.FriendsAddAcceptData{
		FromClientId:   fromClientId,
		FromClientName: fromClientName,
		ToClientId:     this.Id,
	}

	faaMsg := message.MessageDataW{
		Type: message.FriendsAddAccept,
		Data: faaData,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &faaMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", faaMsg)

	return
}

func (this *ClientProcess) sendFriendAddReject(fromClientId string, fromClientName string) (err error) {

	//构建消息
	faaData := message.FriendsAddRejectData{
		FromClientId:   fromClientId,
		FromClientName: fromClientName,
		ToClientId:     this.Id,
	}

	faaMsg := message.MessageDataW{
		Type: message.FriendsAddReject,
		Data: faaData,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &faaMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", faaMsg)

	return
}

// func (this *ClientProcess) sendSmMessageReadNotify(smType uint8, msg interface{}) (err error) {

// 	//构建消息
// 	SMRN_Msg := message.MessageDataW{
// 		Type: message.SmMessageReadNotify_FriendsAddRequest,
// 		Data: msg,
// 	}

// 	//发送消息
// 	err = processpkg.WritePkg(this.Conn, &SMRN_Msg)
// 	if err != nil {
// 		MyLOG.ErrLog("WritePkg err=%v", err)
// 		return
// 	}
// 	MyLOG.Log("发送消息成功: %#v", SMRN_Msg)

// 	return
// }

func (this *ClientProcess) send_SmMessageReadNotify_FriendsAddRequest(msg interface{}) (err error) {

	//构建消息
	SMRN_Msg := message.MessageDataW{
		Type: message.SmMessageReadNotify_FriendsAddRequest,
		Data: msg,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &SMRN_Msg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", SMRN_Msg)

	return
}

func (this *ClientProcess) send_SmMessageReadNotify_FriendsAddAccept(msg interface{}) (err error) {

	//构建消息
	SMRN_Msg := message.MessageDataW{
		Type: message.SmMessageReadNotify_FriendsAddAccept,
		Data: msg,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &SMRN_Msg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", SMRN_Msg)

	return
}

func (this *ClientProcess) send_SmMessageReadNotify_FriendsAddReject(msg interface{}) (err error) {

	//构建消息
	SMRN_Msg := message.MessageDataW{
		Type: message.SmMessageReadNotify_FriendsAddReject,
		Data: msg,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &SMRN_Msg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", SMRN_Msg)

	return
}

func (this *ClientProcess) sendMsgInboxCheckNotify(fromClientId string) (err error) {

	//构建消息
	data := message.MsgInboxCheckNotifyData{
		FromClientId: fromClientId,
		ToClientId:   this.Id,
	}

	Msg := message.MessageDataW{
		Type: message.MsgInboxCheckNotify,
		Data: data,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &Msg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", Msg)

	return
}

func (this *ClientProcess) sendSmsHistoryRequestData(toClientId string, offset uint32, smsNumber int) (err error) {

	//构建消息
	data := message.SmsHistoryRequestData{
		FromClientId: this.Id,
		ToClientId:   toClientId,
		Offset:       offset,
		Number:       smsNumber,
	}

	Msg := message.MessageDataW{
		Type: message.SmsHistoryRequest,
		Data: data,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &Msg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", Msg)

	return
}
