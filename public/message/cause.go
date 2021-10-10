package message

type CauseId uint32

var CauseMap map[CauseId]string = map[CauseId]string{
	ID_INVAILD_MSG_INSUFFICIENT_LENGTH: INVAILD_MSG_INSUFFICIENT_LENGTH,
	ID_INVAILD_MSG_PARSE_FAIL:          INVAILD_MSG_PARSE_FAIL,
	ID_COUNT_NOT_EXIST:                 COUNT_NOT_EXIST,
	ID_COUNT_WRONG_PASSWORD:            COUNT_WRONG_PASSWORD,
	ID_SERVER_ERROR:                    SERVER_ERROR,
	ID_COUNT_EXIST:                     COUNT_EXIST,
	ID_COUNT_ID_LENGTH_WRONG:           COUNT_ID_LENGTH_WRONG,
}

var (
	ID_INVAILD_MSG_INSUFFICIENT_LENGTH CauseId = 1
	INVAILD_MSG_INSUFFICIENT_LENGTH    string  = "无效的信息长度"

	ID_INVAILD_MSG_PARSE_FAIL CauseId = 2
	INVAILD_MSG_PARSE_FAIL    string  = "信息解析失败"

	ID_COUNT_NOT_EXIST CauseId = 3
	COUNT_NOT_EXIST    string  = "用户不存在"

	ID_COUNT_WRONG_PASSWORD CauseId = 4
	COUNT_WRONG_PASSWORD    string  = "用户密码错误"

	ID_SERVER_ERROR CauseId = 5
	SERVER_ERROR    string  = "服务器错误"

	ID_COUNT_EXIST CauseId = 6
	COUNT_EXIST    string  = "用户已存在"

	ID_COUNT_ID_LENGTH_WRONG CauseId = 7
	COUNT_ID_LENGTH_WRONG    string  = "用户 ID 长度错误"
)
