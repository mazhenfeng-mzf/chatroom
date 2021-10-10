package process

import (
	"chatroom/public/message"
)

type ClientsMgr struct {
	//OnlineClients map[string]*message.ClientUserData
	//AllClients    map[string]*message.ClientUserData

	PublicChatroomClients map[string]*message.ClientUserData

	MyFriendsClients map[string]*message.ClientUserData
}

var CliMgr ClientsMgr

func (thiss *ClientsMgr) Init() {
	thiss.MyFriendsClients = make(map[string]*message.ClientUserData, 4096)
	//CliMgr.OnlineClients = make(map[string]*message.ClientUserData, 4096)
	//CliMgr.AllClients = make(map[string]*message.ClientUserData, 4096)
	//CliMgr.PublicChatroomClients = make(map[string]*message.ClientUserData, 4096)
}

func (thiss *ClientsMgr) PublicChatroomClientsInit(clientsList map[string]*message.ClientUserData) {
	thiss.PublicChatroomClients = make(map[string]*message.ClientUserData, 4096)
	thiss.PublicChatroomClients = clientsList
}

func (thiss *ClientsMgr) PublicChatroomClientsUpdate(clients message.UpdatePublicChatClientsData) {
	if clients.InPC {
		cud := message.ClientUserData{
			Id:   clients.Id,
			Name: clients.Name,
		}
		thiss.PublicChatroomClients[clients.Id] = &cud
	} else {
		delete(thiss.PublicChatroomClients, clients.Id)
	}
}

func (thiss *ClientsMgr) FriendsClientsAdd(clientId string, clientName string, online bool) {

	cud := message.ClientUserData{
		Id:     clientId,
		Name:   clientName,
		Online: online,
	}

	thiss.MyFriendsClients[clientId] = &cud
}

func (thiss *ClientsMgr) FriendsClientsAdd_Slice(slice_CUD []*message.ClientUserData) {
	for _, ucd := range slice_CUD {
		thiss.FriendsClientsAdd(ucd.Id, ucd.Name, ucd.Online)
	}
}

func (thiss *ClientsMgr) FriendsClientsDelete(clientId string) {

	delete(thiss.MyFriendsClients, clientId)
}

func (thiss *ClientsMgr) IsMyFriends(clientId string) (yes bool) {

	_, yes = thiss.MyFriendsClients[clientId]
	return
}

// func (thiss *ClientsMgr) OnlineCliAdd(cd *message.ClientUserData) {
// 	thiss.OnlineClients[cd.Id] = cd
// }

// func (thiss *ClientsMgr) OnlineCliInit(cm map[string]*message.ClientUserData) {
// 	thiss.OnlineClients = cm
// }

// func (thiss *ClientsMgr) OnlineCliDisplay(currentUserId string) {
// 	fmt.Println("--- Current Online Users List ---")
// 	i := 1
// 	for key, value := range thiss.OnlineClients {
// 		if currentUserId == key {
// 			continue
// 		}
// 		fmt.Printf("%d. %s(%s)\n", i, value.Name, key)
// 		i++
// 	}
// 	fmt.Println("--- --- --- --- ---")
// }
