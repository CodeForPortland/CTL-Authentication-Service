package connection

import (
	"log"
	"os"
	"os/signal"

	"github.com/CodeForPortland/ctl-authentication-service/authentication"
	"github.com/CodeForPortland/ctl-authentication-service/reporter"

	"github.com/gammazero/nexus/client"
)

type Options struct {
	Realm       string
	Logger      *log.Logger
	Url         string
	Reporter    *reporter.Reporter
}

const (
	realm = "aura.ctl.public"
	wsAddr   = "127.0.0.1:3131"
	wssAddr  = "localhost:4242"
	tcpAddr  = "127.0.0.1:8001"
	tcpsAddr = "localhost:8101"
	unixAddr = "/tmp/example_aura_sock"
)

func RegisterService() {
	const (
		procedureName = "ctl.authentication"
	)
	logger := log.New(os.Stderr, "CTL-AUTH Caller> ", 0)

	clientAddrs := Addresses{
		WsAddr: wsAddr,
		WssAddr: wssAddr,
		TcpAddr: tcpAddr,
		TcpsAddr: tcpsAddr,
		UnixAddr: unixAddr,
	}

	clientCfg := client.Config{
		Realm: realm,
		Logger: logger,
	}

	// Connect serviceClient client with requested socket type and serialization.
	serviceClient, err := NewClient(clientAddrs, clientCfg, "ws")
	if err != nil {
		logger.Fatal(err)
	}
	defer serviceClient.Close()

	serviceClient.Register(procedureName, authentication.AuthenticateHandler, nil)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case <-sigChan:
	case <-serviceClient.Done():
		logger.Print("Router gone, exiting")
		return // router gone, just exit
	}

	if err = serviceClient.Unregister(procedureName); err != nil {
		logger.Println("Failed to unregister procedure:", err)
	}
}