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

// NewMBProvider creates an instance of MBProvider object using provider's client
func NewMBProvider(client *messagebird.Client) *MBProvider {
	provider := &MBProvider{
		Client:  client,
		Done:    make(chan bool, 0),
		ReqChan: make(chan Request, 10),
	}

	// start provider worker
	go provider.run()

	return provider
}

// MBProvider struct implements SMSProvider interface
// that uses MessageBird client as an SMS gateway
type MBProvider struct {
	Client  *messagebird.Client
	Done    chan bool
	ReqChan chan Request
}

// Run will make the MBProvider starts listening to the two channels Done and ReqChan
func (m *MBProvider) run() {
	for {
		select {
		case <-m.Done:
			close(m.ReqChan)
			log.Println("Terminating SMS Worker ...")
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
			response, err := m.Client.NewMessage(originator, recipients, message.Body, params)
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

// Terminate SMS worker
func (m *MBProvider) Terminate() {
	m.Done <- true
}
