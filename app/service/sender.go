package service

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/messagebird/go-rest-api"

	"mbsms-api/app/model"
	"mbsms-api/app/util"
)

// Request is used for making requests to SMSGateway services
type Request struct {
	Payload *model.BaseMessage
	RspChan chan Response
}

// Response is the value returned by sender
type Response interface{}

// Sender defines an interface for SMS Sender
type Sender interface {
	Send(*model.BaseMessage) <-chan Response
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

// Send encapsulate given base message into the Request
// structure and enqueues it for further processing
// returns reposnse chanel
func (s *SMSSender) Send(bm *model.BaseMessage) <-chan Response {
	req := Request{
		Payload: bm,
		RspChan: make(chan Response, 1),
	}
	// enqueue request for processing
	s.ReqChan <- req
	return req.RspChan
}

// Run will make the SMSSender starts listening to the two channels ctx.Done and ReqChan
func (s *SMSSender) Run() {
	for {
		select {
		case <-s.Done:
			close(s.ReqChan)
			log.Println("Terminating SMS Worker ...")
			return
		case req := <-s.ReqChan:
			// process message request and write response to response channel
			req.RspChan <- s.doSend(req.Payload)
		}
	}
}

// doSend performs actual SMS sending
func (s *SMSSender) doSend(bm *model.BaseMessage) Response {
	// create a composite message from provided message text
	cm := util.NewCompositeMessage(bm.Message)

	totalparts := len(cm.MessageParts)
	originator := bm.Originator
	recipients := []string{strconv.Itoa(bm.Recipient)}

	// allow only one call per second towards the external SMS API
	throttle := time.Tick(1 * time.Second)

	// use a WaitGroup to block until all the message parts are sent
	var wg sync.WaitGroup
	for _, mpart := range cm.MessageParts {
		// rate limit SMS sending
		<-throttle

		// increment the waitgroup counter
		wg.Add(1)
		// launch a goroutine to send the message part
		go func(message *util.Message) {
			// decrement the counter when the goroutine completes.
			defer wg.Done()

			var params *messagebird.MessageParams
			// attach UDH parameter only for concatenated sms messages
			if totalparts > 1 {
				params = &messagebird.MessageParams{
					Type: "binary",
					TypeDetails: messagebird.TypeDetails{
						"udh": message.Header.String(),
					},
				}
			}

			// perform actual sending
			response, err := s.SMSClient.NewMessage(originator, recipients, message.Body, params)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("\n%+v\n", response)
		}(mpart)
	}

	// wait for all parts to be sent
	wg.Wait()

	return "Done"
}
