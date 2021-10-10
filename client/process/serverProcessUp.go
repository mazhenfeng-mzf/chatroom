package process

import (
	"chatroom/public/errs"
	"chatroom/public/message"
	. "chatroom/public/tools"
	"encoding/json"
)

func (this *ServerProcess) receiveSmsPublicChatroom() (err error) {
	//Data 反序列化

	var data message.SmsPublicChatroomData
	err = json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("public 信息解析失败")
		err = errs.INVAILD_MSG_PARSE_FAIL
		return
	}
	MyLOG.Log("收到 公共聊天消息 : %#v", data)

	MyView.PcSmsListUpdate(&data)
	MyView.DisplayPcSmsList()
	return
}

func (this *ServerProcess) receiveSmsP2P() (err error) {
	//Data 反序列化

	var data message.SmsP2PData
	err = json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("p2p 信息解析失败")
		err = errs.INVAILD_MSG_PARSE_FAIL
		return
	}
	MyLOG.Log("收到来自 %s(%s) 的p2p 消息 : ", data.FromClientName, data.FromClientId, data.Data)
	MyLOG.Log("ServerProcess.Id: %s", this.Id)
	if data.FromClientId == this.Id {
		MyLOG.Log("收到我自己的消息")
		MyView.P2pSmsListUpdate(data.ToClientId, &data)
		MyView.DisplayP2pSmsList(data.ToClientId)
		return
	}

	if CliMgr.IsMyFriends(data.FromClientId) {
		data.FromClientName = CliMgr.MyFriendsClients[data.FromClientId].Name
	} else {
		data.FromClientName = ""
	}

	if _, yes := MyView.P2pChatroomIsOpen(data.FromClientId); yes {
		MyLOG.Log("当前有打开与 %s 的P2P聊天室", data.FromClientId)
		MyView.P2pSmsListUpdate(data.FromClientId, &data)
		MyView.DisplayP2pSmsList(data.FromClientId)
	} else {
		MyLOG.Log("当前没有打开与 %s 的P2P聊天室", data.FromClientId)
		MySmsBox.Add(&data)
		MySmsBox.Display_MainPage()
	}

	return
}
