package main

import (
	"log"
	"net"
	"net/http"
	"os"

	githubadapter "github.com/vlm326/golang-course-hw-78/internal/collector/adapter/github"
	collectorgrpc "github.com/vlm326/golang-course-hw-78/internal/collector/handler/grpc"
	"github.com/vlm326/golang-course-hw-78/internal/collector/usecase"
	"github.com/vlm326/golang-course-hw-78/internal/shared/grpcjson"
	"github.com/vlm326/golang-course-hw-78/internal/shared/repositoryrpc"
	"google.golang.org/grpc"
)

func main() {
	address := envOrDefault("COLLECTOR_ADDRESS", ":50051")

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("listen collector: %v", err)
	}

	githubClient := githubadapter.NewClient(&http.Client{})
	getRepository := usecase.NewGetRepositoryUseCase(githubClient)
	server := collectorgrpc.NewServer(getRepository)

	grpcServer := grpc.NewServer(grpc.ForceServerCodec(grpcjson.Codec{}))
	repositoryrpc.RegisterRepositoryServiceServer(grpcServer, server)

	log.Printf("collector is listening on %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("serve collector: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
