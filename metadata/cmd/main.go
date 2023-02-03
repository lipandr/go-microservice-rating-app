package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lipandr/go-microservice-rating-app/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"

	"github.com/lipandr/go-microservice-rating-app/metadata/internal/controller/metadata"
	grpcHandler "github.com/lipandr/go-microservice-rating-app/metadata/internal/handler/grpc"
	"github.com/lipandr/go-microservice-rating-app/metadata/internal/repository/memory"
	"github.com/lipandr/go-microservice-rating-app/pkg/discovery"
	"github.com/lipandr/go-microservice-rating-app/pkg/discovery/consul"
)

const serviceName = "metadata"
const consulService = "localhost:8500"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the metadata service on port %d", port)
	registry, err := consul.NewRegistry(consulService)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	repo := memory.New()
	ctrl := metadata.New(repo)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)

	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
