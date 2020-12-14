package main

import (
	"io"
	"net/http"
	"net/rpc"

	"github.com/zebra888/chatserver/common"
)

func main() {
	// create a chatRoom object
	chatRoom := common.NewChatRoom()

	// register chatroom with rpc
	rpc.Register(chatRoom)

	// a few things happened within this call:
	// 1. a defutl http handler is registered
	// 2. a default rpc path to respond to rpc messages is registered
	// 3. a default rpc debug endpoint is registered
	rpc.HandleHTTP()

	// liveness probe
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		io.WriteString(resp, "Chatroom live!")
	})

	// listen
	http.ListenAndServe(":9000", nil)
}
