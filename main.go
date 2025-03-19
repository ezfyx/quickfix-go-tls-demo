package main

import (
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	fix44mdr "github.com/quickfixgo/fix44/marketdatarequest"
	"os"
	"os/signal"
	"quickfix-go-tls-demo/client"
	"quickfix-go-tls-demo/conf"
	"quickfix-go-tls-demo/internal"
	"quickfix-go-tls-demo/utils"
	"syscall"
)

func main() {
	err := conf.LoadConfig()
	if err != nil {
		utils.Fatal("config.LoadConfig()", err)
	}
	cerebro := internal.NewCerebro()
	err = client.Init(cerebro.Input())
	if err != nil {
		utils.Fatal("client.Init()", err)
	}

	req := fix44mdr.New(field.NewMDReqID("MARKETDATAID"),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT),
		field.NewMarketDepth(10),
	)
	g1 := fix44mdr.NewNoMDEntryTypesRepeatingGroup()
	g1.Add().SetMDEntryType(enum.MDEntryType_BID)
	g1.Add().SetMDEntryType(enum.MDEntryType_OFFER)
	req.SetNoMDEntryTypes(g1)
	g2 := fix44mdr.NewNoRelatedSymRepeatingGroup()
	g2.Add().SetSymbol("ETHUSD")
	req.SetNoRelatedSym(g2)
	err = client.MD.SendMsg(req)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
	client.StopAll()
	return
}
