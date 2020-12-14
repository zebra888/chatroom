package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"sync"

	"github.com/zebra888/chatserver/common"
)

func chat(client *rpc.Client, name string) {
	reader := bufio.NewReader(os.Stdin)
	var reply string
	for {
		msg, _ := reader.ReadString('\n')

		// empty string will stop chatting and logout
		if msg == "\n" {
			if err := client.Call("ChatRoom.Logout", name, &reply); err != nil {
				fmt.Println("Error occured in Logout")
			} else {
				fmt.Println(reply)
				return
			}
		}
		if err := client.Call("ChatRoom.Post", common.Chat{Name: name, Message: msg}, nil); err != nil {
			fmt.Println("Error occured in Post")
		}
	}
}

func main() {
	// get rpc client by dailing
	client, _ := rpc.DialHTTP("tcp", "127.0.0.1:9000")

	argsWithoutProg := os.Args[1:]
	var chatterName = "Nobody"
	if len(argsWithoutProg) > 0 {
		chatterName = argsWithoutProg[0]
	}

	var reply string

	// chatter login
	if err := client.Call("ChatRoom.Login", chatterName, &reply); err != nil {
		fmt.Println("Error occured when logging in")
	} else {
		fmt.Printf(reply)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Call("ChatRoom.Listen", chatterName, &reply); err != nil {
				// not sure why direct comparison doesn't work. maybe the RPC call errors were manipulated by the package somehow
				if err.Error() == common.ErrorLoggedOut.Error() {
					return
				} else {
					fmt.Printf("Error occured in Listen: %v\n", err)
				}
			} else {
				fmt.Println(reply)
			}
		}
	}()

	// start chatting
	chat(client, chatterName)

	// wait until chatter logout
	wg.Wait()
}
