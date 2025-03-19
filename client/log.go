package client

import (
	"fmt"
	"github.com/quickfixgo/quickfix"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	SOH = byte(1)
)

type screenLog struct {
	prefix string
}

func (l screenLog) OnIncoming(s []byte) {
	logrus.Infof("in: %s", strings.Replace(string(s), string(SOH), "|", -1))
}

func (l screenLog) OnOutgoing(s []byte) {
	logrus.Info(fmt.Sprintf("out: %s", strings.Replace(string(s), string(SOH), "|", -1)))
}

func (l screenLog) OnEvent(s string) {
	logrus.Infof("event: %s", s)
}

func (l screenLog) OnEventf(format string, a ...interface{}) {
	l.OnEvent(fmt.Sprintf(format, a...))
}

type logFactory struct{}

func (logFactory) Create() (quickfix.Log, error) {
	log := screenLog{"GLOBAL"}
	return log, nil
}

func (logFactory) CreateSessionLog(sessionID quickfix.SessionID) (quickfix.Log, error) {
	log := screenLog{sessionID.String()}
	return log, nil
}

func NewLogFactory() quickfix.LogFactory {
	return logFactory{}
}
