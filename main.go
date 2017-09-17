package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	// Create an instnace of messagebird sms client
	client := messagebird.New(conf.apikey)
	smsProvider := service.NewMBProvider(ctx, client)
	smsService := service.NewSMSService(smsProvider)

	// Print account balance information
	smsService.Balance()

	// Start application using provided SMS service
	// and listen on all host ifaces
	smsApp := app.New(smsService)
	go smsApp.Run("127.0.0.1" + ":" + strconv.Itoa(conf.port))

	// Wait for shutdown signal and react
	<-stop
	// stop http server
	go smsApp.Shutdown(ctx)

	log.Println("Shutting down the server...")
	cancel()

	time.Sleep(1 * time.Second)
	log.Println("Server gracefully stopped")
}
