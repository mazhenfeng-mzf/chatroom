// package serverProcess

// import (
// 	"chatroom/Public/message"
// 	"chatroom/Public/processpkg"
// 	"chatroom/Public/tools"
// 	"chatroom/Server/model"
// 	"encoding/json"
// 	"fmt"
// 	"net"
// 	"time"
// )

// type ClientProcess struct {
// 	Conn net.Conn
// 	// Client_id   string
// 	// Client_name string
// 	// State       int
// 	ClientData *message.ClientUserData
// 	Msg        *message.MessageDataR
// }

// func (this *ClientProcess) lostConn() {
// 	tools.LOG("Close the connect to client: %s, %s, %d", this.ClientData.Id, this.ClientData.Name, this.ClientData.State)
// 	if this.ClientData.State == message.InPC {
// 		this.receiveExitPublicChatroomMsg()
// 	} else if this.ClientData.State == message.Online {
// 		this.logoutProc(false)
// 	}
// }

// func (this *ClientProcess) registerProc() (err error) {

// 	//Data 反序列化
// 	var register_req_d message.RegisterRequestData
// 	err = json.Unmarshal(this.Msg.Data, &register_req_d)
// 	if err != nil {
// 		tools.LOG("注册信息解析失败")
// 		this.write_register_rej(message.SERVER_ERROR)
// 		return
// 	}

// 	err = model.CrDB.CountExist(register_req_d.Id)
// 	if err == nil {
// 		tools.LOG("用户存在")
// 		this.write_register_rej(message.COUNT_EXIST)
// 		return
// 	}

// 	err = model.CrDB.CountRegister(register_req_d.Id, register_req_d.Name, register_req_d.Pwd)
// 	if err != nil {
// 		tools.LOG("注册失败: 数据库错误")
// 		this.write_register_rej(message.SERVER_ERROR)
// 		return
// 	}

// 	this.write_register_accept()

// 	return

// }

// func (this *ClientProcess) loginProc() (err error) {

// 	//Data 反序列化
// 	var login_req_d message.LoginRequestData
// 	err = json.Unmarshal(this.Msg.Data, &login_req_d)
// 	if err != nil {
// 		tools.LOG("登录信息解析失败")
// 		this.write_login_rej(message.SERVER_ERROR)
// 		return
// 	}

// 	err = model.CrDB.CountExist(login_req_d.Id)
// 	if err != nil {
// 		tools.LOG("用户不存在")
// 		this.write_login_rej(message.COUNT_NOT_EXIST)
// 		return
// 	}
// 	var client_info model.Count
// 	err = model.CrDB.CountLogin(login_req_d.Id, login_req_d.Pwd, &client_info)
// 	if err != nil {
// 		tools.LOG("用户密码不正确")
// 		this.write_login_rej(message.COUNT_WRONG_PASSWORD)
// 		return
// 	}
// 	var cd message.ClientUserData
// 	this.ClientData = &cd
// 	this.ClientData.Id = login_req_d.Id
// 	this.ClientData.Name = client_info.CountName
// 	this.ClientData.State = message.Online

// 	//Client login success, there are 3 things server would do:
// 	//1. add this client to server OnlineClient
// 	//2. Send LoginAccept (include inbox)
// 	//3. Send UpdateClientState to other online clients
// 	CliMgr.OnlineCliAdd(this)

// 	CliMgr.UpdateClientStateOnline(this.ClientData)

// 	inboxSmsDataSlice := this.getInboxSmsSlice(this.ClientData.Id)
// 	this.write_login_accept(inboxSmsDataSlice)

// 	//can not delete now, must wait client check Inbox (send CheckInboxNotify message)
// 	//model.CrDB.ClientInboxDelete(this.ClientData.Id)

// 	fmt.Printf("Client login: %s (%s) Conn: %v", this.ClientData.Name, this.ClientData.Id, this.Conn)
// 	return
// }

// func (this *ClientProcess) getInboxSmsSlice(toClientId string) (smsDataSlice []*message.SmsData) {
// 	//cliInboxSlice := make([]*model.ClientInbox, 0, 4096)
// 	cliInboxSlice := model.CrDB.ClientInboxSelect(toClientId)
// 	//var smsDataSlice []*message.SmsData
// 	smsDataSlice = make([]*message.SmsData, 0, 4096)
// 	for _, cliInbox := range cliInboxSlice {

