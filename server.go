package main

import (
	"fmt"
	"github.com/gautamrege/gochat/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type chatServer struct {
}

func (s *chatServer) Chat(ctx context.Context, req *api.ChatRequest) (res *api.ChatResponse, err error) {
	fmt.Printf("\n%s\n> ", fmt.Sprintf("@%s says: \"%s\"", req.From.Name, req.Message))

	// TODO-WORKSHOP-STEP-7: If this is a chat from an unknown user, insert into PeerHandleMap
	_,ok := USERS.Get(req.From.Name)
	if(!ok){

		USERS.Insert(USERS.PeerHandleMap[req.From.Name])
	}
	return &api.ChatResponse{}, nil
}

// gRPC listener - register and start grpc server
func startServer(wg *sync.WaitGroup) {
	defer wg.Done()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", MyHandle.Host, MyHandle.Port))
	if err != nil {
		log.Fatalf("failed to startServer: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterGoChatServer(grpcServer, &chatServer{})

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
