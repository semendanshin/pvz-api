package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/swaggest/swgui/v5emb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"homework/internal/abstractions"
	"homework/internal/infrastructure/server/middleware"
	pvzService "homework/internal/infrastructure/server/services/pvz-service"
	desc "homework/pkg/pvz-service/v1"
	"log"
	"net"
	"net/http"
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

func (s *GRPCServer) Run(ctx context.Context, host string, grpcPort, httpPort int) error {
	// Create a new server instance
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.StdLogging,
			middleware.NewErrorMiddleware(),
		),
	)

	// Register the service
	desc.RegisterPvzServiceServer(srv, pvzService.NewPVZService(s.useCase))

	// Reflect the service
	reflection.Register(srv)

	// Create gateway
	gatewayMux := runtime.NewServeMux()
	err := desc.RegisterPvzServiceHandlerFromEndpoint(
		ctx,
		gatewayMux,
		fmt.Sprintf("%s:%d", host, grpcPort),
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	)
	if err != nil {
		return err
	}

	// Create swagger ui
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/swagger", func(w http.ResponseWriter, request *http.Request) {
		http.ServeFile(w, request, "pkg/pvz-service/v1/pvz-service.swagger.json")
	})
	httpMux.Handle("/docs/", v5emb.NewHandler(
		"PVZ Service",
		"/swagger",
		"/docs/",
	))
	httpMux.Handle("/", gatewayMux)

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, httpPort),
		Handler: httpMux,
	}

	// Start the gateway and swagger ui
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

		<-stop

		log.Println("Stopping gracefully...")

		err := httpSrv.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}

		srv.Stop()
	}()

	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, grpcPort))
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