// 		var smsdata message.SmsData
// 		smsdata.Data = cliInbox.Data
// 		smsdata.FromClientId = cliInbox.FromClientId
// 		smsdata.FromClientName = cliInbox.FromClientName
// 		smsdata.Time, _ = time.Parse("2006-01-02 15:04:05", cliInbox.Time)

// 		smsDataSlice = append(smsDataSlice, &smsdata)
// 	}
// 	return
// }

// func (this *ClientProcess) logoutProc(withComplete bool) (err error) {

// 	var logout_d message.LogoutData
// 	err = json.Unmarshal(this.Msg.Data, &logout_d)
// 	if err != nil {
// 		tools.LOG("Receive Logout message err: %v", err)
// 		return
// 	}

// 	this.ClientData.State = message.Offline
// 	//1. Delete this client in ClientMgr
// 	CliMgr.OnlineCliDelete(this)

// 	//2. Send UpdateClientState to other online clients
// 	CliMgr.UpdateClientStateOnline(this.ClientData)

// 	//3. Send LogoutComplete to this client
// 	if withComplete {
// 		this.write_logout_complete()
// 	}
// 	fmt.Printf("Client logout: %s (%s) Conn: %v", this.ClientData.Name, this.ClientData.Id, this.Conn)
// 	return
// }

// func (this *ClientProcess) write_login_rej(cause string) {
// 	login_rej_data := message.LoginRejectData{
// 		Cause: cause,
// 	}

// 	rstMsg := message.MessageDataW{
// 		Type: message.LoginReject,
// 		Data: login_rej_data,
// 	}

// 	processpkg.WritePkg(this.Conn, &rstMsg)

// }

// func (this *ClientProcess) write_logout_complete() {
// 	logout_complete_data := message.LogoutCompleteData{}

// 	rstMsg := message.MessageDataW{
// 		Type: message.LogoutComplete,
// 		Data: logout_complete_data,
// 	}

// 	processpkg.WritePkg(this.Conn, &rstMsg)

// }

// func (this *ClientProcess) write_login_accept(inboxSmsDataSlice []*message.SmsData) {

// 	login_accept_data := message.LoginAcceptData{
// 		Cause:             "Login Accept",
// 		MyselfClient:      this.ClientData,
// 		ClientOnline:      CliMgr.OnlineClientsForClient,
// 		InboxSmsDataSlice: inboxSmsDataSlice,
// 	}

// 	rstMsg := message.MessageDataW{
// 		Type: message.LoginAccept,
// 		Data: login_accept_data,
// 	}
// 	tools.LOG("inboxSmsDataSlice: %v", inboxSmsDataSlice)
// 	processpkg.WritePkg(this.Conn, &rstMsg)

// }

// func (this *ClientProcess) write_register_rej(cause string) {
// 	register_rej_data := message.RegisterRejectData{
// 		Cause: cause,
// 	}

// 	rstMsg := message.MessageDataW{
// 		Type: message.LoginReject,
// 		Data: register_rej_data,
// 	}
// 	tools.LOG("用户注册失败: %v", cause)
// 	processpkg.WritePkg(this.Conn, &rstMsg)

// }

// func (this *ClientProcess) write_register_accept() {
// 	register_accept_data := message.RegisterAcceptData{
// 		Cause: "注册成功",
// 	}

// 	rstMsg := message.MessageDataW{
// 		Type: message.RegisterAccept,
// 		Data: register_accept_data,
// 	}
// 	tools.LOG("用户注册成功")
// 	processpkg.WritePkg(this.Conn, &rstMsg)

// }

// func (this *ClientProcess) receiveSearchClientRequest() (err error) {
// 	var scr_req_d message.SearchClientRequestData
// 	err = json.Unmarshal(this.Msg.Data, &scr_req_d)
// 	if err != nil {
// 		tools.LOG("receiveSearchClientRequest json.Unmarshal fail")
// 		this.write_search_client_fail(message.SERVER_ERROR)
// 		return
// 	}

// 	err = model.CrDB.CountExist(scr_req_d.ClientId)
// 	if err != nil {
// 		this.write_search_client_fail(message.COUNT_NOT_EXIST)
// 		return
// 	}

// 	count, err := model.CrDB.CountSearch(scr_req_d.ClientId)
// 	if err != nil {
// 		this.write_search_client_fail(message.SERVER_ERROR)
// 		return
// 	}

