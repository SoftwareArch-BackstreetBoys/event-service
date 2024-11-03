package main

import (
	"fmt"
	"log"
	"net"
	"server/configs"
	"server/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc/reflection"
)

func main() {
	s := grpc.NewServer()

	// Register the health check service.
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	// Set the server's health status to SERVING for the main service and "grpc.health.v1.Health".
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("EventService", grpc_health_v1.HealthCheckResponse_SERVING)

	//run database
	configs.ConnectDB()

	// Enable reflection for grpcurl and other tools to access service descriptors
	reflection.Register(s)

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
