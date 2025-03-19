package client

import (
	"errors"
	qk "github.com/quickfixgo/quickfix"
	"github.com/sirupsen/logrus"
	"quickfix-go-tls-demo/conf"
	"quickfix-go-tls-demo/types"
	"time"
)

var (
	MD *QfixInitiator
	OM *QfixInitiator
)

type QfixInitiator struct {
	logonUser     string
	logonPassword string
	sessionId     qk.SessionID
	loggedIn      bool
	initiator     *qk.Initiator
	outChan       chan interface{}
}

func (i *QfixInitiator) OnCreate(sessionID qk.SessionID) {
	i.sessionId = sessionID
}

func (i *QfixInitiator) OnLogon(sessionID qk.SessionID) {
	i.loggedIn = true
}

func (i *QfixInitiator) OnLogout(sessionID qk.SessionID) {
	i.loggedIn = false
}

func (i *QfixInitiator) FromAdmin(msg *qk.Message, sessionID qk.SessionID) (reject qk.MessageRejectError) {
	return nil
}

func (i *QfixInitiator) ToAdmin(msg *qk.Message, sessionID qk.SessionID) {
	if (*qk.Message).IsMsgTypeOf(msg, "A") {
		if i.logonUser != "" && i.logonPassword != "" {
			msg.Body.SetString(553, i.logonUser)
			msg.Body.SetString(554, i.logonPassword)
		} else {
			logrus.Error("Logon user and password not set!!!")
		}
	}
}

func (i *QfixInitiator) ToApp(msg *qk.Message, sessionID qk.SessionID) (err error) {
	return nil
}

func (i *QfixInitiator) FromApp(msg *qk.Message, sessionID qk.SessionID) (reject qk.MessageRejectError) {
	i.outChan <- &types.Message{SessionID: sessionID, Message: msg}
	return nil
}

func (i *QfixInitiator) OnEvent(msg *qk.Message, sessionID qk.SessionID) (reject qk.MessageRejectError) {
	return nil
}

func (i *QfixInitiator) IsLoggedIn() bool {
	return i.loggedIn
}

func (i *QfixInitiator) Start() error {
	err := i.initiator.Start()
	if err != nil {
		return err
	}
	return nil
}

func (i *QfixInitiator) Stop() {
	i.initiator.Stop()
}

func (i *QfixInitiator) SendMsg(message qk.Messagable) error {
	if !i.loggedIn {
		return errors.New("NOT_LOGGED_IN")
	}
	err := qk.SendToTarget(message, i.sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (i *QfixInitiator) setInitiator(initiator *qk.Initiator) {
	i.initiator = initiator
}

func newClient(settings *qk.Settings, msgChan chan interface{}) (*QfixInitiator, error) {
	var err error
	var logonUser, logonPassword string
	if settings.GlobalSettings().HasSetting("LogonUserName") {
		if logonUser, err = settings.GlobalSettings().Setting("LogonUserName"); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("LogonUserName not set")
	}
	if settings.GlobalSettings().HasSetting("LogonPassword") {
		if logonPassword, err = settings.GlobalSettings().Setting("LogonPassword"); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("LogonPassword not set")
	}

	initiator := &QfixInitiator{
		logonUser:     logonUser,
		logonPassword: logonPassword,
		loggedIn:      false,
		outChan:       msgChan,
	}
	ini, err := qk.NewInitiator(initiator, qk.NewMemoryStoreFactory(), settings, NewLogFactory())
	if err != nil {
		return nil, err
	}
	initiator.setInitiator(ini)

	err = initiator.Start()
	if err != nil {
		return nil, err
	}
	return initiator, nil
}

func Init(outChan chan interface{}) (err error) {
	// md
	MD, err = newClient(conf.GetMdSettings(), outChan)
	if err != nil {
		return
	}
	if !waitForTrue(func() bool {
		return MD.IsLoggedIn()
	}, 10*time.Second) {
		return errors.New("MD failed logged in after 10 seconds")
	}
	// om
	OM, err = newClient(conf.GetOmSettings(), outChan)
	if err != nil {
		return
	}
	if !waitForTrue(func() bool {
		return OM.IsLoggedIn()
	}, 10*time.Second) {
		return errors.New("OM failed logged in after 10 seconds")
	}
	return
}

func InitMD(outChan chan interface{}) (err error) {
	MD, err = newClient(conf.GetMdSettings(), outChan)
	if err != nil {
		return
	}
	if !waitForTrue(func() bool {
		return MD.IsLoggedIn()
	}, 10*time.Second) {
		return errors.New("MD failed logged in after 10 seconds")
	}
	return
}

func InitOM(outChan chan interface{}) (err error) {
	OM, err = newClient(conf.GetOmSettings(), outChan)
	if err != nil {
		return
	}
	if !waitForTrue(func() bool {
		return OM.IsLoggedIn()
	}, 10*time.Second) {
		return errors.New("OM failed logged in after 10 seconds")
	}
	return
}

func StopAll() {
	if MD != nil {
		MD.Stop()
	}
	if OM != nil {
		OM.Stop()
	}
}

func waitForTrue(checkFunc func() bool, timeout time.Duration) bool {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ticker.C:
			if checkFunc() {
				return true
			}
		case <-timeoutTimer.C:
			return false
		}
	}
}
