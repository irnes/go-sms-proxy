## SMS REST API Proxy

Demonstrate a way how to implement a REST API proxy service that provide an
unique interface for sending SMS messages through various SMS gateways

The current version relies on usage of MessageBird SMS Gateway

### Help
```
Usage of ./go-sms-proxy
  -apikey string
    	API Key for SMS Gateway (MessageBird)
  -port int
    	Port to listen (default 8080)
```

### Usage:
Start REST API proxy service on the local host
```
./go-sms-proxy -apikey=<apikey>
```

and use provided helper script to send a message through REST API service
```
./test/send-sms.sh +387######## "This is a test message."
