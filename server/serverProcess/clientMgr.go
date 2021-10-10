// package serverProcess

// import (
// 	"chatroom/Public/message"
// 	"chatroom/Public/processpkg"
// 	"chatroom/Public/tools"
// )

// type ClientsMgr struct {
// 	OnlineClients                  map[string]*ClientProcess
// 	AllClients                     map[string]*ClientProcess
// 	OnlineClientsForClient         map[string]*message.ClientUserData
// 	PublicChatroomClients          map[string]*ClientProcess
// 	PublicChatroomClientsForClient map[string]*message.ClientUserData
// }

// var CliMgr ClientsMgr

// func (thiss *ClientsMgr) Init() {
// 	//clients map stored and used in server local
// 	CliMgr.OnlineClients = make(map[string]*ClientProcess, 4096)
// 	CliMgr.AllClients = make(map[string]*ClientProcess, 4096)
// 	CliMgr.PublicChatroomClients = make(map[string]*ClientProcess, 4096)

// 	//When a client login, would send this message(current online clients) to it
// 	CliMgr.OnlineClientsForClient = make(map[string]*message.ClientUserData, 4096)
// 	CliMgr.PublicChatroomClientsForClient = make(map[string]*message.ClientUserData, 4096)
// }

// func (thiss *ClientsMgr) OnlineCliAdd(cp *ClientProcess) {
// 	thiss.OnlineClients[cp.ClientData.Id] = cp

// 	var olCliforClient message.ClientUserData
// 	olCliforClient.Id = cp.ClientData.Id
// 	olCliforClient.Name = cp.ClientData.Name
// 	olCliforClient.State = message.Online
// 	thiss.OnlineClientsForClient[cp.ClientData.Id] = &olCliforClient
// }

// func (thiss *ClientsMgr) OnlineCliDelete(cp *ClientProcess) {
// 	delete(thiss.OnlineClients, cp.ClientData.Id)
// 	delete(thiss.OnlineClientsForClient, cp.ClientData.Id)
// }

// func (thiss *ClientsMgr) PublicCrCliAdd(cp *ClientProcess) {
// 	thiss.PublicChatroomClients[cp.ClientData.Id] = cp

// 	var puCliForClient message.ClientUserData
// 	puCliForClient.Id = cp.ClientData.Id
// 	puCliForClient.Name = cp.ClientData.Name
// 	thiss.PublicChatroomClientsForClient[cp.ClientData.Id] = &puCliForClient
// }

// func (thiss *ClientsMgr) PublicCrCliDelete(cp *ClientProcess) {
// 	delete(thiss.PublicChatroomClients, cp.ClientData.Id)
// 	delete(thiss.PublicChatroomClientsForClient, cp.ClientData.Id)
// }

// func (thiss *ClientsMgr) UpdateClientStateOnline(cd *message.ClientUserData) {

// 	udMsg := message.MessageDataW{
// 		Type: message.UpdateClientState,
// 		Data: cd,
// 	}
// 	tools.LOG("update client %s (%s) state %s", cd.Name, cd.Id, cd.State)

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
// 	tools.LOG("update PublicChatroom client %s (%s) state %s", cd.Name, cd.Id, cd.State)

// 	for id, value := range thiss.PublicChatroomClients {
// 		if cd.Id == id {
// 			continue
// 		}
// 		processpkg.WritePkg(value.Conn, &udMsg)
// 	}

// }
