package connection

import (
	"bytes"
	"context"
	"crypto/sha256"
	"ctl-auth/reporter"
	"encoding/base64"
	"github.com/gammazero/nexus/wamp"
	"log"
	"os"
	"strings"

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
		procedureName = "register.service"
		serviceName = "ctl.authentication"
	)
	logger := log.New(os.Stderr, "CTL-AUTH Caller> ", 0)

	clientAddrs := Addresses{
		WsAddr: wsAddr,
		WssAddr: wssAddr,
		TcpAddr:tcpAddr,
		TcpsAddr: tcpsAddr,
		UnixAddr: unixAddr,
	}

	clientCfg := client.Config{
		Realm: realm,
		Logger: logger,
	}

	// Connect caller client with requested socket type and serialization.
	caller, err := NewClient(logger, clientAddrs, clientCfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer caller.Close()
	ctx := context.Background()

	result, err := caller.Call(
		ctx, procedureName, nil, nil, wamp.Dict{"procName": serviceName}, "")
	if err != nil {
		logger.Println("Err of reg", err)
	}
	logger.Println("Result of reg", result)
	var authenticationHandler client.InvocationHandler = func(i context.Context, lists wamp.List, dicts wamp.Dict, dicts2 wamp.Dict) (result *client.InvokeResult) {
		logger.Println("Authentication invoked \n", i, lists, dicts, dicts2)
		return &client.InvokeResult{Args: wamp.List{"Tada!"}}
	}

	caller.Register("ctl.authentication", authenticationHandler, nil)
}

func progressiveCall(logger log.Logger, caller *client.Client, ctx context.Context) {
	const (
		procedureName = "Unknown"
		chunkSize = 64
	)

	// The progress handler accumulates the chunks of data as they arrive.  It
	// also progressively calculates a sha256 hash of the data as it arrives.
	var chunks []string
	h := sha256.New()
	progHandler := func(result *wamp.Result) {
		// Received another chunk of data, computing hash as chunks received.
		chunk := result.Arguments[0].(string)
		logger.Println("Received", len(chunk), "bytes (as progressive result)")
		chunks = append(chunks, chunk)
		h.Write([]byte(chunk))
	}



	// Call the example procedure, specifying the size of chunks to send as
	// progressive results.
	result, err := caller.CallProgress(
		ctx, procedureName, nil, wamp.List{chunkSize}, nil, "", progHandler)
	if err != nil {
		logger.Println("Failed to call procedure:", err)
		return
	}

	// As a final result, the callee returns the base64 encoded sha256 hash of
	// the data.  This is decoded and compared to the value that the caller
	// calculated.  If they match, then the caller recieved the data correctly.
	hashB64 := result.Arguments[0].(string)
	calleeHash, err := base64.StdEncoding.DecodeString(hashB64)
	if err != nil {
		logger.Println("decode error:", err)
		return
	}

	// Check if rceived hash matches the hash computed over the received data.
	if !bytes.Equal(calleeHash, h.Sum(nil)) {
		logger.Println("Hash of received data does not match")
		return
	}
	logger.Println("Correctly received all data:")
	logger.Println("----------------------------")
	logger.Println(strings.Join(chunks, ""))
}