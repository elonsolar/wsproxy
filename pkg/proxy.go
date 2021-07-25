package pkg

import (
	"fmt"
	"io"
	"net"

	"github.com/gorilla/websocket"
)

// a proxy manage  connection  from client  port <--cross websocket---> server port
type TcpProxy struct {
	Name          string
	ClientPort    int
	RemotePort    int
	RemoteAddress string
}

// 服务端建立内部连接
func (tp *TcpProxy) Serve(con net.Conn, wsCon *websocket.Conn)  {

	var copyWs2Tcp = func(con net.Conn, wsCon *websocket.Conn) {
		fmt.Println("ws ---->tcp")
		for {
			_, buff, err := wsCon.ReadMessage()
			if err != nil {
				wsCon.Close()
			}
			fmt.Println(" read from ws", string(buff))
			_, err = con.Write(buff)
			if err != nil {
				con.Close()
			}

		}
	}

	var copyTcp2Ws = func(wsCon *websocket.Conn, con net.Conn) {
		fmt.Println("tcp ---->ws")
		buff := make([]byte, 16*1024)
		for {
			n, err := con.Read(buff)
			if n > 0 {
				err = wsCon.WriteMessage(websocket.BinaryMessage, buff[:n])
				if err != nil {
					wsCon.Close()
				}
			}
			if err != nil {
				if err != io.EOF {
					con.Close()
					return
				}
			}
		}

	}
	go copyWs2Tcp(con, wsCon)
	go copyTcp2Ws(wsCon, con)
}
