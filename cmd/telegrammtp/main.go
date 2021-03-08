package main

import (
	"fmt"
	"time"

	"github.com/xelaj/mtproto/telegram"

	"multitool/pkg/config"
	"multitool/pkg/mtp"
)

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

	client.IfShortMessage = func(message *telegram.UpdateShortMessage) {
		client.DeleteMessage(message.ID)
		client.SendMessage(message.UserID, "go: "+message.Message)
		client.SendMedia(message.UserID)
	}
	client.IfNewMessage = func(message *telegram.MessageObj) {
		peer := message.PeerID.(*telegram.PeerUser)
		client.DeleteMessage(message.ID)
		client.SendMessage(peer.UserID, "go: "+message.Message)
		client.SendMedia(peer.UserID)
	}

	<-time.After(time.Hour)
}
