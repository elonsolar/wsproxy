package pkg

import (
	"fmt"
	"go.uber.org/zap"
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
		defer con.Close()
		defer wsCon.Close()
		fmt.Println("ws ---->tcp")
		for {
			_, buff, err := wsCon.ReadMessage()
			if err != nil {
				Logger.Error("read from websocket",zap.Error(err))
				return
			}
			if len(buff)>0{
				fmt.Println(" read from ws", string(buff))
				_, err = con.Write(buff)
				if err != nil {
					Logger.Error("write to tcp connection",zap.Error(err))
					return
				}
			}

		}
	}

	var copyTcp2Ws = func(wsCon *websocket.Conn, con net.Conn) {
		defer wsCon.Close()
		defer con.Close()
		fmt.Println("tcp ---->ws")
		buff := make([]byte, 16*1024)
		for {
			n, err := con.Read(buff)
			if n > 0 {
				err = wsCon.WriteMessage(websocket.BinaryMessage, buff[:n])
				if err != nil {
					Logger.Error("write to websocket",zap.Error(err))
					return
				}
			}
			if err != nil {
				if err != io.EOF {
					Logger.Error("read from tcp connection",zap.Error(err))
					return
				}
			}
		}

	}
	go copyWs2Tcp(con, wsCon)
	go copyTcp2Ws(wsCon, con)
}
