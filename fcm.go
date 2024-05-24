package fcm_gateway

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
	"time"
)

type Client struct {
	fcmClient *messaging.Client
}
type ClientConfig struct {
	Cred string
}
type ClientSendResults struct {
	SuccessCount int
	FailureCount int
	Results      []*SendResult
}
type SendResult struct {
	Token     string
	Success   bool
	MessageId string
	Error     error
}

func New(cfg *ClientConfig) (*Client, error) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, &firebase.Config{}, option.WithCredentialsJSON([]byte(cfg.Cred)))
	if err != nil {
		return nil, err
	}
	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{fcmClient: fcmClient}, nil
}

func (client *Client) SendMessage(ctx context.Context, body string, title string, tokens []string, data map[string]string) (*ClientSendResults, error) {
	cont, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	response, err := client.fcmClient.SendEachForMulticast(cont, &messaging.MulticastMessage{
		Data: data,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
			//ImageURL: "https://mcs.nlmk.com/img/app-store.svg",
		},
		Tokens: tokens,
	})
	if err != nil {
		return nil, err
	}
	var sendResults []*SendResult
	for i, token := range tokens {
		success := false
		messageID := ""
		var err error

		if resp := response.Responses[i]; resp != nil {
			success = resp.Success
			messageID = resp.MessageID
			err = resp.Error
		}
		sendResults = append(sendResults, &SendResult{
			Token:     token,
			Success:   success,
			MessageId: messageID,
			Error:     err,
		})
	}

	return &ClientSendResults{
		SuccessCount: response.SuccessCount,
		FailureCount: response.FailureCount,
		Results:      sendResults,
	}, nil
}
