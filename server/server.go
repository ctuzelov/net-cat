package server

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"

	assist "net-cat/tools"
)

type Chat struct {
	users   []Usr
	channel chan Message
	history []Message
	mutex   sync.Mutex
}

type Usr struct {
	name string
	conn net.Conn
}

type Message struct {
	user Usr
	rpl  string
}

func (chat *Chat) Massenger() {
	var NewMessage Message

	for {
		NewMessage = <-chat.channel

		chat.mutex.Lock()
		for _, usr := range chat.users {
			if NewMessage.user != usr {
				fmt.Fprint(usr.conn, "\n"+NewMessage.rpl)
				chat.PrintBaseMessage(usr)
			}
		}
		chat.mutex.Unlock()
	}
}

func (chat *Chat) ClientServe(conn net.Conn) {
	userPort := bufio.NewScanner(conn)

	assist.Welcome(conn)

	var name string

	for userPort.Scan() {
		name = userPort.Text()

		err := chat.ValidName(name)
		if err != nil {
			fmt.Fprint(conn, err)
			fmt.Fprint(conn, "[ENTER YOUR NAME]:")
		} else {
			break
		}
	}

	// overall 10 users are acceptable even if there is the limit of 9
	if len(chat.users) > 9 {
		GetResponseChatIsFull(conn, name)
		return
	}

	chat.mutex.Lock()
	chat.restoreHistory(conn)
	chat.mutex.Unlock()

	NewUsr := chat.AddNewUsr(conn, name)

	var NewMessage Message

	NewMessage.user = NewUsr

	chat.PrintBaseMessage(NewUsr)

	for userPort.Scan() {
		NewMessage.rpl = userPort.Text()

		chat.PrintBaseMessage(NewUsr)

		if len(NewMessage.rpl) == 0 {
			continue
		}

		chat.mutex.Lock()
		chat.history = append(chat.history, Message{NewUsr, GetRegularTextMessage(NewMessage)})
		chat.mutex.Unlock()

		chat.channel <- Message{NewUsr, GetRegularTextMessage(NewMessage)}
	}

	chat.DeleteTheUsr(NewUsr)
}

func (chat *Chat) ValidName(name string) error {
	switch {
	case len(name) < 3:
		return fmt.Errorf("Your name should be loneger... \n")
	case chat.Existance(name):
		return fmt.Errorf("This name is already occupied \n")
	default:
		return nil
	}
}

func (chat *Chat) Existance(name string) bool {
	for _, usr := range chat.users {
		if usr.name == name {
			return true
		}
	}
	return false
}

func (chat *Chat) AddNewUsr(conn net.Conn, name string) Usr {
	NewUser := new(Usr)
	NewUser.conn = conn
	NewUser.name = name

	chat.mutex.Lock()
	chat.users = append(chat.users, *NewUser)
	chat.history = append(chat.history, Message{*NewUser, GetClientAddMessage(*NewUser)})
	chat.mutex.Unlock()

	chat.channel <- Message{*NewUser, GetClientAddMessage(*NewUser)}

	return *NewUser
}

func (chat *Chat) DeleteTheUsr(user Usr) {
	chat.mutex.Lock()
	chat.history = append(chat.history, Message{user, GetClientDeleteMessage(user)})
	chat.mutex.Unlock()

	chat.channel <- Message{user, GetClientDeleteMessage(user)}

	chat.mutex.Lock()
	for i, v := range chat.users {
		if v == user {
			chat.users = append(chat.users[:i], chat.users[i+1:]...)
		}
	}
	chat.mutex.Unlock()
}

func GetRegularTextMessage(msg Message) string {
	return fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]: %s\n", msg.user.name, msg.rpl)
}

func (ServerChat *Chat) PrintBaseMessage(client Usr) {
	fmt.Fprint(client.conn, fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]:", client.name))
}

func GetClientAddMessage(client Usr) string {
	return fmt.Sprintf("%s has joined our chat...\n", client.name)
}

func GetClientDeleteMessage(client Usr) string {
	return fmt.Sprintf("%s has left our chat...\n", client.name)
}

func GetResponseChatIsFull(conn net.Conn, name string) {
	fmt.Fprint(conn, fmt.Sprintf("%s, chat is full, try later... \n", name))
}

func (ServerChat *Chat) restoreHistory(conn net.Conn) {
	history := ServerChat.history
	for _, v := range history {
		fmt.Fprint(conn, v.rpl)
	}
}
