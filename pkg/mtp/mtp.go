package mtp

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/xelaj/mtproto/telegram"
)

type Client struct {
	tgClient       *telegram.Client
	Users          map[int32]*telegram.UserObj
	mux            sync.Mutex
	IfShortMessage func(message *telegram.UpdateShortMessage)
	IfNewMessage   func(message *telegram.MessageObj)
}

func New(tgClient *telegram.Client) *Client {
	c := &Client{
		tgClient: tgClient,
	}
	err := c.getUsers()
	if err != nil {
		panic(err)
	}
	c.RunCustomUpdatesHandler()
	return c
}

func (c *Client) RunCustomUpdatesHandler() {
	c.tgClient.AddCustomServerRequestHandler(func(message interface{}) (processed bool) {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println("Recovered with error:", err)
			}
		}()
		fmt.Printf("message type: %T\n", message)
		switch m := message.(type) {
		case *telegram.UpdateShortMessage:
			fmt.Printf("short message %+v\n", m)
			c.IfShortMessage(m)
		case *telegram.UpdateShort:
			fmt.Printf("message type: %T\n", m.Update)
			fmt.Printf("short %+v\n", m)
			switch upd := m.Update.(type) {
			case *telegram.UpdateUserStatus:
				fmt.Printf("user status %+v\n", upd)
			default:
			}
		case *telegram.UpdatesObj:
			fmt.Printf("updates obj message %+v\n", m)
			for _, upd := range m.Updates {
				switch u := upd.(type) {
				case *telegram.UpdateDeleteMessages:
					fmt.Println("deleted messages: ", u.Messages)
				case *telegram.UpdateNewMessage:
					msg := u.Message.(*telegram.MessageObj)
					fmt.Printf("new message: %v\n", msg)
					c.IfNewMessage(msg)
				default:
					fmt.Printf("updates obj type %T\n", u)
				}
			}
		}
		return true
	})
	_, err := c.tgClient.UpdatesGetState()
	if err != nil {
		panic(err)
	}
}

func (c *Client) getUsers() error {
	c.Users = make(map[int32]*telegram.UserObj)
	chats, err := c.tgClient.ContactsGetContacts(0)
	if err != nil {
		return fmt.Errorf("get contacts err: %w", err)
	}
	cts, _ := chats.(*telegram.ContactsContactsObj)
	for _, user := range cts.Users {
		usr := user.(*telegram.UserObj)
		c.Users[usr.ID] = usr
	}
	return nil
}

func (c *Client) SendMessage(userID int32, message string) {
	go func() {
		err := c.sendMessage(userID, message)
		if err != nil {
			panic(err)
		}
	}()
}

func (c *Client) sendMessage(userID int32, message string) error {
	selectedUser, ok := c.Users[userID]
	if !ok {
		return errors.New("user not found")
	}
	c.mux.Lock()
	defer c.mux.Unlock()

	_, err := c.tgClient.MessagesSendMessage(&telegram.MessagesSendMessageParams{
		NoWebpage:  false,
		Silent:     false,
		Background: false,
		ClearDraft: false,
		Peer: &telegram.InputPeerUser{
			UserID:     selectedUser.ID,
			AccessHash: selectedUser.AccessHash,
		},
		ReplyToMsgID: 0,
		Message:      message,
		RandomID:     int64(rand.Int()),
		ScheduleDate: 0,
	})
	if err != nil {
		return fmt.Errorf("sending message error: %w", err)
	}
	return nil
}

func (c *Client) DeleteMessage(messageID int32) {
	go func() {
		err := c.deleteMessage(messageID)
		if err != nil {
			panic(err)
		}
	}()
}

func (c *Client) deleteMessage(messageID int32) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, err := c.tgClient.MessagesDeleteMessages(true, []int32{messageID})
	if err != nil {
		return fmt.Errorf("get contacts err: %w", err)
	}
	return nil
}

func (c *Client) SendMedia(userID int32) {
	go func() {
		err := c.sendMedia(userID)
		if err != nil {
			panic(err)
		}
	}()
}

func (c *Client) sendMedia(userID int32) error {
	selectedUser, ok := c.Users[userID]
	if !ok {
		return errors.New("user not found")
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	_, err := c.tgClient.MessagesSendMedia(&telegram.MessagesSendMediaParams{
		Silent:     false,
		Background: false,
		ClearDraft: false,
		Peer: &telegram.InputPeerUser{
			UserID:     selectedUser.ID,
			AccessHash: selectedUser.AccessHash,
		},
		ReplyToMsgID: 0,
		Media: &telegram.InputMediaDocumentExternal{
			URL:        "https://images.ctfassets.net/hrltx12pl8hq/4plHDVeTkWuFMihxQnzBSb/aea2f06d675c3d710d095306e377382f/shutterstock_554314555_copy.jpg",
			TtlSeconds: 100,
		},
		Message:      "",
		RandomID:     int64(rand.Int()),
		ReplyMarkup:  nil,
		Entities:     nil,
		ScheduleDate: 0,
	})
	return err
}
