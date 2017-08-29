curl -i -H 'Content-Type: application/json'  \
-d '{"recipient": 38761475148, "originator":"IRNES-API", "message":"Hello from SMS API!"}' \
http://127.0.0.1:8080/messages
