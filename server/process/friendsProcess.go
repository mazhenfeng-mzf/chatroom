package process

import (
	"chatroom/public/message"
	"chatroom/public/processpkg"
	. "chatroom/public/tools"
	"chatroom/server/model"
	"encoding/json"
	"time"
)

func (this *ClientProcess) receiveFriendsAddRequest() (causeId message.CauseId) {

	//Data 反序列化
	var req_d message.FriendsAddRequestData
	err := json.Unmarshal(this.Msg.Data, &req_d)
	if err != nil {
		MyLOG.ErrLog("FriendsAddRequest信息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		return
	}

	msg := message.MessageDataW{
		Type: message.FriendsAddRequest,
		Data: req_d,
	}
	toClientId := req_d.ToClientId
	fromClientId := req_d.FromClientId
	farMsg := req_d.Message

	if cliProc, yes := CliMgr.ClientIsOnline(toClientId); yes {
		MyLOG.Log("系统消息目标用户id %s 在线", toClientId)
		processpkg.WritePkg(cliProc.Conn, &msg)
	} else {
		MyLOG.Log("系统消息目标用户id %s 不在线, 消息存进数据库", toClientId)
		model.CrDB.FriendsBoxAdd(fromClientId, toClientId, farMsg, message.FriendsAddRequest)
	}

	return
}

func (this *ClientProcess) receiveFriendsAddAccept() (causeId message.CauseId) {

	//Data 反序列化
	var acc_d message.FriendsAddAcceptData
	err := json.Unmarshal(this.Msg.Data, &acc_d)
	if err != nil {
		MyLOG.ErrLog("FriendsAddAccept信息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		return
	}

	fromClientId := acc_d.FromClientId
	//toClientId := acc_d.ToClientId
	toCountId := this.Model_Id
	timeString := time.Now().Format("2006-01-02 15:04:05")

	MyLOG.Log("model.CrDB.FriendsAdd(%d, %s, %s)", toCountId, fromClientId, timeString)
	causeId = model.CrDB.FriendsAdd(toCountId, fromClientId, timeString)
	if causeId != 0 {
		MyLOG.ErrLog("FriendsInsert fail")
		return
	}

	acc_d.ToClientName = this.Client_name

	//check if online
	_, online := CliMgr.ClientIsOnline(acc_d.ToClientId)
	acc_d.Online = online

	//transfer the FriendsAddAccept to peer
	this.handle_FriendsAddAccept(&acc_d)
	// msg := message.MessageDataW{
	// 	Type: message.FriendsAddAccept,
	// 	Data: acc_d,
	// }
	// if cliProc, yes := CliMgr.ClientIsOnline(fromClientId); yes {
	// 	MyLOG.Log("系统消息目标用户id %s 在线", fromClientId)
	// 	processpkg.WritePkg(cliProc.Conn, &msg)
	// } else {
	// 	MyLOG.Log("系统消息目标用户id %s 不在线", fromClientId)
	// 	model.CrDB.FriendsBoxAdd(fromClientId, toClientId, "", message.FriendsAddAccept)
	// }

	return
}

func (this *ClientProcess) handle_FriendsAddAccept(acc_d *message.FriendsAddAcceptData) (causeId message.CauseId) {

	//转发 FriendsAddAccept 消息给目标
	acc_msg := message.MessageDataW{
		Type: message.FriendsAddAccept,
		Data: acc_d,
	}

	cliProc, yes := CliMgr.ClientIsOnline(acc_d.FromClientId)
	if yes {
		MyLOG.Log("系统消息目标用户id %s 在线", acc_d.FromClientId)
		processpkg.WritePkg(cliProc.Conn, &acc_msg)
	} else {
		MyLOG.Log("系统消息目标用户id %s 不在线", acc_d.FromClientId)
		model.CrDB.FriendsBoxAdd(acc_d.FromClientId, acc_d.ToClientId, "", message.FriendsAddAccept)
	}

	notify_data := message.FriendsAddNotifyData{
		FromClientId:   acc_d.FromClientId,
		FromClientName: acc_d.FromClientName,
		ToClientId:     acc_d.ToClientId,
		ToClientName:   acc_d.ToClientName,
		Online:         yes,
	}

	//回复 FriendsAddNotify 给用户
	notify_msg := message.MessageDataW{
		Type: message.FriendsAddNotify,
		Data: notify_data,
	}
	processpkg.WritePkg(this.Conn, &notify_msg)

	return
}

func (this *ClientProcess) receiveFriendsAddReject() (causeId message.CauseId) {

	//Data 反序列化
	var rej_d message.FriendsAddRejectData
	err := json.Unmarshal(this.Msg.Data, &rej_d)
	if err != nil {
		MyLOG.ErrLog("FriendsAddReject信息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		return
	}

	fromClientId := rej_d.FromClientId
	toClientId := rej_d.ToClientId
	toCountId := this.Model_Id
	timeString := time.Now().Format("2006-01-02 15:04:05")

	causeId = model.CrDB.FriendsAdd(toCountId, fromClientId, timeString)
	if causeId != 0 {
		MyLOG.ErrLog("FriendsInsert fail")
		return
	}

	rej_d.ToClientName = this.Client_name
	msg := message.MessageDataW{
		Type: message.FriendsAddReject,
		Data: rej_d,
	}

	if cliProc, yes := CliMgr.ClientIsOnline(fromClientId); yes {
		MyLOG.Log("系统消息目标用户id %s 在线", fromClientId)
		processpkg.WritePkg(cliProc.Conn, &msg)
	} else {
		MyLOG.Log("系统消息目标用户id %s 不在线", fromClientId)
		model.CrDB.FriendsBoxAdd(fromClientId, toClientId, "", message.FriendsAddReject)
	}

	return
}

func handleSmMessageReadNotify_FriendsAddRequest(msg message.FriendsAddRequestData) {
	model.CrDB.FriendsBoxDelete(msg.FromClientId, msg.ToClientId, msg.Message, message.FriendsAddRequest)
}

func handleSmMessageReadNotify_FriendsAddAccept(msg message.FriendsAddAcceptData) {
	model.CrDB.FriendsBoxDelete(msg.FromClientId, msg.ToClientId, "", message.FriendsAddAccept)
}

func handleSmMessageReadNotify_FriendsAddReject(msg message.FriendsAddRejectData) {
	model.CrDB.FriendsBoxDelete(msg.FromClientId, msg.ToClientId, "", message.FriendsAddReject)
}
