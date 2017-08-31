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
	"mbsms-api/app/service"
)

func main() {
	fmt.Println("### SMS API ###")

	// Set up monitoring of operating system signals
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// Create an instnace of messagebird sms client
	//client := messagebird.New("j5ONQuMMG09WNSaFvZawtoWvc")
	client := messagebird.New("test_mCqng0op0JjXkPNe5jEkHZcaO")
	provider := service.NewMBProvider(client)
	service := service.NewSMSService(provider)

	// Start application using provided SMS service
	smsApp := app.New(service)
	go smsApp.Run(":8080")

	// Wait for shutdown signal and react
	<-stop

	log.Println("Shutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	service.Terminate()
	//smsApp.Shutdown(ctx)
	_ = ctx

	time.Sleep(1 * time.Second)
	log.Println("Server gracefully stopped")
}

//cudh := [5]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
// udh := NewUserDataHeader()
// fmt.Printf("%x\n", udh)
// fmt.Printf("%x\n", udh)
