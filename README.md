# sss
Simple syslog server

##### Build and run
```
go build .
./sss  --listen-tcp 0.0.0.0 --listen-udp 0.0.0.0 --loki-consumer-url http://192.168.1.150:3100 --loglevel DEBUG
```
