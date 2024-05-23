package firebase

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"google.golang.org/api/option"
	"time"
)

type Client struct {
	fcmClient *messaging.Client
}
type ClientConfig struct {
	Creditions string
}

func New(cfg *ClientConfig) (*Client, error) {
	ctx := context.Background()
	decodeString, err := base64.StdEncoding.DecodeString(cfg.Creditions)
	if err != nil {
		return nil, err
	}
	opt1 := option.WithCredentialsJSON(decodeString)
	app, err := firebase.NewApp(ctx, &firebase.Config{}, opt1)
	if err != nil {
		return nil, err
	}
	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{fcmClient: fcmClient}, nil
}

func (client *Client) SendMessage(ctx context.Context, body string, title string, tokens []string, data map[string]string) (*messaging.BatchResponse, error) {
	startTime := time.Now()
	message := &messaging.MulticastMessage{
		Data: data,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
			//ImageURL: "https://mcs.nlmk.com/img/app-store.svg",
		},
		Tokens: tokens,
	}
	cont, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	response, err := client.fcmClient.SendEachForMulticast(cont, message)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Successfully sent message: %d, Failure sent message: %d, error: %v, time: %s\n", response.SuccessCount, response.FailureCount, err, time.Since(startTime))
	return response, nil
}
