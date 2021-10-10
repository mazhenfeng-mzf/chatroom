package process

import (
	"chatroom/public/message"
)

type ClientsMgr struct {
	OnlineClients map[string]*ClientProcess
	AllClients    map[string]*ClientProcess

	OnlineClientsForClient map[string]*message.ClientUserData

	PublicChatroomClients map[string]*ClientProcess

	PublicChatroomClientsForClient map[string]*message.ClientUserData
}

var CliMgr ClientsMgr

func (thiss *ClientsMgr) Init() {
	//clients map stored and used in server local
	CliMgr.OnlineClients = make(map[string]*ClientProcess, 4096)
	CliMgr.AllClients = make(map[string]*ClientProcess, 4096)
	CliMgr.PublicChatroomClients = make(map[string]*ClientProcess, 4096)

	//When a client login, would send this message(current online clients) to it
	CliMgr.OnlineClientsForClient = make(map[string]*message.ClientUserData, 4096)

	///When a client enter PublicChatroom, would send this message(current inPC clients) to it
	CliMgr.PublicChatroomClientsForClient = make(map[string]*message.ClientUserData, 4096)
}

func (thiss *ClientsMgr) OnlineCliAdd(cp *ClientProcess) {
	thiss.OnlineClients[cp.Client_id] = cp

	// var olCliforClient message.ClientUserData
	// olCliforClient.Id = cp.Client_id
	// olCliforClient.Name = cp.Client_name
	// thiss.OnlineClientsForClient[cp.Client_id] = &olCliforClient
}

func (thiss *ClientsMgr) OnlineCliDelete(clientId string) {
	delete(thiss.OnlineClients, clientId)
	// delete(thiss.OnlineClientsForClient, clientId)
}

func (thiss *ClientsMgr) PublicCrCliAdd(cp *ClientProcess) {
	thiss.PublicChatroomClients[cp.Client_id] = cp

	var puCliForClient message.ClientUserData
	puCliForClient.Id = cp.Client_id
	puCliForClient.Name = cp.Client_name
	thiss.PublicChatroomClientsForClient[cp.Client_id] = &puCliForClient
}

func (thiss *ClientsMgr) PublicCrCliDelete(cp *ClientProcess) {
	delete(thiss.PublicChatroomClients, cp.Client_id)
	delete(thiss.PublicChatroomClientsForClient, cp.Client_id)
}

func (this *ClientsMgr) ClientIsOnline(clientId string) (cliProc *ClientProcess, yes bool) {
	cliProc, yes = this.OnlineClients[clientId]
	return
}

// func (thiss *ClientsMgr) UpdateClientStateOnline(cd *message.ClientUserData) {

// 	udMsg := message.MessageDataW{
// 		Type: message.UpdateClientState,
// 		Data: cd,
// 	}
// 	MyLOG.Log("update client %s (%s) state %s", cd.Name, cd.Id, cd.State)

// 	for id, value := range thiss.OnlineClients {
// 		if cd.Id == id {
// 			continue
// 		}
// 		processpkg.WritePkg(value.Conn, &udMsg)
// 	}

// }

// func (thiss *ClientsMgr) UpdateClientStateInPC(cd *message.ClientUserData) {

// 	udMsg := message.MessageDataW{
// 		Type: message.UpdatePublicChatroomClients,
// 		Data: cd,
// 	}
// 	MyLOG.Log("update PublicChatroom client %s (%s) state %s", cd.Name, cd.Id, cd.State)

// 	for id, value := range thiss.PublicChatroomClients {
// 		if cd.Id == id {
// 			continue
// 		}
// 		processpkg.WritePkg(value.Conn, &udMsg)
// 	}

// }
