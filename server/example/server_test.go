package example

import (
	"github.com/elonsolar/wsproxy/server"
	"net/http"
	"testing"
)

// go test
func TestServer(t *testing.T){
	http.HandleFunc("/ws",server.WsProxy)

	err:=http.ListenAndServe(":8080",nil)
	if err!=nil{
		panic(err)
	}

}
