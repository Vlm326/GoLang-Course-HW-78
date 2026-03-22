package main

import (
	"log"
	"net/http"
	"os"
	"time"

	httpHandler "github.com/vlm326/golang-course-hw-78/internal/gateway/handler/http"
	gatewayusecase "github.com/vlm326/golang-course-hw-78/internal/gateway/usecase"
	"github.com/vlm326/golang-course-hw-78/internal/shared/grpcjson"
	"github.com/vlm326/golang-course-hw-78/internal/shared/repositoryrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	httpAddress := envOrDefault("API_GATEWAY_ADDRESS", ":8080")
	collectorAddress := envOrDefault("COLLECTOR_GRPC_ADDRESS", "localhost:50051")

	conn, err := grpc.NewClient(
		collectorAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(grpcjson.Codec{})),
	)
	if err != nil {
		log.Fatalf("dial collector: %v", err)
	}
	defer conn.Close()

	client := repositoryrpc.NewRepositoryServiceClient(conn)
	getRepository := gatewayusecase.NewGetRepositoryUseCase(client)
	handler := httpHandler.NewHandler(getRepository)

	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("api gateway is listening on %s", httpAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve api gateway: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
