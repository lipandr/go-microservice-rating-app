package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/lipandr/go-microservice-rating-app/gen"
	"github.com/lipandr/go-microservice-rating-app/movie/internal/controller/movie"
	metadataGateway "github.com/lipandr/go-microservice-rating-app/movie/internal/gateway/metadata/grpc"
	ratingGateway "github.com/lipandr/go-microservice-rating-app/movie/internal/gateway/rating/grpc"
	grpcHandler "github.com/lipandr/go-microservice-rating-app/movie/internal/handler/grpc"
	"github.com/lipandr/go-microservice-rating-app/pkg/discovery"
	"github.com/lipandr/go-microservice-rating-app/pkg/discovery/consul"
)

const serviceName = "movie"
const consulService = "localhost:8500"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
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

	metadataGW := metadataGateway.New(registry)
	ratingGW := ratingGateway.New(registry)
	ctrl := movie.New(ratingGW, metadataGW)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to liten: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
