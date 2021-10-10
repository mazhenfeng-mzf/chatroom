package model

import (
	//_"github.com/jmoiron/sqlx"

	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"
	"sort"

	_ "github.com/go-sql-driver/mysql"
)

type SmsHistory struct {
	Id             string `db:"id"`
	FromClientId   string `db:"from_client_id"`
	FromClientName string `db:"count_name"`
	ToClientId     string `db:"to_client_id"`
	Data           string `db:"data"`
	Time           string `db:"time"`
}

type SliceSmsHistory []*message.SmsP2PData

//for sort.Sort
func (hs SliceSmsHistory) Len() int {
	return len(hs)
}
func (hs SliceSmsHistory) Less(i, j int) bool {
	return hs[i].Time < hs[j].Time
}
func (hs SliceSmsHistory) Swap(i, j int) {
	hs[i], hs[j] = hs[j], hs[i]
}

func (thiss *ChatroomDB) SmsHistoryAdd(fromClientId string, toClientId string, data string, timeString string) (causeId message.CauseId) {
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

	causeId = thiss.SmsHistoryInsert(fromCountId, toCountId, data, timeString)
	if causeId != 0 {
		MyLOG.Log("SmsHistoryInsert fail")
		return
	}
	return
}

func (thiss *ChatroomDB) SmsHistoryInsert(fromCountId uint16, toCountId uint16, data string, timeString string) (causeId message.CauseId) {
	stmt, err := thiss.DB.Prepare(`INSERT sms_history (from_count_id,to_count_id,data,time) values (?,?,?,?)`)
	if err != nil {
		MyLOG.ErrLog("Mysql Prepare fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	_, err = stmt.Exec(fromCountId, toCountId, data, timeString)
	if err != nil {
		MyLOG.ErrLog("Mysql Exec fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	return
}

func (thiss *ChatroomDB) SmsHistorySelect(ClientId1 string, ClientId2 string, number int, offset uint32) (smsHistorySlice SliceSmsHistory) {
	smsHistorySlice = make(SliceSmsHistory, 0, 4096)

	CountId1, causeId := thiss.GetCountId(ClientId1)
	if causeId != 0 {
		MyLOG.Log("GetCountId wish clientId \"%s\" fail", ClientId1)
		return
	}

	CountId2, causeId := thiss.GetCountId(ClientId2)
	if causeId != 0 {
		MyLOG.Log("GetCountId wish clientId \"%s\" fail", ClientId2)
		return
	}

	msql := fmt.Sprintf(`
	select c.client_id as from_client_id, c.count_name as from_count_name, t.to_client_id, t.to_count_name, t.data, t.time  from (
	
		select a.*,b.client_id as to_client_id, b.count_name as to_count_name 
		from sms_history a 
		inner join account b on a.to_count_id = b.id 
		
		where (from_count_id = %d and to_count_id = %d) or (from_count_id = %d and to_count_id = %d)
		
		)t
	
	inner join account c on t.from_count_id = c.id

	ORDER BY t.time desc limit %d, %d`, CountId1, CountId2, CountId2, CountId1, offset, number)

	MyLOG.Log("msql: %v", msql)

	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		return
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		var sms message.SmsP2PData
		err = rows.Scan(&sms.FromClientId, &sms.FromClientName, &sms.ToClientId, &sms.ToClientName, &sms.Data, &sms.Time)
		if err != nil {
			MyLOG.ErrLog("Scan error: %v", err)
			return
		}
		smsHistorySlice = append(smsHistorySlice, &sms)
	}
	sort.Sort(smsHistorySlice)
	MyLOG.Log("Search smsHistorySlice: %v", smsHistorySlice)
	return
}

func (thiss *ChatroomDB) SmsHistorySum(ClientId1 string, ClientId2 string) (sum uint32) {

	CountId1, causeId := thiss.GetCountId(ClientId1)
	if causeId != 0 {
		MyLOG.Log("GetCountId wish clientId \"%s\" fail", ClientId1)
		return
	}

	CountId2, causeId := thiss.GetCountId(ClientId2)
	if causeId != 0 {
		MyLOG.Log("GetCountId wish clientId \"%s\" fail", ClientId2)
		return
	}

	sum = 0
	msql := fmt.Sprintf(`select Count(id)
	from sms_history 
	where (from_count_id = %d and to_count_id = %d) or (from_count_id = %d and to_count_id = %d)
	`, CountId1, CountId2, CountId2, CountId1)

	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		return
	} else {
		defer rows.Close()
	}

	rows.Next()
	rows.Scan(&sum)
	MyLOG.Log("SmsHistorySum: %v", sum)
	return
}
