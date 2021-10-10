package model

import (
	. "chatroom/public/tools"
	. "chatroom/server/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type ChatroomDB struct {
	DB *sql.DB
}

var CrDB ChatroomDB

const (
	DB_INIT_account string = `
CREATE TABLE account (
	id int(11) NOT NULL AUTO_INCREMENT,
	client_id varchar(11) NOT NULL,
	count_name varchar(255) DEFAULT NULL,
	count_pwd varchar(255) DEFAULT NULL,
	PRIMARY KEY (id,client_id) USING BTREE,
	KEY id (id)
  ) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
`
	DB_INIT_clientinbox string = `
CREATE TABLE clientinbox (
	id int(255) NOT NULL AUTO_INCREMENT,
	to_count_id int(255) NOT NULL,
	from_count_id int(255) NOT NULL,
	data varchar(255) DEFAULT NULL,
	time datetime(6) DEFAULT NULL,
	PRIMARY KEY (id) USING BTREE,
	KEY FK_ID_2 (to_count_id),
	KEY index_time (time),
	CONSTRAINT FK_ID_1 FOREIGN KEY (to_count_id) REFERENCES account (id),
	CONSTRAINT FK_ID_2 FOREIGN KEY (to_count_id) REFERENCES account (id)
  ) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
`
	DB_INIT_friend_box string = `
  CREATE TABLE friend_box (
	id int(11) NOT NULL AUTO_INCREMENT,
	from_count_id int(255) DEFAULT NULL,
	to_count_id int(255) DEFAULT NULL,
	message varchar(10000) DEFAULT NULL,
	type int(255) DEFAULT NULL,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8;
`
	DB_INIT_friends string = `
  CREATE TABLE friends (
	id int(255) NOT NULL AUTO_INCREMENT,
	count_id int(255) NOT NULL,
	friend_count_id int(255) DEFAULT NULL,
	time varchar(10000) DEFAULT NULL,
	PRIMARY KEY (id) USING BTREE,
	KEY FK_ID_friends_1 (count_id),
	KEY FK_ID_friends_2 (friend_count_id),
	CONSTRAINT FK_ID_friends_1 FOREIGN KEY (count_id) REFERENCES account (id),
	CONSTRAINT FK_ID_friends_2 FOREIGN KEY (friend_count_id) REFERENCES account (id)
  ) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
`
	DB_INIT_sms_history string = `
  CREATE TABLE sms_history (
	id int(255) NOT NULL AUTO_INCREMENT,
	from_count_id int(255) NOT NULL,
	to_count_id int(255) NOT NULL,
	data varchar(255) DEFAULT NULL,
	time datetime(6) DEFAULT NULL,
	PRIMARY KEY (id) USING BTREE,
	KEY index_sms_time (time),
	KEY FK_ID_history_1 (from_count_id),
	KEY FK_ID_history_2 (to_count_id),
	CONSTRAINT FK_ID_history_1 FOREIGN KEY (from_count_id) REFERENCES account (id),
	CONSTRAINT FK_ID_history_2 FOREIGN KEY (to_count_id) REFERENCES account (id)
  ) ENGINE=InnoDB AUTO_INCREMENT=73 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
`
)

func (thiss *ChatroomDB) Connect() (err error) {

	db, err := sql.Open("mysql", MyConfig.DB_CONNECT_STRING)
	if err != nil {
		MyLOG.ErrLog("连接数据库失败")
		return
	}

	db.SetMaxOpenConns(MyConfig.DB_MAX_OPEN_CONNS)
	db.SetMaxIdleConns(MyConfig.DB_MAX_IDLE_CONNS)
	thiss.DB = db
	return
}

func (thiss *ChatroomDB) Init() (err error) {

	for _, db_init := range [...]string{DB_INIT_account, DB_INIT_clientinbox, DB_INIT_friend_box, DB_INIT_friends, DB_INIT_sms_history} {
		_, err = thiss.DB.Exec(db_init)
		if err != nil {
			MyLOG.ErrLog("Mysql Exec fail: %v", err)
			return
		}
	}
	return
}
