package main

import (
	"flag"
	"log"
	"net"

	authpb "github.com/perlinson/gocraft-server/internal/proto/auth"
	blockpb "github.com/perlinson/gocraft-server/internal/proto/block"
	playerpb "github.com/perlinson/gocraft-server/internal/proto/player"
	Store "github.com/perlinson/gocraft-server/internal/store"

	"github.com/perlinson/gocraft-server/internal/services"
	"google.golang.org/grpc"
)

// var (
// 	listenAddr = flag.String("l", "0.0.0.0:8421", "listen address")
// )

func main() {
	flag.Parse()

	// 注释掉InitStore调用，根据用户要求不需要关心这部分
	store, err := Store.InitStore()
	if err != nil {
		log.Fatal(err)
	}
	// 创建gRPC服务器
	grpcServer := grpc.NewServer()

	// 初始化各服务
	blockService := services.NewBlockService(store)
	playerService := services.NewPlayerService(nil) // 暂时传入nil
	authService := services.NewAuthService(nil)     // 暂时传入nil

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
