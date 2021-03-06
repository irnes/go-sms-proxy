package service

import (
	"go-sms-proxy/app/model"
)

// Request is used for making requests to SMSGateway services
type Request struct {
	Payload *model.BaseMessage
	RspChan chan Response
}

// Response is the value returned by sender
type Response interface{}

// SMSProvider defines an interface for SMS Sender
type SMSProvider interface {
	Balance()
	Send(*model.BaseMessage) <-chan Response
}

// NewSMSService creates new instance of sms service using given provider
func NewSMSService(provider SMSProvider) *SMSService {
	service := &SMSService{SMSProvider: provider}
	return service
}

// SMSService uses given sms provider to send text messages
type SMSService struct {
	SMSProvider
}
