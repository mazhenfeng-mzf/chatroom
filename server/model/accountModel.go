package model

import (
	"chatroom/public/message"
	. "chatroom/public/tools"

	//. "chatroom/server/config"
	//"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Count struct {
	Id        uint16 `db:"id"`
	CountId   string `db:"client_id"`
	CountName string `db:"count_name"`
	CountPwd  string `db:"count_pwd"`
}

func (thiss *ChatroomDB) CountExist(client_id string) (yes bool, causeId message.CauseId) {

	msql := fmt.Sprintf("select * from account where client_id = %s", client_id)
	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	} else {
		defer rows.Close()
	}

	isexist := rows.Next()
	if !isexist {
		yes = false
		return
	}

	yes = true
	return
}

func (thiss *ChatroomDB) CountLogin(client_id string, count_pwd string, count *Count) (causeId message.CauseId) {

	msql := fmt.Sprintf("select * from account where client_id = %s and count_pwd = %s", client_id, count_pwd)
	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	} else {
		defer rows.Close()
	}

	ok := rows.Next()
	if !ok {
		MyLOG.ErrLog("rows Next fail: %v", err)
		causeId = message.ID_COUNT_WRONG_PASSWORD
		return
	}

	err = rows.Scan(&count.Id, &count.CountId, &count.CountName, &count.CountPwd)
	if err != nil {
		MyLOG.ErrLog("Scan error: %v", err)
		causeId = message.ID_SERVER_ERROR
	}
	MyLOG.Log("count: %v", count)
	return

}

func (thiss *ChatroomDB) CountRegister(client_id string, count_name string, count_pwd string) (causeId message.CauseId) {

	stmt, err := thiss.DB.Prepare(`INSERT account (client_id,count_name,count_pwd) values (?,?,?)`)
	if err != nil {
		MyLOG.ErrLog("Mysql Prepare fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	_, err = stmt.Exec(client_id, count_name, count_pwd)

	return

}

func (thiss *ChatroomDB) CountSearch(client_id string) (count *Count, causeId message.CauseId) {

	msql := fmt.Sprintf("select id, client_id, count_name from account where client_id = %s", client_id)
	//MyLOG.Log("msql: %v", msql)
	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	} else {
		defer rows.Close()
	}

	ok := rows.Next()
	if !ok {
		causeId = message.ID_COUNT_NOT_EXIST
		return
	}
	count = new(Count)
	err = rows.Scan(&count.Id, &count.CountId, &count.CountName)
	if err != nil {
		MyLOG.ErrLog("Scan error: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	MyLOG.Log("Search count: %v", count)
	return

}

func (thiss *ChatroomDB) GetCountId(client_id string) (Id uint16, causeId message.CauseId) {

	msql := fmt.Sprintf("select id from account where client_id = %s", client_id)

	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	} else {
		defer rows.Close()
	}

	ok := rows.Next()
	if !ok {
		MyLOG.ErrLog("用户 %s 不存在", client_id)
		causeId = message.ID_COUNT_NOT_EXIST
		return
	}

	err = rows.Scan(&Id)
	if err != nil {
		MyLOG.ErrLog("Scan error: %v", err)
		causeId = message.ID_SERVER_ERROR
	}
	return
}
