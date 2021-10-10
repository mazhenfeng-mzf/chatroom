package message

//Message Type const
const (
	LoginRequest                          uint8 = 1
	LoginAccept                           uint8 = 2
	LoginReject                           uint8 = 3
	RegisterRequest                       uint8 = 4
	RegisterAccept                        uint8 = 5
	RegisterReject                        uint8 = 6
	DeRegister                            uint8 = 7
	Logout                                uint8 = 8
	EnterPCRequest                        uint8 = 9
	EnterPCAccept                         uint8 = 10
	ExitPC                                uint8 = 11
	UpdatePublicChatClients               uint8 = 12
	SearchClientResquest                  uint8 = 13
	SearchClientResponse                  uint8 = 14
	SearchClientFail                      uint8 = 15
	FriendsAddRequest                     uint8 = 16
	FriendsAddAccept                      uint8 = 17
	FriendsAddReject                      uint8 = 18
	FriendsAddNotify                      uint8 = 19
	SmMessageReadNotify_FriendsAddRequest uint8 = 20
	SmMessageReadNotify_FriendsAddAccept  uint8 = 21
	SmMessageReadNotify_FriendsAddReject  uint8 = 22
	MsgInboxCheckNotify                   uint8 = 23
	SmsHistoryRequest                     uint8 = 24
	SmsHistoryResponse                    uint8 = 25

	FriendOnOffLine uint8 = 24
)

//message for socket transfer
type MessageSocket struct {
	Type   uint8
	Length uint32
	Data   []byte
}

//read message from socket, use this to store
type MessageDataR struct {
	Type uint8
	Data []byte
}

//use this to store message which will write to socket
type MessageDataW struct {
	Type uint8
	Data interface{}
}

//Following is the MessageDataW.Data
type LoginRequestData struct {
	Id  string
	Pwd string
}

type LoginAcceptData struct {
	MyselfClient *ClientUserData
	MsgInbox     []*SmsP2PData

	Slice_FriendsAddRequest []*FriendsAddRequestData
	Slice_FriendsAddAccept  []*FriendsAddAcceptData
	Slice_FriendsAddReject  []*FriendsAddRejectData
	//ClientOnline map[string]*ClientUserData

	Friends []*ClientUserData
}

type LoginRejectData struct {
	Cause CauseId
}

type RegisterRequestData struct {
	Id   string
	Name string
	Pwd  string
}

type RegisterAcceptData struct {
}

type RegisterRejectData struct {
	Cause CauseId
}

type LogoutData struct {
	Id string
}

type EnterPCData struct {
}

type ExitPCData struct {
}

type EnterPCAcceptData struct {
	PcClientList map[string]*ClientUserData
}

type UpdatePublicChatClientsData struct {
	Id   string
	Name string
	InPC bool
}

type SearchClientResquestData struct {
	Id string
}

type SearchClientResponseData struct {
	ClientUserData
}

type SearchClientFailData struct {
	CauseId
}

type FriendsAddRequestData struct {
	FromClientId   string
	FromClientName string
	ToClientId     string
	ToClientName   string
	Message        string
}

type FriendsAddAcceptData struct {
	FromClientId   string // who trigger FriendsAddRequest
	FromClientName string
	ToClientId     string
	ToClientName   string
	Online         bool
}

type FriendsAddNotifyData struct {
	FromClientId   string // who trigger FriendsAddRequest
	FromClientName string
	ToClientId     string
	ToClientName   string
	Online         bool
}

type FriendsAddRejectData struct {
	FromClientId   string // who trigger FriendsAddRequest
	FromClientName string
	ToClientId     string
	ToClientName   string
}

type SmMessageReadNotify_FriendsAddRequest_Data struct {
	FriendsAddRequestData
}

type SmMessageReadNotify_FriendsAddAccept_Data struct {
	FromClientId   string
	FromClientName string
	ToClientId     string
	ToClientName   string
}

type SmMessageReadNotify_FriendsAddReject_Data struct {
	FromClientId   string
	FromClientName string
	ToClientId     string
	ToClientName   string
}

type MsgInboxCheckNotifyData struct {
	FromClientId string
	ToClientId   string
}

type FriendOnOffLineData struct {
	ClientUserData
}

type SmsHistoryRequestData struct {
	FromClientId string
	ToClientId   string
	Offset       uint32
	Number       int
}

type SmsHistoryResponseData struct {
	FromClientId    string
	ToClientId      string
	SmsHistorySlice []*SmsP2PData
	SmsHistorySum   uint32
}
