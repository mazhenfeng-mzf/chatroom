package process

import (
	"chatroom/public/errs"
	"chatroom/public/message"
	. "chatroom/public/tools"
	"encoding/json"
)

// func (this *ServerProcess) receiveSm() (err error) {
// 	//Data 反序列化

// 	//如果是 FriendAddAccept 消息, 则更新好友列表

// 	MySmMsgBox.Add(this.Msg)
// 	MyLOG.Log("finish MySmMsgBox.Add")
// 	MySmMsgBox.Display_MainPage()

// 	return
// }

func (this *ServerProcess) receive_FriendsAddRequest() (err error) {

	var data message.FriendsAddRequestData
	err = json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("FriendsAddRequest 信息解析失败")
		err = errs.INVAILD_MSG_PARSE_FAIL
		return
	}

	MySmMsgBox.Add(message.FriendsAddRequest, data)
	MySmMsgBox.Display_MainPage()
	return
}

func (this *ServerProcess) receive_FriendsAddAccept() (err error) {

	var data message.FriendsAddAcceptData
	err = json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("FriendsAddAccept 信息解析失败")
		err = errs.INVAILD_MSG_PARSE_FAIL
		return
	}

	CliMgr.FriendsClientsAdd(data.ToClientId, data.ToClientName, data.Online)

	MySmMsgBox.Add(message.FriendsAddAccept, data)
	MySmMsgBox.Display_MainPage()
	return
}

func (this *ServerProcess) receive_FriendsAddNotify() (err error) {

	var data message.FriendsAddNotifyData
	err = json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("FriendsAddNotify 信息解析失败")
		err = errs.INVAILD_MSG_PARSE_FAIL
		return
	}

	CliMgr.FriendsClientsAdd(data.FromClientId, data.FromClientName, data.Online)

	MySmMsgBox.Add(message.FriendsAddAccept, data)
	MySmMsgBox.Display_MainPage()
	return
}

func (this *ServerProcess) receive_FriendsAddReject() (err error) {

	var data message.FriendsAddRejectData
	err = json.Unmarshal(this.Msg.Data, &data)
	if err != nil {
		MyLOG.ErrLog("FriendsAddReject 信息解析失败")
		err = errs.INVAILD_MSG_PARSE_FAIL
		return
	}

	MySmMsgBox.Add(message.FriendsAddReject, data)
	MySmMsgBox.Display_MainPage()
	return
}
