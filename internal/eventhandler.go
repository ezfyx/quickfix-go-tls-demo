package internal

import (
	"github.com/sirupsen/logrus"
	"quickfix-go-tls-demo/types"
)

var log = logrus.WithField("pkg", "internal")

type Cerebro struct {
	inChan chan interface{}
}

func (c *Cerebro) eventLoop() {
	for {
		select {
		case msg := <-c.inChan:
			switch msg.(type) {
			case *types.Message:
				m := msg.(*types.Message)
				log.Infof("session: %s, msg: %s", m.SessionID, m.Message.String())
			}
		default:
		}
	}
}

func (c *Cerebro) Input() chan interface{} {
	return c.inChan
}

func NewCerebro() *Cerebro {
	c := &Cerebro{
		inChan: make(chan interface{}, 1000),
	}
	go c.eventLoop()
	return c
}
