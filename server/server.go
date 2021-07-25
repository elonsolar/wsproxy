package main

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
	. "github.com/elonsolar/wsproxy/pkg"
	"github.com/gorilla/websocket"
)

// 每个control 管理多个连接
// 每个客户端 可以有多个 连接
//
type Control struct {
	RunId  string
	conMap map[string]*websocket.Conn
}

//   make it easy
func WsProxy(writer http.ResponseWriter, request *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}
	go handleConnection(conn)
}

// first of all login
func handleConnection(con *websocket.Conn) {

	var pxy TcpProxy = TcpProxy{
		Name:          "ssh",
		ClientPort:    9121,
		RemotePort:    22,
		RemoteAddress: "dev.codenai.com",
	}
	err := con.ReadJSON(&pxy)
	if err != nil {
		Zapper.Error("login err", zap.Error(err))
		return
	}
	fmt.Println(pxy)
	localCon, err := net.Dial("tcp", fmt.Sprintf("%s:%d", pxy.RemoteAddress, pxy.RemotePort))
	if err != nil {
		panic(err)
	}
	

	pxy.Serve(localCon, con)
}
