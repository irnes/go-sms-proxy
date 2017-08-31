## MessageBird SMS REST API Client


### Help
```
Usage of ./mbsms-api:
  -apikey string
    	API Key for MessageBird (default "test_mCqng0op0JjXkPNe5jEkHZcaO")
  -port int
    	Port to listen (default 8080)
```

### Usage:
Start REST API service on the local host
```
./mbsms-api -apikey=<apikey>
```

and use provided helper script to send a message through REST API service
```
./test/send-sms.sh +38761475148 "This is a test message."
