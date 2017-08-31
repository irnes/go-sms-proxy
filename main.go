package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/messagebird/go-rest-api"

	"mbsms-api/app"
	"mbsms-api/app/service"
)

// Config holds basic app config
type Config struct {
	apikey string
	port   int
}

var conf Config

func init() {
	flag.StringVar(&conf.apikey, "apikey", "test_mCqng0op0JjXkPNe5jEkHZcaO", "API Key for MessageBird")
	flag.IntVar(&conf.port, "port", 8080, "Port to listen")

	flag.Parse()
}

func main() {
	fmt.Println("### SMS REST API ###")

	// Set up monitoring of operating system signals
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// Create an instnace of messagebird sms client
	//client := messagebird.New("j5ONQuMMG09WNSaFvZawtoWvc")
	client := messagebird.New(conf.apikey)
	smsProvider := service.NewMBProvider(client)
	smsService := service.NewSMSService(smsProvider)

	// Print account balance information
	smsService.Balance()

	// Start application using provided SMS service
	smsApp := app.New(smsService)
	go smsApp.Run(":" + strconv.Itoa(conf.port)) // listen on all host ifaces

	// Wait for shutdown signal and react
	<-stop

	log.Println("Shutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	smsService.Terminate()
	smsApp.Shutdown(ctx)
	_ = ctx

	log.Println("Server gracefully stopped")
}
