package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lipandr/go-microservice-rating-app/movie/internal/controller/movie"
	metadataGateway "github.com/lipandr/go-microservice-rating-app/movie/internal/gateway/metadata/http"
	ratingGateway "github.com/lipandr/go-microservice-rating-app/movie/internal/gateway/rating/http"
	httpHandler "github.com/lipandr/go-microservice-rating-app/movie/internal/handler/http"
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
	h := httpHandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
