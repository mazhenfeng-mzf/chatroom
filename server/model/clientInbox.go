package model

import (
	//_"github.com/jmoiron/sqlx"

	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"
	"sort"

	_ "github.com/go-sql-driver/mysql"
)

type ClientInbox struct {
	Id          string `db:"id"`
	FromCountId string `db:"from_count_id"`
	//FromClientName string `db:"count_name"`
	ToCountId string `db:"to_count_id"`
	Data      string `db:"data"`
	Time      string `db:"time"`

	FromClientId   string
	FromClientName string
	ToClientId     string
	ToClientName   string
}

type SliceClientInbox []*ClientInbox

//for sort.Sort
func (hs SliceClientInbox) Len() int {
	return len(hs)
}
func (hs SliceClientInbox) Less(i, j int) bool {
	return hs[i].Time < hs[j].Time
}
func (hs SliceClientInbox) Swap(i, j int) {
	hs[i], hs[j] = hs[j], hs[i]
}

type SliceMsgInbox []*message.SmsP2PData

//for sort.Sort
func (hs SliceMsgInbox) Len() int {
	return len(hs)
}
func (hs SliceMsgInbox) Less(i, j int) bool {
	return hs[i].Time < hs[j].Time
}
func (hs SliceMsgInbox) Swap(i, j int) {
	hs[i], hs[j] = hs[j], hs[i]
}

func (thiss *ChatroomDB) ClientInboxInsert(fromClientId string, fromCountId uint16, toClientId string, data string, timeString string) (causeId message.CauseId) {

	if fromCountId == 0 {
		fromCountId, causeId = thiss.GetCountId(fromClientId)
		if causeId != 0 {
			MyLOG.ErrLog("GetCountId wish clientId \"%s\" fail", fromClientId)
			return
		}
	}

	toCountId, causeId := thiss.GetCountId(toClientId)
	MyLOG.Log("GetCountId - causeId: %s", causeId)
	if causeId != 0 {
		MyLOG.ErrLog("GetCountId wish clientId \"%s\" fail", toClientId)
		return
	}

	stmt, err := thiss.DB.Prepare(`INSERT clientInbox (to_count_id,from_count_id,data,time) values (?,?,?,?)`)
	if err != nil {
		MyLOG.ErrLog("Mysql Prepare fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	_, err = stmt.Exec(toCountId, fromCountId, data, timeString)
	if err != nil {
		MyLOG.ErrLog("Mysql Exec fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	return
}

func (thiss *ChatroomDB) ClientInboxSelect(toClientId string, toCountId uint16) (sliceMsg SliceMsgInbox, causeId message.CauseId) {

	if toCountId == 0 {
		toCountId, causeId = thiss.GetCountId(toClientId)
		if causeId != 0 {
			MyLOG.ErrLog("GetCountId wish clientId \"%s\" fail", toClientId)
			return
		}
	}

	sliceMsg = make(SliceMsgInbox, 0, 4096)
	msql := fmt.Sprintf(`
	select c.client_id as from_client_id, c.count_name as from_count_name, t.to_client_id, t.data, t.time  from (
	
		select a.*,b.client_id as to_client_id, b.count_name as to_count_name 
		from clientInbox a 
		inner join account b on a.to_count_id = b.id
		)t
	
	inner join account c on t.from_count_id = c.id
	
	where to_count_id = %d`, toCountId)
	//MyLOG.Log("ClientInboxSelect - msql: %v", msql)
	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		var msgP2p message.SmsP2PData

		err = rows.Scan(&msgP2p.FromClientId, &msgP2p.FromClientName, &msgP2p.ToClientId, &msgP2p.Data, &msgP2p.Time)
		if err != nil {
			MyLOG.ErrLog("Scan error: %v", err)
			causeId = message.ID_SERVER_ERROR
			return
		}
		sliceMsg = append(sliceMsg, &msgP2p)
	}
	sort.Sort(sliceMsg)
	MyLOG.Log("Search ClientInbox: %v", sliceMsg)
	return
}

func (thiss *ChatroomDB) ClientInboxDelete(fromClientId string, toClientId string) (causeId message.CauseId) {

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

	//MyLOG.Log("Delete from clientInbox where to_count_id = %d", toCountId)
	stmt, err := thiss.DB.Prepare(`Delete from clientInbox where from_count_id = ? and to_count_id = ?`)
	if err != nil {
		MyLOG.ErrLog("Mysql Prepare fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	_, err = stmt.Exec(fromCountId, toCountId)

	return
}
