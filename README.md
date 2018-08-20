# golang-p2p-chat

console chat with auto reconnection and TLS encryption written in go

### Installing

```
go get github.com/faroyam/golang-p2p-chat
```

## Getting Started

1. Generate SSL keys for TSL 
```
bash makeCerts.sh
```
2. Run chat
```
go build chat.go && ./chat <username> <remote IP> 
```