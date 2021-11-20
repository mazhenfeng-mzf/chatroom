

**chatroom**聊天室
 
- **基于golang socket 编程** 
- **基于window GUI 界面 walk 模块** 
- **使用外部数据库 mysql** 

-------------------

# 示范
![chatroom](https://github.com/mazhenfeng-mzf/chatroom/blob/master/chatroom_01.mp4)


# 使用说明
一个服务器, 多客户端模式

## 服务器

### 准备一个 mysql 数据库 
需要手动创建一个数据库
比如

![image](https://user-images.githubusercontent.com/63535556/136687762-530a4a0d-3f12-4dd9-a3e5-9bf2bdf77a21.png)

### 修改配置文件 chatroomServer.conf
```
# mysql db 连接设置
db_user = root
db_pwd = mzf123
db_protocol = tcp
db_ip = localhost
db_port = 3306
db_name = chatroom

# 服务器监听端口
listen_port = 8888  
```

### 编译服务器
自行设置 GOPATH
```
MZF@DESKTOP-HDOQO35 MINGW64 ~/go/src/chatroom (master)
$ cd ./main/

MZF@DESKTOP-HDOQO35 MINGW64 ~/go/src/chatroom/main (master)
$ pwd
/c/Users/MZF/go/src/chatroom/main

MZF@DESKTOP-HDOQO35 MINGW64 ~/go/src/chatroom/main (master)
$ go env |grep GOPATH
C:\Users\MZF\go;

```

```
MZF@DESKTOP-HDOQO35 MINGW64 ~/go/src/chatroom/main (master)
$ go build  -o serverDbInit.exe serverDB_INIT.go

```

```
MZF@DESKTOP-HDOQO35 MINGW64 ~/go/src/chatroom/main (master)
$ go build -ldflags="-H windowsgui" -o server.exe serverMain.go
```

### 初始化数据库
```
MZF@DESKTOP-HDOQO35 MINGW64 ~/go/src/chatroom/main (master)
$ ./serverDbInit.exe
读取配置文件 ./chatroomServer.conf 成功: &config.Config{DB_USER:"root", DB_PWD:"mzf123", DB_PROTOCOL:"tcp", DB_IP:"localhost", DB_PORT:"3306", DB_NAME:"test", DB_CONNECT_STRING:"root:mzf123@tcp(localhost:3306)/test", LISTEN_PORT:"8888", DB_MAX_OPEN_CONNS:10, DB_MAX_IDLE_CONNS:5}
连接 DB 成功
初始化 DB 成功
```

### 运行server
```
MZF@DESKTOP-HDOQO35 MINGW64 ~/go/src/chatroom/main (master)
$ ./server.exe
日志文件不存在
打开日志文件 C://Users/MZF/go/src/chatroom/main/logServer/clientLog-2021-10-10-15-57-46.txt
<info>-<serverMain.go><23><func main>: 读取配置文件 ./chatroomServer.conf 成功: &config.Config{DB_USER:"root", DB_PWD:"mzf123", DB_PROTOCOL:"tcp", DB_IP:"localhost", DB_PORT:"3306", DB_NAME:"test", DB_CONNECT_STRING:"root:mzf123@tcp(localhost:3306)/test", LISTEN_PORT:"8888", DB_MAX_OPEN_CONNS:10, DB_MAX_IDLE_CONNS:5}

server start listening
waitting client1

```

## 客户端


### 修改配置文件 chatroomClient.conf
```
# 服务器地址和端口
server_ip = 192.168.1.103
server_port = 8888
```

### 编译客户端
```
MZF@DESKTOP-HDOQO35 MINGW64 ~/go/src/chatroom/main (master)
$ go build -ldflags="-H windowsgui" -o client.exe clientMain.go
```

### 运行客户端
```
MZF@DESKTOP-HDOQO35 MINGW64 ~/go/src/chatroom/main (master)
$ ./client.exe
```


![image](https://user-images.githubusercontent.com/63535556/136687778-38327907-eba2-4aba-96b5-6878103b6936.png)
