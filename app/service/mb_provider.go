package service

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/messagebird/go-rest-api"

	"go-sms-proxy/app/model"
	"go-sms-proxy/app/util"
)

// NewMBProvider creates an instance of MBProvider object using provider's client
func NewMBProvider(ctx context.Context, client *messagebird.Client) *MBProvider {
	provider := &MBProvider{
		Client:  client,
		ReqChan: make(chan Request, 10),
	}

	// start provider worker
	go provider.run(ctx)

	return provider
}

// MBProvider struct implements SMSProvider interface
// that uses MessageBird client as an SMS gateway
type MBProvider struct {
	Client  *messagebird.Client
	ReqChan chan Request
}

// Run will make the MBProvider starts listening to the two channels Done and ReqChan
func (m *MBProvider) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(m.ReqChan)
			log.Println("Terminating SMS Provider ...")
			return
		case req := <-m.ReqChan:
			// process message request and write response to response channel
			req.RspChan <- m.doSend(req.Payload)
		}
	}
}

// doSend performs actual SMS sending
func (m *MBProvider) doSend(bm *model.BaseMessage) Response {
	// create a composite message from provided message text
	cm := util.NewCompositeMessage(bm.Message)

	totalparts := len(cm.MessageParts)
	originator := bm.Originator
	recipients := []string{bm.Recipient}

	// allow only one call per second towards the external SMS API
	throttle := time.Tick(1 * time.Second)

	// total number of sent message parts
	var totalSent int32

	// use a WaitGroup to block until all the message parts are sent
	var wg sync.WaitGroup
	for _, mpart := range cm.MessageParts {
		// rate limit SMS sending
		<-throttle

		// increment the waitgroup counter
		wg.Add(1)
		// launch a goroutine to send the message part
		go func(message *util.Message) {
			// decrement the counter when the goroutine completem.
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
			_, err := m.Client.NewMessage(originator, recipients, message.Body, params)
			if err != nil {
				// TODO collect errors and insert them into the response
				log.Print("Error: ", err)
			}

			atomic.AddInt32(&totalSent, 1)
		}(mpart)
	}

	// wait for all parts to be sent
	wg.Wait()

	if int(totalSent) != len(cm.MessageParts) {
		return map[string]interface{}{
			"Status": "Failed",
		}
	}

	return map[string]interface{}{
		"Status":         "Success",
		"TotalSentParts": totalSent,
	}
}

// Balance returns the account balance information
func (m *MBProvider) Balance() {
	// Request the balance information, returned as a Balance object.
	balance, err := m.Client.Balance()
	if err != nil {
		// messagebird.ErrResponse means custom JSON errors.
		if err == messagebird.ErrResponse {
			for _, mbError := range balance.Errors {
				log.Printf("Error: %#v\n", mbError)
			}
		}
		return
	}
	log.Printf("Balance: %+v\n", balance)
}

// Send encapsulates given base message into the Request
// structure and enqueues it for further processing
// returns reponse channel
func (m *MBProvider) Send(bm *model.BaseMessage) <-chan Response {
	req := Request{
		Payload: bm,
		RspChan: make(chan Response, 1),
	}
	// enqueue request for processing
	m.ReqChan <- req
	return req.RspChan
}
