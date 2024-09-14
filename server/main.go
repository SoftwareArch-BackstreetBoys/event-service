package main

import (
	"fmt"
	"log"
	"net"
	"server/configs"
	"server/services"

	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()

	//run database
	configs.ConnectDB()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	services.RegisterEventServiceServer(s, services.NewEventServiceServer())

	fmt.Println("gRPC server listening on port 50051")
	err = s.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
