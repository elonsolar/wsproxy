package client

import (
	"fmt"
	"net"

	. "github.com/elonsolar/wsproxy/pkg"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Config struct {
	ServerAddr string
	Proxy      []*TcpProxy
}

func Proxy(pxyCfg *Config) {
	for _, pxy := range pxyCfg.Proxy {

		tcpaddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("127.0.0.1:%d", pxy.ClientPort))

		if err != nil {
			panic(err)
		}
		listenner, err := net.ListenTCP("tcp", tcpaddr)

		if err != nil {
			panic(err)
		}
		go handleListener(pxyCfg.ServerAddr, pxy, listenner)
	}

}

func handleListener(serverAddr string, pxy *TcpProxy, l *net.TCPListener) {
	for {
		con, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(serverAddr, pxy, con)
	}
}

func handleConnection(serverAddr string, pxy *TcpProxy, con net.Conn) {

	wsCon, err := getWsConnection(serverAddr, pxy)
	if err != nil {
		Logger.Error("ws error", zap.Error(err))
		return
	}
	pxy.Serve(con, wsCon)
}

func getWsConnection(serverAddr string, pxy *TcpProxy) (*websocket.Conn, error) {

	c, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("dial websocket err:%w", err)
	}

	err = c.WriteJSON(pxy)
	if err != nil {
		return nil, fmt.Errorf("init websocket connection err:%w", err)
	}

	return c, nil
}
