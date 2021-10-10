package model

import (
	//_"github.com/jmoiron/sqlx"

	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type FriendsBox struct {
	Id          int    `db:"id"`
	FromCountId int    `db:"from_count_id"`
	ToCountId   int    `db:"to_count_id"`
	fa_Message  string `db:"message"`
	Type        int    `db:"type"`
}

func (thiss *ChatroomDB) FriendsBoxAdd(fromClientId string, toClientId string, fa_message string, msg_type uint8) (causeId message.CauseId) {
	fromCountId, causeId := thiss.GetCountId(fromClientId)
	if causeId != 0 {
		MyLOG.ErrLog("GetCountId wish clientId \"%s\" fail", fromClientId)
		return
	}

	toCountId, causeId := thiss.GetCountId(toClientId)
	if causeId != 0 {
		MyLOG.ErrLog("GetCountId wish clientId \"%s\" fail", toClientId)
		return
	}

	causeId = thiss.FriendsBoxInsert(fromCountId, toCountId, fa_message, msg_type)
	if causeId != 0 {
		MyLOG.ErrLog("FriendsBoxInsert fail")
		return
	}

	return
}

func (thiss *ChatroomDB) FriendsAddRequestGet(toCountId uint16) (sliceFAR []*message.FriendsAddRequestData, causeId message.CauseId) {

	sliceFAR, causeId = thiss.FriendsAddRequestSelect(toCountId)
	if causeId != 0 {
		MyLOG.ErrLog("FriendsAddRequestSelect fail")
		return
	}
	thiss.FriendsAddRequestDelete(toCountId)
	return
}

func (thiss *ChatroomDB) FriendsAddAcceptGet(fromCountId uint16) (sliceFAA []*message.FriendsAddAcceptData, causeId message.CauseId) {
	sliceFAA, causeId = thiss.FriendsAddAcceptSelect(fromCountId)
	if causeId != 0 {
		MyLOG.ErrLog("FriendsAddAcceptSelect fail")
		return
	}
	thiss.FriendsAddAcceptDelete(fromCountId)
	return
}

func (thiss *ChatroomDB) FriendsAddRejectGet(fromCountId uint16) (sliceFARJ []*message.FriendsAddRejectData, causeId message.CauseId) {
	sliceFARJ, causeId = thiss.FriendsAddRejectSelect(fromCountId)
	if causeId != 0 {
		MyLOG.ErrLog("FriendsAddRejectSelect fail")
		return
	}
	thiss.FriendsAddRejectDelete(fromCountId)
	return
}

func (thiss *ChatroomDB) FriendsBoxInsert(fromCountId uint16, toCountId uint16, fa_message string, msg_type uint8) (causeId message.CauseId) {

	stmt, err := thiss.DB.Prepare(`INSERT friend_box (from_count_id,to_count_id,message,type) values (?,?,?,?)`)
	if err != nil {
		MyLOG.ErrLog("Mysql Prepare fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	_, err = stmt.Exec(fromCountId, toCountId, fa_message, msg_type)
	if err != nil {
		MyLOG.ErrLog("Mysql Exec fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	return
}

func (thiss *ChatroomDB) FriendsAddRequestSelect(toCountId uint16) (sliceFAR []*message.FriendsAddRequestData, causeId message.CauseId) {

	sliceFAR = make([]*message.FriendsAddRequestData, 0, 4096)
	msql := fmt.Sprintf(`
	select t.from_client_id, t.from_client_name, ac.client_id as to_client_id, ac.count_name as to_client_name, t.message from (
		
		select f.*, ac.client_id as from_client_id, ac.count_name as from_client_name
		from friend_box f 
		inner join account ac on f.from_count_id = ac.id
	  )t

	inner join account ac on t.to_count_id = ac.id

	where t.to_count_id = %d and t.type = %d`, toCountId, message.FriendsAddRequest)
	// MyLOG.Log(msql)
	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		var far message.FriendsAddRequestData

		err = rows.Scan(&far.FromClientId, &far.FromClientName, &far.ToClientId, &far.ToClientName, &far.Message)
		if err != nil {
			MyLOG.ErrLog("Scan error: %v", err)
			causeId = message.ID_SERVER_ERROR
			return
		}
		sliceFAR = append(sliceFAR, &far)
	}
	//sort.Sort(sliceMsg)
	MyLOG.Log("FriendsBoxSelect: %v", sliceFAR)
	return
}

func (thiss *ChatroomDB) FriendsAddRequestDelete(toCountId uint16) {

	msql := fmt.Sprintf("Delete from friend_box where to_count_id = ? and type = ?")
	stmt, err := thiss.DB.Prepare(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Prepare fail: %v", err)
		return
	}
	_, err = stmt.Exec(toCountId, message.FriendsAddRequest)
	return
}

func (thiss *ChatroomDB) FriendsAddAcceptSelect(fromCountId uint16) (sliceFAA []*message.FriendsAddAcceptData, causeId message.CauseId) {

	//fromCountId: who trigger the FriendsAddRequest

	sliceFAA = make([]*message.FriendsAddAcceptData, 0, 4096)
	msql := fmt.Sprintf(`
	select t.from_client_id, t.from_client_name, ac.client_id as to_client_id, ac.count_name as to_client_name from (
		
		select f.*, ac.client_id as from_client_id, ac.count_name as from_client_name
		from friend_box f 
		inner join account ac on f.from_count_id = ac.id
	  )t

	inner join account ac on t.to_count_id = ac.id

	where t.from_count_id = %d and t.type = %d`, fromCountId, message.FriendsAddAccept)
	//MyLOG.Log(msql)
	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		var faa message.FriendsAddAcceptData

		err = rows.Scan(&faa.FromClientId, &faa.FromClientName, &faa.ToClientId, &faa.ToClientName)
		if err != nil {
			MyLOG.ErrLog("Scan error: %v", err)
			causeId = message.ID_SERVER_ERROR
			return
		}
		sliceFAA = append(sliceFAA, &faa)
	}
	//sort.Sort(sliceMsg)
	MyLOG.Log("FriendsBoxSelect: %v", sliceFAA)
	return
}

func (thiss *ChatroomDB) FriendsAddAcceptDelete(fromCountId uint16) {

	msql := fmt.Sprintf("delete from friend_box where from_count_id = ? and type = ?")
	stmt, err := thiss.DB.Prepare(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Prepare fail: %v", err)
		return
	}
	_, err = stmt.Exec(fromCountId, message.FriendsAddAccept)
	return
}

func (thiss *ChatroomDB) FriendsAddRejectSelect(fromCountId uint16) (sliceFARJ []*message.FriendsAddRejectData, causeId message.CauseId) {

	//fromCountId: who trigger the FriendsAddRequest

	sliceFARJ = make([]*message.FriendsAddRejectData, 0, 4096)
	msql := fmt.Sprintf(`
	select t.from_client_id, t.from_client_name, ac.client_id as to_client_id, ac.count_name as to_client_name from (
		
		select f.*, ac.client_id as from_client_id, ac.count_name as from_client_name
		from friend_box f 
		inner join account ac on f.from_count_id = ac.id
	  )t

	inner join account ac on t.to_count_id = ac.id

	where t.from_count_id = %d and t.type = %d`, fromCountId, message.FriendsAddReject)
	//MyLOG.Log(msql)
	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		var faa message.FriendsAddRejectData

		err = rows.Scan(&faa.FromClientId, &faa.FromClientName, &faa.ToClientId, &faa.ToClientName)
		if err != nil {
			MyLOG.ErrLog("Scan error: %v", err)
			causeId = message.ID_SERVER_ERROR
			return
		}
		sliceFARJ = append(sliceFARJ, &faa)
	}
	//sort.Sort(sliceMsg)
	MyLOG.Log("FriendsBoxSelect: %v", sliceFARJ)
	return
}

func (thiss *ChatroomDB) FriendsAddRejectDelete(fromCountId uint16) {

	msql := fmt.Sprintf("delete from friend_box where from_count_id = ? and type = ?")
	stmt, err := thiss.DB.Prepare(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Prepare fail: %v", err)
		return
	}
	_, err = stmt.Exec(fromCountId, message.FriendsAddReject)
	return
}

func (thiss *ChatroomDB) FriendsBoxDelete(fromClientId string, toClientId string, msg string, msgType uint8) (causeId message.CauseId) {

	fromCountId, causeId := thiss.GetCountId(fromClientId)
	if causeId != 0 {
		MyLOG.Log("GetCountId wish clientId \"%s\" fail", fromClientId)
		return
	}

	toCountId, causeId := thiss.GetCountId(toClientId)
	if causeId != 0 {
		MyLOG.Log("GetCountId wish clientId \"%s\" fail", toClientId)
		return
	}

	msql := fmt.Sprintf("delete from friend_box where from_count_id = ? and to_count_id = ? and message = ? and type = ?")
	MyLOG.Log(msql)
	stmt, err := thiss.DB.Prepare(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Prepare fail: %v", err)
		return
	}
	_, err = stmt.Exec(fromCountId, toCountId, msg, msgType)
	return
}
