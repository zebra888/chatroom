package common

import (
	"errors"
	"fmt"
	"sync"
)

var ErrorLoggedOut = errors.New("User has already logged out")
var ErrorAlreadyLoggedIn = errors.New("User has already logged in")
var ErrorNotLoggedIn = errors.New("Chatter does not exist or not logged in")

type ChatRoom struct {
	chatterDB map[string]bool
	msgChan   map[string]chan string
	mu        sync.RWMutex
}

type Chat struct {
	Name, Message string
}

/*func (c *ChatRoom) Register(name string, pwd string) error {
	return errors.New("Register -- Not Implemented")
}*/

// login to chat room. only provide name for now. will add password later
func (c *ChatRoom) Login(name string, reply *string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	chatterLoggedIn, ok := c.chatterDB[name]
	if ok {
		if chatterLoggedIn {
			return ErrorAlreadyLoggedIn
		}
	}

	// new chatter or existing chatter login: init channel and set login status
	c.msgChan[name] = make(chan string)
	c.chatterDB[name] = true
	*reply = fmt.Sprintf("Welcome %s\n", name)

	return nil
}

func (c *ChatRoom) Logout(name string, reply *string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	chatterLoggedIn, ok := c.chatterDB[name]
	if !ok || chatterLoggedIn == false {
		return ErrorNotLoggedIn
	}
	c.chatterDB[name] = false
	close(c.msgChan[name])
	*reply = fmt.Sprintf("Bye %s\n", name)
	return nil
}

// postreply can be ignored, it's for matching RPC call function format
func (c *ChatRoom) Post(chat Chat, reply *string) error {
	// since login is already checked, no need to check existence here
	// broadcast message to everyone in room except sender itself
	c.mu.RLock()
	defer c.mu.RUnlock()

	for chatter, ch := range c.msgChan {
		if chatter != chat.Name && c.chatterDB[chatter] == true {
			// client should already be listening in go routine so this should not block
			ch <- fmt.Sprintf("%s: %s", chat.Name, chat.Message)
			fmt.Printf("posted message %s to %s\n", chat.Message, chatter)
		}
	}
	*reply = "Posted"
	return nil
}

func (c *ChatRoom) Listen(name string, reply *string) error {
	c.mu.RLock()
	ch, _ := c.msgChan[name]
	c.mu.RUnlock()

	// each call of Listen retrives one message from channel, because rpc doesn't support streaming
	msg, ok := <-ch
	if !ok {
		return ErrorLoggedOut
	}
	*reply = msg
	return nil
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		chatterDB: make(map[string]bool),
		msgChan:   make(map[string]chan string),
	}
}
