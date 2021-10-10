package config

import (
	. "chatroom/public/tools"

	"gopkg.in/ini.v1"
)

type Config struct { //配置文件要通过tag来指定配置文件中的名称
	// db_user = root
	// db_pwd = mzf123
	// db_protocol = tcp
	// db_ip = localhost
	// db_port = 3306
	// db_name = chatroom

	SERVER_IP   string `ini:"server_ip"`
	SERVER_PORT string `ini:"server_port"`
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
	return
}
