package tcp

//CallBack 是一个回调接口，用于连接的各种事件处理
type CallBack interface {
	//链接建立回调
	OnConnected(conn *TCPConn)
	//消息处理回调
	OnMessage(conn *TCPConn, p Packet)
	//链接断开回调
	OnDisconnected(conn *TCPConn)
	//错误回调
	OnError(err error)
}
