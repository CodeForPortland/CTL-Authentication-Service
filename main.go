package main

import (
	"ctl-auth/reporter"
	"ctl-auth/wamp"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "ctl-auth: ", 0)
	rptChan := make(chan reporter.Report, 0)
	mainReporter := reporter.Reporter{
		Updater: rptChan,
	}
	conOpts := wamp.ConnectionOptions{
		Url: "ws://127.0.0.1:3131",
		Realm: "aura.registration",
		Logger: logger,
		Reporter: &mainReporter,
	}
	go wamp.ConnectService(conOpts)


	for onNewReport := range rptChan {
		logger.Println(onNewReport.ReportType().String())
	}
}
