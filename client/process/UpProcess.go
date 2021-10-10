package process

import (
	"chatroom/public/message"
	"chatroom/public/processpkg"
	. "chatroom/public/tools"
	"strings"
)

func (this *ClientProcess) SendPCSms(data string) (err error) {

	data = strings.TrimSpace(data)
	if data == "" {
		MyView.PcCR.PcInput.SetText("")
	}

	//构建登录消息
	smsPCData := message.SmsPublicChatroomData{
		FromClientId:   this.Id,
		FromClientName: this.Name,
		Data:           data,
	}

	smsPCMsg := message.MessageDataW{
		Type: message.SmsPublicChatroom,
		Data: smsPCData,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &smsPCMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", smsPCMsg)
	MyView.PcCR.PcInput.SetText("")
	return
}

func (this *ClientProcess) SendP2PSms(toClientId string, data string) (err error) {

	data = strings.TrimSpace(data)
	if data == "" {
		MyView.PcCR.PcInput.SetText("")
	}

	//构建登录消息
	smsP2pData := message.SmsP2PData{
		FromClientId: this.Id,
		ToClientId:   toClientId,
		Data:         data,
	}

	smsPCMsg := message.MessageDataW{
		Type: message.SmsP2P,
		Data: smsP2pData,
	}

	//发送消息
	err = processpkg.WritePkg(this.Conn, &smsPCMsg)
	if err != nil {
		MyLOG.ErrLog("WritePkg err=%v", err)
		return
	}
	MyLOG.Log("发送消息成功: %#v", smsPCMsg)
	MyView.P2PChatRoomList[toClientId].Input.SetText("")
	return
}
