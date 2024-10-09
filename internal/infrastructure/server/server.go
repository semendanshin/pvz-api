package server

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"homework/internal/abstractions"
	pvz_service "homework/internal/infrastructure/server/services/pvz-service"
	desc "homework/pkg/pvz-service/v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type GRPCServer struct {
	useCase abstractions.IPVZOrderUseCase
}

func NewGRPCServer(useCase abstractions.IPVZOrderUseCase) *GRPCServer {
	return &GRPCServer{
		useCase: useCase,
	}
}

func (s *GRPCServer) Run(host string, port int) error {

	// Create a new server instance
	srv := grpc.NewServer()

	// Register the service
	desc.RegisterPvzServiceServer(srv, pvz_service.NewPVZService(s.useCase))

	// Reflect the service
	reflection.Register(srv)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

		<-stop

		log.Println("Stopping gracefully...")
		srv.Stop()
	}()

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	log.Println("Starting server...")

	// Start the server
	if err := srv.Serve(lis); err != nil {
		return err
	}

	return nil
}
