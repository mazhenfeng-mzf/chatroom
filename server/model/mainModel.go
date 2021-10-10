package model

import "database/sql"

// const (
// 	connectDB string = "root:mzf123@tcp(localhost:3306)/chatroom"
// )

type ChatroomDB struct {
	DB *sql.DB
}

var CrDB ChatroomDB
