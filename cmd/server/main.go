package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	blockpb "github.com/perlinson/gocraft-server/proto/block"
	playerpb "github.com/perlinson/gocraft-server/proto/player"
	authpb "github.com/perlinson/gocraft-server/proto/auth"
	"github.com/perlinson/gocraft-server/services"
)

var (
	listenAddr = flag.String("l", "0.0.0.0:8421", "listen address")
)

func main() {
	flag.Parse()

	err := InitStore()
	if err != nil {
		log.Fatal(err)
	}
	// 创建gRPC服务器
	grpcServer := grpc.NewServer()

	// 初始化各服务
	blockService := services.NewBlockService()
	playerService := services.NewPlayerService()
	authService := services.NewAuthService()

	// 注册服务
	blockpb.RegisterBlockServiceServer(grpcServer, blockService)
	playerpb.RegisterPlayerServiceServer(grpcServer, playerService)
	authpb.RegisterAuthServiceServer(grpcServer, authService)

	// 启动监听
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	log.Println("Server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
