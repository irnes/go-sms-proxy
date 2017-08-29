package app

import (
	"log"
	"strconv"

	"github.com/messagebird/go-rest-api"
)

// Request is used for making requests to SMSGateway services
type Request struct {
	Payload *BaseMessage
	RspChan chan Response
}

// Response is the value returned by sender
type Response interface{}

//{"recipient":31612345678,"originator":"MessageBird","message":"This is a test message."}

// Sender defines an interface for SMS Sender
type Sender interface {
	Send(*BaseMessage) <-chan Response
}

// NewSMSSender creates an instance of SMSSender object using provided sms client
func NewSMSSender(client *messagebird.Client) *SMSSender {
	sender := &SMSSender{
		SMSClient: client,
		Done:      make(chan bool, 0),
		ReqChan:   make(chan Request, 10),
	}

	// start sender worker
	go sender.Run()

	return sender
}

// SMSSender is service for sending text messages through provided SMSClient
type SMSSender struct {
	SMSClient *messagebird.Client
	Done      chan bool
	ReqChan   chan Request
}

// Terminate SMS worker
func (s *SMSSender) Terminate() {
	s.Done <- true
}

// Send encapsulate given payload into the Request
// structure and enques it for further processing
// returns reposnse chanel
func (s *SMSSender) Send(payload *BaseMessage) <-chan Response {
	req := Request{
		Payload: payload,
		RspChan: make(chan Response, 1),
	}
	// enqueue request for processing
	s.ReqChan <- req
	return req.RspChan
}

// Run will make the SMS start listening to the two channels ctx.Done and ReqChan
func (s *SMSSender) Run() {
	for {
		select {
		case <-s.Done:
			close(s.ReqChan)
			log.Println("Terminating SMS Worker ...")
			return
		case req := <-s.ReqChan:
			req.RspChan <- s.process(req.Payload)
		}
	}
}

func (s *SMSSender) process(bm *BaseMessage) Response {
	originator := bm.Originator
	recipients := []string{strconv.Itoa(bm.Recipient)}

	message, _ := s.SMSClient.NewMessage(originator, recipients, bm.Message, nil)

	return message
}
