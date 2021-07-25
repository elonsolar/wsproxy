package main

import (
	"flag"
	"fmt"
	"net"
	. "github.com/elonsolar/wsproxy/pkg"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/websocket"
)

type Config struct {
	//BindPort   int
	ServerAddr string
	Proxy      []*TcpProxy
}

var (
	cfg string
)

func init() {
	flag.StringVar(&cfg, "cfg", "./client.toml", "")
}

func main() {
	flag.Parse()

	var pxyCfg Config
	if _, err := toml.DecodeFile(cfg, &pxyCfg); err != nil {
		panic(err)
	}

	fmt.Println(pxyCfg)
	proxy(&pxyCfg)
	select {}
}

func proxy(pxyCfg *Config) {
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

	pxy.Serve(con, getWsConnection(serverAddr, pxy))
}

func getWsConnection(serverAddr string, pxy *TcpProxy) *websocket.Conn {

	c, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		panic(err)
	}

	err = c.WriteJSON(pxy)
	if err != nil {
		panic(err)
	}

	return c
}
