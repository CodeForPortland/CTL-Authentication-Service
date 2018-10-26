package wamp

import (
	"ctl-auth/reporter"
	"github.com/gammazero/nexus/wamp"
	"log"
	"os"
	"os/signal"

	"github.com/gammazero/nexus/client"
)

type ConnectionOptions struct {
	Realm       string
	Logger      *log.Logger
	Url         string
	Reporter    *reporter.Reporter
}

func ConnectService(opt ConnectionOptions) {
	const RegistrationTopic = "Register::Service"

	opt.Reporter.MakeReport(reporter.INITIALIZING,"Starting WAMP router Connection.", nil)
	cfg := client.Config{
		Realm: opt.Realm,
		Logger: opt.Logger,

	}

	subscriber, err := client.ConnectNet(opt.Url, cfg)
	if err != nil {
		opt.Logger.Fatal(err)
	}
	defer subscriber.Close()

	//routerEvtHandler := func(args wamp.List, kwargs wamp.Dict, details wamp.Dict) {
	//	opt.Reporter.MakeReport(reporter.LISTENING, "WAMP router connection established and topic subscription returned.", nil)
	//}

	regEvtHandler := func(args wamp.List, kwargs wamp.Dict, details wamp.Dict) {
		//if args[0] == "REGISTERED" {
		//	topic := kwargs["topic"].(string)
		//	subscriber.Subscribe(topic, routerEvtHandler, nil)
		//	opt.Reporter.MakeReport(reporter.EVENT, "WAMP router is registering authentication component.", map[string]interface{}{
		//		"args": args, "kwargs": kwargs, "details": details,
		//	})
		//}
		opt.Logger.Println("Router connection received event.")
	}

	regDetails := wamp.Dict{
		"id": "CTLAuthService0",
		"topic": "CTL.Authentication",
		"realm": "ctl.authentication",
	}

	err = subscriber.Subscribe(RegistrationTopic, regEvtHandler, regDetails)
	if err != nil {
		opt.Logger.Fatal(err)
	}
	err = subscriber.Publish(RegistrationTopic, wamp.Dict{"acknowledge": true}, nil, regDetails)
	if err != nil {
		opt.Logger.Fatal(err)
	}
	opt.Logger.Println("Test")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	select {
		case<-sigChan:
		case<-subscriber.Done():
			opt.Logger.Println("Sub done?")
			opt.Reporter.MakeReport(reporter.CLOSING, "Closing WAMP router connection.", nil)
			return
	}

	if err = subscriber.Unsubscribe(RegistrationTopic); err != nil {
		opt.Logger.Fatal("Failed to unsubscribe", RegistrationTopic, err)
	}
}
