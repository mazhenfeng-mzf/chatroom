package process

import (
	"chatroom/public/message"
	"chatroom/public/processpkg"
	. "chatroom/public/tools"
	"chatroom/server/model"
	"encoding/json"
	"time"
)

func (this *ClientProcess) ReceiveSmsPCMsg() (causeId message.CauseId) {

	//Data 反序列化
	var smsPC_d message.SmsPublicChatroomData
	err := json.Unmarshal(this.Msg.Data, &smsPC_d)
	if err != nil {
		MyLOG.ErrLog("Public Sms 信息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		this.write_login_rej(causeId)
		return
	}
	smsPC_d.Time = time.Now()

	this.PcBroadcastSms(&smsPC_d)

	return
}

func (this *ClientProcess) PcBroadcastSms(smsPcData *message.SmsPublicChatroomData) (causeId message.CauseId) {

	Msg := message.MessageDataW{
		Type: message.SmsPublicChatroom,
		Data: smsPcData,
	}

	for _, cliPro := range CliMgr.PublicChatroomClients {

		processpkg.WritePkg(cliPro.Conn, &Msg)

		MyLOG.Log("广播消息给用户 %s: %s", cliPro.Client_name, smsPcData.Data)
	}
	MyLOG.Log("广播完成: %s", smsPcData.Data)
	return
}

func (this *ClientProcess) ReceiveSmsP2PMsg() (causeId message.CauseId) {

	//Data 反序列化
	var smsP2p_d message.SmsP2PData
	err := json.Unmarshal(this.Msg.Data, &smsP2p_d)
	if err != nil {
		MyLOG.ErrLog("p2p Sms 信息解析失败")
		causeId = message.ID_INVAILD_MSG_PARSE_FAIL
		this.write_login_rej(causeId)
		return
	}
	toClientId := smsP2p_d.ToClientId
	fromClientId := smsP2p_d.FromClientId
	smsP2p_d.FromClientName = CliMgr.OnlineClients[fromClientId].Client_name
	smsP2p_d.Time = time.Now().Format("2006-01-02 15:04:05")
	//fromClientName := CliMgr.OnlineClients[fromClientId].Client_name

	model.CrDB.SmsHistoryAdd(fromClientId, toClientId, smsP2p_d.Data, smsP2p_d.Time)

	CliMgr.OnlineClients[fromClientId].writeSmsP2pMsg(&smsP2p_d)

	if cliProc, yes := CliMgr.ClientIsOnline(toClientId); yes {
		MyLOG.Log("Sms 消息目标用户id %s 在线", toClientId)
		cliProc.writeSmsP2pMsg(&smsP2p_d)
	} else {
		MyLOG.Log("Sms 消息目标用户id %s 不在线", toClientId)
		//MyLOG.Log(fromClientId, this.Model_Id, toClientId)
		model.CrDB.ClientInboxInsert(fromClientId, this.Model_Id, toClientId, smsP2p_d.Data, smsP2p_d.Time)
		//MyLOG.Log("完成 ClientInboxInsert: %s", toClientId)
	}

	return
}

func (this *ClientProcess) writeSmsP2pMsg(smsP2p_d *message.SmsP2PData) {

	Msg := message.MessageDataW{
		Type: message.SmsP2P,
		Data: smsP2p_d,
	}

	processpkg.WritePkg(this.Conn, &Msg)

	MyLOG.Log("发送 p2p 消息给: %s", this.Client_name)
	return
}
