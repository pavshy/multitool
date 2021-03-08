package mtp

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xelaj/mtproto"
	utils "github.com/xelaj/mtproto/examples/example_utils"
	"github.com/xelaj/mtproto/telegram"

	"multitool/pkg/config"
)

func SignIn(appConf *config.App) (*telegram.Client, error) {

	client, err := telegram.NewClient(telegram.ClientConfig{
		SessionFile:     appConf.SessionFile,
		ServerHost:      appConf.TgProdHost,
		PublicKeysFile:  appConf.PublicKeys,
		AppID:           appConf.AppID,
		AppHash:         appConf.AppHash,
		InitWarnChannel: true,
	})
	if err != nil {
		panic(err)
	}
	client.Warnings = make(chan error)
	utils.ReadWarningsToStdErr(client.Warnings)

	_, err = os.Stat(appConf.SessionFile)
	if err == nil {
		println("You've already signed in!")
		return client, nil
	}

	setCode, err := client.AuthSendCode(
		appConf.PhoneNumber, int32(appConf.AppID), appConf.AppHash, &telegram.CodeSettings{},
	)
	fmt.Println("phone err:", setCode, err)

	// this part shows how to deal with errors (if you want of course. No one
	// like errors, but the can be return sometimes)
	if err != nil {
		errResponse := &mtproto.ErrResponseCode{}
		if !errors.As(err, &errResponse) {
			// some strange error, looks like a bug actually
			fmt.Println(err)
			return nil, err
		} else {
			if errResponse.Message == "AUTH_RESTART" {
				println("Oh crap! You accidentaly restart authorization process!")
				println("You should login only once, if you'll spam 'AuthSendCode' method, you can be")
				println("timeouted to loooooooong long time. You warned.")
			} else if errResponse.Message == "FLOOD_WAIT_X" {
				println("No way... You've reached flood timeout! Did i warn you? Yes, i am. That's what")
				println("happens, when you don't listen to me...")
				println()
				timeoutDuration := time.Second * time.Duration(errResponse.AdditionalInfo.(int))

				println("Repeat after " + timeoutDuration.String())
			} else {
				println("Oh crap! Got strange error:")
				fmt.Println(errResponse)
			}

			return nil, err
		}
	}

	fmt.Print("Auth code: ")
	code, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	code = strings.ReplaceAll(code, "\n", "")

	auth, err := client.AuthSignIn(
		appConf.PhoneNumber,
		setCode.PhoneCodeHash,
		code,
	)
	if err == nil {
		fmt.Println(auth)

		fmt.Println("Success! You've signed in!")
		return client, nil
	}
	//
	//// if you don't have password protection — THAT'S ALL! You're already logged in.
	//// but if you have 2FA, you need to make few more steps:
	//
	//// could be some errors:
	//errResponse := &mtproto.ErrResponseCode{}
	//ok := errors.As(err, &errResponse)
	//// checking that error type is correct, and error msg is actualy ask for password
	//if !ok || errResponse.Message != "SESSION_PASSWORD_NEEDED" {
	//	fmt.Println("SignIn failed:", err)
	//	println("\n\nMore info about error:")
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//fmt.Print("Password:")
	//password, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	//password = strings.ReplaceAll(password, "\n", "")
	//
	//accountPassword, err := client.AccountGetPassword()
	//if err != nil {
	//	return err
	//}
	//
	//// GetInputCheckPassword is fast response object generator
	//inputCheck, err := telegram.GetInputCheckPassword(password, accountPassword)
	//if err != nil {
	//	return err
	//}
	//
	//auth, err = client.AuthCheckPassword(inputCheck)
	//if err != nil {
	//	return err
	//}

	fmt.Println(auth)
	fmt.Println("Success! You've signed in!")
	return client, nil
}
