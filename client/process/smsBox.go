package process

import (
	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"
)

// 存放 sms 消息的结构体
type SmsBox struct {
	SmsMap map[string][]*message.SmsP2PData //string is Client Id

	SmsSum uint16
}

var MySmsBox SmsBox

func (this *SmsBox) Init() {

	// smsList := make([]*message.SmsP2PData, 0, 4096)
	this.SmsMap = make(map[string][]*message.SmsP2PData, 4096)
	this.SmsSum = 0
}

func (this *SmsBox) Add(smsP2p_d *message.SmsP2PData) {

	//MyLOG.Log("SmsBox Before Add: %#v", this.SmsMap)

	fromClientId := smsP2p_d.FromClientId

	smsList, exist := this.SmsMap[fromClientId]
	//MyLOG.Log("smsList Before Add: %#v", smsList)
	if !exist {
		smsList = make([]*message.SmsP2PData, 0, 4096)

		MyView.SmsBox_forClient_Create(fromClientId)
	}

	smsList = append(smsList, smsP2p_d)
	//MyLOG.Log("smsList After Add: %#v", smsList)

	this.SmsMap[fromClientId] = smsList
	MyLOG.Log("SmsBox After Add: %#v", this.SmsMap)

	this.SmsSum += 1
}

func (this *SmsBox) Add_Slice(smsP2p_d_slice []*message.SmsP2PData) {
	for _, smsP2p := range smsP2p_d_slice {
		this.Add(smsP2p)
	}
	MyLOG.Log("SmsBox After Add_Slice: %#v", this.SmsMap)
}

func (this *SmsBox) Get(fromClientId string) (smsList []*message.SmsP2PData) {

	smsList = this.SmsMap[fromClientId]

	smsLen := len(smsList)

	delete(this.SmsMap, fromClientId)

	this.SmsSum -= uint16(smsLen)

	return
}

func (this *SmsBox) Display_MainPage() {
	MyView.MainPage_SmsBox.SetText(fmt.Sprintf("收件箱(%d)", this.SmsSum))
}

func (this *SmsBox) ClickSmsInbox(fromClientId string, fromClientName string) {

	// clear MyView.SmsBoxMap
	MyView.SmsBoxMap[fromClientId].PushBtnAssign.SetText(fmt.Sprintf("%s(%s)(%d)", fromClientName, fromClientId, 0)) //更新收件箱
	MyView.SmsBox_forClient_Delete(fromClientId)

	// clear MySms.SmsMap
	smsList := this.Get(fromClientId)
	this.Display_MainPage() //刷新主界面
	MyView.OpenP2PChatRoom(fromClientId, fromClientName, smsList)
}
