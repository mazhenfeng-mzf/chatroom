package processpkg

import (
	"chatroom/public/errs"
	"chatroom/public/message"
	. "chatroom/public/tools"
	"encoding/json"
	"net"
)

func WritePkg(conn net.Conn, msg *message.MessageDataW) (err error) {

	//0. 将 消息Data 序列化, 获取长度

	//MyLOG.Log("消息数据格式化前: %#v", msg.Data)
	data_byte, err := json.Marshal(msg.Data)
	if err != nil {
		MyLOG.Log("json格式化失败,message: %v \n err: %v", msg.Data, err)
		return
	}
	//MyLOG.Log("消息数据格式化后: %#v", data_byte)

	length := len(data_byte)

	//1 构建结构体
	msg_sock := message.MessageSocket{
		Type:   msg.Type,
		Length: uint32(length),
		Data:   data_byte,
	}
	MyLOG.Log("发送消息类型是 %d, byte长度是 %d", msg.Type, length)
	//2 结构体序列化
	msg_sock_byte, err := json.Marshal(msg_sock)
	if err != nil {
		MyLOG.Log("json格式化失败,message: %v \n err: %v", msg_sock, err)
		return
	}
	//3 发送
	_, err = conn.Write(msg_sock_byte)
	if err != nil {
		MyLOG.Log("发送消息失败: %v \n err: %v", msg_sock, err)
		return
	}
	//MyLOG.Log("发送消息成功: , %#v(%T)", msg_sock, msg_sock)
	//MyLOG.Log("发送消息byte成功: , %v", msg_sock_byte)
	return
}

func ReadPkg(conn net.Conn, msg *message.MessageDataR) (err error) {

	//1 获取全部信息
	buf := make([]byte, 65535)
	n, err := conn.Read(buf)
	if err != nil {
		MyLOG.Log("1. 读取消息失败, err: %v", err)
		return
	}
	//MyLOG.Log("1. 读取数据成功: %v", buf[:n])
	//1.1  反序列化
	var msg_sock message.MessageSocket
	err = json.Unmarshal(buf[:n], &msg_sock)
	if err != nil {
		MyLOG.Log("1. 反序列化失败, err: %v", err)
		return
	}
	//MyLOG.Log("1.1. 反序列化后: %#v", msg_sock)

	//1.2 获取Type 和 Length
	msg_type := msg_sock.Type
	length := msg_sock.Length

	MyLOG.Log("读取消息类型是 %d, byte长度是 %d", msg_type, length)

	//2 根据lenght截取Data
	if length > uint32(len(msg_sock.Data)) {
		MyLOG.Log("Data 长度不足")
		err = errs.INVAILD_MSG_INSUFFICIENT_LENGTH
		return
	}
	all_data_byte := msg_sock.Data
	data_byte := all_data_byte[:length]

	// 给到msg
	msg.Type = msg_type
	msg.Data = data_byte

	//MyLOG.Log("3. 总: 读取消息成功: (Type: %d) %v", msg.Type, msg.Data)
	return
}
