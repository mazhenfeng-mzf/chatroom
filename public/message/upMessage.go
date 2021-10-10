package message

import "time"

//Message Type const
const (
	Sms               uint8 = 100
	SmsPublicChatroom uint8 = 101
	SmsP2P            uint8 = 102
)

// type SmsData struct {
// 	FromClientId string
// 	ToClientId   string
// 	Data         string
// }

type SmsPublicChatroomData struct {
	FromClientId   string
	FromClientName string
	Data           string
	Time           time.Time
}

type SmsP2PData struct {
	FromClientId   string
	FromClientName string
	ToClientId     string
	ToClientName   string
	Data           string
	Time           string
}
