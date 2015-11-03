package tcpserver

type Handler interface {
	HandlerMessage(conn *Conn, data []byte)

	Connect(*Conn)

	Disconnect(*Conn)
}
