package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xelaj/mtproto/telegram"

	"multitool/pkg/config"
	"multitool/pkg/mtp"
)

var chatID int32 = -1 // Input desired chat id here

func main() {
	appConf, err := config.ReadFromFile("static/config.yml")
	if err != nil {
		panic(fmt.Errorf("error reading config file: %w", err))
	}
	tgClient, err := mtp.SignIn(appConf)
	if err != nil {
		panic(fmt.Errorf("error signing in: %w", err))
	}
	defer tgClient.Stop()

	client := mtp.New(tgClient)

	for _, user := range client.Users {
		fmt.Println(user.ID, user.Username)
	}

	chatsInterface, err := tgClient.MessagesGetAllChats([]int32{})
	if err != nil {
		panic(fmt.Errorf("getting chats error: %w", err))
	}
	chats, _ := chatsInterface.(*telegram.MessagesChatsObj)
	for _, chatInterface := range chats.Chats {
		chat, _ := chatInterface.(*telegram.ChatObj)
		fmt.Println(chat.ID, chat.Title)
	}

	for i := 0; i <= 100; i++ {
		client.SendMessage(chatID, strconv.Itoa(i))
		time.Sleep(time.Millisecond * 10)
	}

	//client.IfShortMessage = func(message *telegram.UpdateShortMessage) {
	//	client.DeleteMessage(message.ID)
	//	client.SendMessage(message.UserID, "go: "+message.Message)
	//	client.SendMedia(message.UserID)
	//}
	//client.IfNewMessage = func(message *telegram.MessageObj) {
	//	peer := message.PeerID.(*telegram.PeerUser)
	//	client.DeleteMessage(message.ID)
	//	client.SendMessage(peer.UserID, "go: "+message.Message)
	//	client.SendMedia(peer.UserID)
	//}
}
