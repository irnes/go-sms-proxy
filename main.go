package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/messagebird/go-rest-api"

	"mbsms-api/app"
)

func main() {
	fmt.Println("### SMS API ###")

	// Set up monitoring of operating system signals
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// Create an instnace of messagebird sms client
	//smsClient := messagebird.New("j5ONQuMMG09WNSaFvZawtoWvc")
	smsClient := messagebird.New("test_mCqng0op0JjXkPNe5jEkHZcaO")
	smsSender := app.NewSMSSender(smsClient)

	// Start application using provided SMS sender
	smsApp := app.New(smsSender)
	go smsApp.Run(":8080")

	// Wait for shutdown signal and react
	<-stop

	log.Println("Shutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	smsSender.Terminate()
	//smsApp.Shutdown(ctx)
	_ = ctx

	time.Sleep(1 * time.Second)
	log.Println("Server gracefully stopped")
}

//cudh := [5]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
// udh := NewUserDataHeader()
// fmt.Printf("%x\n", udh)
// fmt.Printf("%x\n", udh)
