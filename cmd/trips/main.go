package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	apiv1 "github.com/sharifahmad2061/trip-grpc-go/api/gen/go"
	"github.com/sharifahmad2061/trip-grpc-go/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	//conf := config.Load()
	//fmt.Println(conf)
	socket, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	apiv1.RegisterTripsServer(server, &service.TripsServiceImpl{})
	reflection.Register(server)

	go func() {
		if err := server.Serve(socket); err != nil {
			log.Fatalf("failed to serve: %v", err)
		} else {
			log.Println("gRPC server is running on port 50051")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	server.GracefulStop()
	log.Println("Server stopped")
}
