curl -i -H 'Content-Type: application/json'  \
-d '{"recipient": 38761475148, "originator":"IRNES-API", "message":"Hello from SMS API!"}' \
http://127.0.0.1:8080/messages

#curl -i -H 'Content-Type: application/json'  \
#-d '{"recipient": 38761475148, "originator":"IRNES-API", "message":"With MessageBird, we send 1M messages in 10 minutes. Every second counts, MessageBird is literally saving lives. 0123456789 abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ~!@#$%^&*()-_+={}[]\\|<,>.?/\";:"}' \
#http://127.0.0.1:8080/messages


