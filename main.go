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

	ctx, cancel := context.WithCancel(context.Background())
	// Set up monitoring of operating system signals
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// Create an instnace of messagebird sms client
	client := messagebird.New(conf.apikey)
	smsProvider := service.NewMBProvider(ctx, client)
	smsService := service.NewSMSService(smsProvider)

	// Print account balance information
	smsService.Balance()

	// Start application using provided SMS service
	// and listen on all host ifaces
	smsApp := app.New(smsService)
	go smsApp.Run("0.0.0.0" + ":" + strconv.Itoa(conf.port))

	// Wait for shutdown signal and react
	<-stop

	log.Println("Shutting down the server...")

	// stop http server
	lastctx, _ := context.WithTimeout(ctx, 1*time.Second)
	smsApp.Shutdown(lastctx)

	cancel()
	time.Sleep(1 * time.Second)

	log.Println("Server gracefully stopped")
}
