package internal

import (
	sdk "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/sms"
	log "github.com/sirupsen/logrus"
)

type SmsClient interface {
	sendSms(msg *Messages) error
}

type DefaultSmsClient struct {
	c *sdk.Client
}

func NewDefaultClient(accessKey string) *DefaultSmsClient {
	return &DefaultSmsClient{
		c: sdk.New(accessKey),
	}
}

func (cl *DefaultSmsClient) sendSms(msg *Messages) error {
	message, err := sms.Create(cl.c, msg.Originator, msg.Recipients, msg.Message,nil)
	if err!= nil {
		return err
	}
	log.Println("Messagebird Message:", message)
	return nil
}