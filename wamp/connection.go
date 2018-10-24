package wamp

import (
	"ctl-auth/reporter"
	"github.com/gammazero/nexus/wamp"
	"log"

	"github.com/gammazero/nexus/client"
)

type ConnectionOptions struct {
	Realm       string
	Logger      *log.Logger
	Url         string
	Reporter    *reporter.Reporter
}

func connect(opt ConnectionOptions) {
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

	routerEvtHandler := func(args wamp.List, kwargs wamp.Dict, details wamp.Dict) {
		opt.Reporter.MakeReport(reporter.LISTENING, "WAMP router connection established and topic subscrition returned.", nil)

	}

	RegHandler := func(args wamp.List, kwargs wamp.Dict, details wamp.Dict) {
		if args[0] == "REGISTERED" {
			topic := kwargs["topic"].(string)
			subscriber.Subscribe(topic, routerEvtHandler, nil)
			opt.Reporter.MakeReport(reporter.EVENT, "WAMP router is registering authentication component.", map[string]interface{}{
				"args": args, "kwargs": kwargs, "details": details,
			})
		}
		opt.Logger.Println("Router connection received event.")
	}

	regDetails := wamp.Dict{
		"id": "CTLAuthService0",
		"topic": "Authentication",
	}

	subscriber.Subscribe("Register", RegHandler, regDetails)

	defer opt.Reporter.MakeReport(reporter.CLOSING, "Closing WAMP router connection.", nil)
}
