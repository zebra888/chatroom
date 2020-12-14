package common

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	c := NewChatRoom()
	var reply string
	chatter := "Guest"
	err := c.Login(chatter, &reply)
	assert.Equal(t, err, nil)
	assert.Equal(t, reply, fmt.Sprintf("Welcome %s\n", chatter))

	err = c.Login("Guest", &reply)
	assert.Equal(t, err, ErrorAlreadyLoggedIn)
}

func TestLogout(t *testing.T) {
	c := NewChatRoom()
	var reply string
	chatter := "Guest"
	err := c.Logout(chatter, &reply)
	assert.Equal(t, err, ErrorNotLoggedIn)
	c.Login(chatter, &reply)
	err = c.Logout(chatter, &reply)
	assert.Equal(t, err, nil)
	assert.Equal(t, reply, fmt.Sprintf("Bye %s\n", chatter))
}

func TestPostandListen(t *testing.T) {
	c := NewChatRoom()
	var reply1, reply2 string
	chatter1 := "Guest1"
	chatter2 := "Guest2"

	chat1 := Chat{Name: chatter1, Message: fmt.Sprintf("Hi, I'm %s\n", chatter1)}
	chat2 := Chat{Name: chatter2, Message: fmt.Sprintf("Hi, I'm %s\n", chatter2)}

	_ = c.Login(chatter1, &reply1)
	_ = c.Login(chatter2, &reply2)

	// start listening first
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := c.Listen(chatter1, &reply1)
		assert.Equal(t, err, nil)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := c.Listen(chatter2, &reply2)
		assert.Equal(t, err, nil)
	}()

	var postreply string

	// postreply can be ignored, it's for matching RPC call function format
	err := c.Post(chat1, &postreply)
	assert.Equal(t, err, nil)
	err = c.Post(chat2, &postreply)
	assert.Equal(t, err, nil)

	wg.Wait()
	// chatter should receive messae posted by the others
	assert.Equal(t, reply1, fmt.Sprintf("%s: %s", chat2.Name, chat2.Message))
	assert.Equal(t, reply2, fmt.Sprintf("%s: %s", chat1.Name, chat1.Message))
}
