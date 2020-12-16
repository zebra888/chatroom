package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"

	"github.com/marcusolsson/tui-go"
	"github.com/zebra888/chatroom/common"
)

const DEBUG = false

func listen(history *tui.Box, ui tui.UI, client *rpc.Client, name string) {
	var chat common.Chat
	go func() {
		for {
			if err := client.Call("ChatRoom.Listen", name, &chat); err != nil {
				// not sure why direct comparison doesn't work. maybe the RPC call errors were manipulated by the package somehow
				if err.Error() == common.ErrorLoggedOut.Error() || err == rpc.ErrShutdown {
					return
				} else {
					fmt.Printf("Error occured in Listen: %v\n", err)
				}
			} else {
				history.Append(tui.NewHBox(
					tui.NewLabel(chat.SendTime),
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", chat.Name))),
					tui.NewLabel(chat.Message),
					tui.NewSpacer(),
				))
				ui.Repaint()

			}
		}
	}()
}

func main() {

	// get rpc client by dailing
	client, _ := rpc.DialHTTP("tcp", "127.0.0.1:9000")

	argsWithoutProg := os.Args[1:]
	var chatterName = "Nobody"
	if len(argsWithoutProg) > 0 {
		chatterName = argsWithoutProg[0]
	}

	if DEBUG {
		f, _ := os.Create(fmt.Sprintf("debug_%s.log", chatterName))
		defer f.Close()
		logger := log.New(f, "", log.LstdFlags)
		tui.SetLogger(logger)
	}

	// chatter login
	var reply string
	if err := client.Call("ChatRoom.Login", chatterName, &reply); err != nil {
		fmt.Println("Error occured when logging in")
	} else {
		fmt.Printf(reply)
	}

	// start TUI
	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		sendTime := time.Now().Format("15:04")
		if err := client.Call("ChatRoom.Post", common.Chat{Name: chatterName, Message: e.Text(), SendTime: sendTime}, nil); err != nil {
			fmt.Println("Error occured in Post")
		}
		//fmt.Printf("input history %v\n", history)
		history.Append(tui.NewHBox(
			tui.NewLabel(sendTime),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", chatterName))),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		input.SetText("")
	})

	ui, err := tui.New(chat)
	if err != nil {
		log.Fatal(err)
	}
	// start listening to get other's message
	listen(history, ui, client, chatterName)
	ui.SetKeybinding("Esc", func() {
		if err := client.Call("ChatRoom.Logout", chatterName, &reply); err != nil {
			fmt.Println("Error occured in Logout")
		} else {
			inputBox.Append(tui.NewHBox(
				tui.NewLabel(reply),
			))
		}
		ui.Quit()
	})

	// config TUI
	// onsubmit:
	// call post (wrap up a function),
	// while listening, append to history box
	ui.Run()
}
