package main

import "net/http"

func main(){
	http.HandleFunc("/ws",WsProxy)

	err:=http.ListenAndServe(":8080",nil)
	if err!=nil{
		panic(err)
	}

}
