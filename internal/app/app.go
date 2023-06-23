package app

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/begenov/backend/internal/config"
	"github.com/begenov/backend/internal/delivery/gapi"
	"github.com/begenov/backend/internal/repository"
	"github.com/begenov/backend/internal/service"
	"github.com/begenov/backend/pb"
	"github.com/begenov/backend/pkg/auth"
	"github.com/begenov/backend/pkg/db"
	"github.com/begenov/backend/pkg/hash"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func Run(cfg *config.Config) error {
	db, err := db.NewDB(cfg.Postgres.Driver, cfg.Postgres.DSN)
	if err != nil {
		return err
	}

	hash := hash.NewHash()

	token, err := auth.NewManager(cfg.JWT.TokenSymmetricKey)
	if err != nil {
		return err
	}

	repo := repository.NewRepository(db)

	service := service.NewService(repo, hash, token, cfg.JWT.AccessTokenDuration)

	// handler := delivery.NewHandler(service, token)

	// srv := server.NewServer(cfg, handler.Init(cfg))
	go runGatewayServer(cfg, service, token)
	runGrpcServer(cfg, service, token)
	// go func() {
	// 	if err = srv.Run(); err != nil {
	// 		log.Fatalf("error occurred while running http server: %s\n", err.Error())
	// 	}
	// }()

	// log.Println("Server started")

	// quit := make(chan os.Signal, 1)

	// signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// <-quit

	// const timeout = 5 * time.Second

	// ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	// defer shutdown()

	// if err := srv.Stop(ctx); err != nil {
	// 	log.Printf("failed to stop server: %v", err)
	// }

	return nil
}

func runGrpcServer(cfg *config.Config, service *service.Service, token auth.TokenManager) {
	server := gapi.NewHandler(service, token)

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "localhost:"+cfg.Server.GrpcAddr)
	if err != nil {
		log.Fatal("cannot create listener", err)
	}
	log.Println("starting grpc server:" + cfg.Server.GrpcAddr)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server", err)
	}
}

func runGatewayServer(cfg *config.Config, service *service.Service, token auth.TokenManager) {
	server := gapi.NewHandler(service, token)

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("error ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", "localhost:"+cfg.Server.Addr)
	if err != nil {
		log.Fatal("cannot create listener", err)
	}

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal(err)
	}

}
