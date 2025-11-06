package server

import (
	"log"
	"net"

	"example.com/user/internal/config"
	"example.com/user/internal/repository"
	"example.com/user/internal/service"
	pb "example.com/user/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server wraps the gRPC server with configuration
type Server struct {
	grpcServer *grpc.Server
	userSvc    *service.UserService
	config     *config.Config
}

// New creates a new gRPC server instance
func New() *Server {
	cfg := config.Load()
	
	// Initialize repository
	userRepo := repository.NewInMemoryUserRepository()
	
	// Initialize service
	userSvc := service.NewUserService(userRepo)
	
	// Create gRPC server with options
	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(cfg.Server.MaxConcurrentStreams),
		grpc.MaxRecvMsgSize(cfg.Server.MaxMessageSize),
		grpc.MaxSendMsgSize(cfg.Server.MaxMessageSize),
	)
	
	// Register services
	pb.RegisterUserServiceServer(grpcServer, userSvc)
	reflection.Register(grpcServer)
	
	return &Server{
		grpcServer: grpcServer,
		userSvc:    userSvc,
		config:     cfg,
	}
}

// Start starts the gRPC server on the configured port
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.config.Server.Port)
	if err != nil {
		return err
	}
	
	log.Printf("🚀 gRPC Server started on %s", s.config.Server.Port)
	log.Printf("📍 Health Check: grpc_health_probe -addr=%s", s.config.Server.Port)
	log.Printf("📍 API Discovery: grpcurl -plaintext %s list", s.config.Server.Port)
	
	return s.grpcServer.Serve(lis)
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	log.Println("🛑 Shutting down gRPC server...")
	s.grpcServer.GracefulStop()
}