# Live Chat App

Real-time chat application built with pure Go and gRPC. It allows users to join chat rooms, send messages, and leave the chat in a synchronous manner.


## Getting Started

### Server
```
go run app/cmd/server/main.go
2023/12/25 23:52:29 gRPC Server Listing on: localhost:50051...
```

### Client

```
go run app/cmd/client/main.go
Enter username: doe
Enter chat room: manga
[server] : 'doe' join the chat
[doe] : hello, world
[jack] : hi doe
...
[doe] : /exit
[server] doe left the chat
```


## Reference
- [gRPC](https://grpc.io/)
- [Buf](https://buf.build/)
- [Connect RPC](https://connectrpc.com/)
