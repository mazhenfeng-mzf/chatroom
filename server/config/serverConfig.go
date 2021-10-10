package config

import (
	. "chatroom/public/tools"
	"fmt"

	"gopkg.in/ini.v1"
)

type Config struct { //配置文件要通过tag来指定配置文件中的名称
	// db_user = root
	// db_pwd = mzf123
	// db_protocol = tcp
	// db_ip = localhost
	// db_port = 3306
	// db_name = chatroom

	DB_USER           string `ini:"db_user"`
	DB_PWD            string `ini:"db_pwd"`
	DB_PROTOCOL       string `ini:"db_protocol"`
	DB_IP             string `ini:"db_ip"`
	DB_PORT           string `ini:"db_port"`
	DB_NAME           string `ini:"db_name"`
	DB_CONNECT_STRING string

	LISTEN_PORT string `ini:"listen_port"`
}

var MyConfig *Config

// func main() {
// 	config, err := ReadConfig(filepath) //也可以通过os.arg或flag从命令行指定配置文件路径
// 	if err != nil {
// 		MyLOG.ErrLog("ReadConfig \"%s\" error: %v", filepath, err)
// 	}
// }

//读取配置文件并转成结构体
func ReadConfig(path string) (err error) {
	//var config Config
	MyConfig = new(Config)

	conf, err := ini.Load(path) //加载配置文件
	if err != nil {
		MyLOG.ErrLog("load config file \"%s\" error: %v", path, err)
		return
	}
	conf.BlockMode = false
	err = conf.MapTo(MyConfig) //解析成结构体
	if err != nil {
		MyLOG.ErrLog("mapto config file \"%s\" error: %v", path, err)
		return
	}
	//MyLOG.Log(MyConfig.DB_USER)
	//"root:mzf123@tcp(localhost:3306)/chatroom"
	MyConfig.DB_CONNECT_STRING = fmt.Sprintf("%s:%s@%s(%s:%s)/%s", MyConfig.DB_USER, MyConfig.DB_PWD, MyConfig.DB_PROTOCOL, MyConfig.DB_IP, MyConfig.DB_PORT, MyConfig.DB_NAME)
	//MyLOG.Log(MyConfig.DB_CONNECT_STRING)
	return
}
