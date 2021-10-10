package model

import (
	//_"github.com/jmoiron/sqlx"

	"chatroom/public/message"
	. "chatroom/public/tools"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Friends struct {
	Id            int    `db:"id"`
	CountId       int    `db:"count_id"`
	FriendCountId int    `db:"friend_count_id"`
	Time          string `db:"time"`
}

func (thiss *ChatroomDB) FriendsAdd(CountId uint16, friendClientId string, timeString string) (causeId message.CauseId) {
	friendCountId, causeId := thiss.GetCountId(friendClientId)
	if causeId != 0 {
		MyLOG.ErrLog("GetCountId wish clientId \"%s\" fail", friendCountId)
		return
	}

	if yes, causeId := thiss.FriendsIsExist(CountId, friendCountId); !yes {
		MyLOG.Log("FriendsIsExist(%s, %s): %v", CountId, friendCountId, yes)
		if causeId != 0 {
			MyLOG.ErrLog("Check FriendsIsExist err")
			return causeId
		}
		causeId = thiss.FriendsInsert(CountId, friendCountId, timeString)
		if causeId != 0 {
			MyLOG.ErrLog("Check FriendsInsert err")
			return causeId
		}
	}

	if yes, causeId := thiss.FriendsIsExist(friendCountId, CountId); !yes {
		if causeId != 0 {
			MyLOG.ErrLog("Check FriendsIsExist err")
			return causeId
		}
		causeId = thiss.FriendsInsert(friendCountId, CountId, timeString)
		if causeId != 0 {
			MyLOG.ErrLog("Check FriendsInsert err")
			return causeId
		}
	}

	return
}

func (thiss *ChatroomDB) FriendsInsert(CountId uint16, friendCountId uint16, timeString string) (causeId message.CauseId) {

	stmt, err := thiss.DB.Prepare(`INSERT friends (count_id,friend_count_id,time) values (?,?,?)`)
	if err != nil {
		MyLOG.ErrLog("Mysql Prepare fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	_, err = stmt.Exec(CountId, friendCountId, timeString)
	if err != nil {
		MyLOG.ErrLog("Mysql Exec fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	}
	return
}

func (thiss *ChatroomDB) FriendsSelect(CountId uint16) (sliceFriends []*message.ClientUserData, causeId message.CauseId) {

	sliceFriends = make([]*message.ClientUserData, 0, 4096)
	msql := fmt.Sprintf(`
	select friend_client_id, friend_name  from (
	
		select f.*,ac.client_id as friend_client_id, ac.count_name as friend_name 
		from friends f 
		inner join account ac on ac.id = f.friend_count_id
		)t
	
	where count_id = %d`, CountId)

	rows, err := thiss.DB.Query(msql)
	if err != nil {
		MyLOG.ErrLog("Mysql Query fail: %v", err)
		causeId = message.ID_SERVER_ERROR
		return
	} else {
		defer rows.Close()
	}

	for rows.Next() {
		var friend message.ClientUserData

		err = rows.Scan(&friend.Id, &friend.Name)
		if err != nil {
			MyLOG.ErrLog("Scan error: %v", err)
			causeId = message.ID_SERVER_ERROR
			return
		}
		sliceFriends = append(sliceFriends, &friend)
	}
	//sort.Sort(sliceMsg)
	MyLOG.Log("FriendsSelect: %v", sliceFriends)
	return
}

func (thiss *ChatroomDB) FriendsIsExist(countId uint16, friendCountId uint16) (yes bool, causeId message.CauseId) {

	msql := fmt.Sprintf("select id from friends where count_id = %d and friend_count_id = %d", countId, friendCountId)
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