// 	cli_data := message.ClientUserData{
// 		Id:    count.CountId,
// 		Name:  count.CountName,
// 		State: message.Offline,
// 	}

// 	// cli_data, ok := CliMgr.OnlineClientsForClient[scr_req_d.ClientId]
// 	// if !ok {
// 	// 	this.write_search_client_fail(message.COUNT_NOT_ONLINE)
// 	// 	return
// 	// }

// 	this.write_search_client_response(&cli_data)
// 	return
// }

// func (this *ClientProcess) write_search_client_fail(cause string) {

// 	search_client_fail_data := message.SearchClientFailData{
// 		Cause: cause,
// 	}

// 	rstMsg := message.MessageDataW{
// 		Type: message.SearchClientFail,
// 		Data: search_client_fail_data,
// 	}
// 	tools.LOG("write_search_client_fail: %v", rstMsg)
// 	processpkg.WritePkg(this.Conn, &rstMsg)

// }

// func (this *ClientProcess) write_search_client_response(cd *message.ClientUserData) {
// 	tools.LOG("write_search_client_fail")
// 	search_client_response_data := message.SearchClientResponseData{
// 		ClientUserData: *cd,
// 	}

// 	rstMsg := message.MessageDataW{
// 		Type: message.SearchClientResponse,
// 		Data: search_client_response_data,
// 	}
// 	processpkg.WritePkg(this.Conn, &rstMsg)

// }

// func (this *ClientProcess) receiveCheckInboxNotify() (err error) {
// 	var cin_d message.CheckInboxNotifyData
// 	err = json.Unmarshal(this.Msg.Data, &cin_d)
// 	if err != nil {
// 		tools.LOG("receiveCheckInboxNotify json.Unmarshal fail")
// 		return
// 	}

// 	err = model.CrDB.ClientInboxDelete(cin_d.ClientId)
// 	if err != nil {
// 		tools.LOG("ClientInboxDelete fail, ClientId: %v", cin_d.ClientId)
// 		return
// 	}

// 	return
// }

// func (this *ClientProcess) receiveGetSmsHistoryRequest() (err error) {
// 	var gshr_d message.GetSmsHistoryRequestData
// 	err = json.Unmarshal(this.Msg.Data, &gshr_d)
// 	if err != nil {
// 		tools.LOG("receiveCheckInboxNotify json.Unmarshal fail")
// 		return
// 	}

// 	var sum uint32
// 	sum = 0
// 	sum_response := false
// 	if gshr_d.SumRequest {
// 		sum = model.CrDB.SmsHistoryCount(this.ClientData.Id, gshr_d.ClientId)
// 		sum_response = true
// 	}

// 	smsDataSlice := this.getSmsHistorySlice(gshr_d.ClientId, gshr_d.Number, gshr_d.Offset)

// 	this.writeGetSmsHistoryReponse(smsDataSlice, sum_response, sum)

// 	return
// }

// func (this *ClientProcess) writeGetSmsHistoryReponse(smsDataSlice []*message.SmsData, sum_response bool, sum uint32) (err error) {

// 	search_client_response_data := message.GetSmsHistoryResponseData{
// 		SmsDataSlice: smsDataSlice,
// 		SumNumber:    sum,
// 		SumResponse:  sum_response,
// 	}

// 	rstMsg := message.MessageDataW{
// 		Type: message.GetSmsHistoryResponse,
// 		Data: search_client_response_data,
// 	}
// 	processpkg.WritePkg(this.Conn, &rstMsg)
// 	return
// }

// func (this *ClientProcess) getSmsHistorySlice(toClientId string, number int, offset uint32) (smsDataSlice []*message.SmsData) {

// 	smsHistorySlice := model.CrDB.SmsHistorySelect(this.ClientData.Id, toClientId, number, offset)

// 	smsDataSlice = make([]*message.SmsData, 0, 4096)
// 	for _, smsHistory := range smsHistorySlice {

// 		var smsdata message.SmsData
// 		smsdata.Data = smsHistory.Data
// 		smsdata.FromClientId = smsHistory.FromClientId
// 		smsdata.FromClientName = smsHistory.FromClientName
// 		smsdata.Time, _ = time.Parse("2006-01-02 15:04:05", smsHistory.Time)

// 		smsDataSlice = append(smsDataSlice, &smsdata)
// 	}
// 	return
// }
