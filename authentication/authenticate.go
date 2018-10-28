package authentication

import (
	"context"
	"crypto/rand"
	"github.com/dgrijalva/jwt-go"
	"github.com/gammazero/nexus/client"
	"github.com/gammazero/nexus/wamp"
	"log"
	"os"
)

type Credentials struct {
	username string
	password string
}

var logger = log.New(os.Stdout, "Authentication> ", log.LstdFlags)

var AuthenticateHandler client.InvocationHandler = func(context context.Context, lists wamp.List, dicts wamp.Dict, dicts2 wamp.Dict) (result *client.InvokeResult) {
	logger.Println("Auth Invoked", dicts, dicts2, lists, context)

	var (
		username = dicts["username"]
		password = dicts["password"]
	)

	if username != nil && password != nil {

		if len(username.(string)) > 0 && len(password.(string)) > 0 {

			// Needs to Auth

			// Makes a JWT
			key := make([]byte, 128)
			rand.Read(key)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"foo": "bar",
				"signedIn": true,
				//"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
			})

			// Sign and get the complete encoded token as a string using the secret
			tokenString, err := token.SignedString(key)

			logger.Println(tokenString, err)

			return &client.InvokeResult{Kwargs: wamp.Dict{"jwt": tokenString}}
		}
	}
	return &client.InvokeResult{Err: wamp.ErrCanceled}
}

