package repository

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/rimantoro/event_driven/profiler/entities/gowa"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	"go.uber.org/zap"

	"github.com/rimantoro/event_driven/profiler/shared/bootstrap"
)

type waRepository struct {
	// DB         dbHelper.Database
	// Collection dbHelper.Collection
}

func NewWaRepository() gowa.Repository {
	return &waRepository{}
}

func (w *waRepository) SendMessage(c context.Context, number string, message string) (string, error) {

	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		bootstrap.App.Logger.Error("error creating whatsapp connection", zap.Error(err))
		return "", err
	}

	err = login(wac)
	if err != nil {
		bootstrap.App.Logger.Error("error logging in whatsapp", zap.Error(err))
		return "", err
	}

	<-time.After(3 * time.Second)

	previousMessage := "😘"
	quotedMessage := proto.Message{
		Conversation: &previousMessage,
	}

	ContextInfo := whatsapp.ContextInfo{
		QuotedMessage:   &quotedMessage,
		QuotedMessageID: "",
		Participant:     "", //Whot sent the original message
	}

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: number + "@s.whatsapp.net",
		},
		ContextInfo: ContextInfo,
		Text:        message,
	}

	msgId, err := wac.Send(msg)
	if err != nil {
		bootstrap.App.Logger.Error("error sending message of whatsapp", zap.Error(err))
		return msgId, err
	} else {
		return msgId, nil
	}

}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	// file, err := os.Open(os.TempDir() + "~/whatsappSession.gob")
	file, err := os.Open(bootstrap.App.Config.GetString("wa.session_path") + "whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	// file, err := os.Create(os.TempDir() + "~/whatsappSession.gob")
	file, err := os.Create(bootstrap.App.Config.GetString("wa.session_path") + "whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}
