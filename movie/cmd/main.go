package main

import (
	"log"
	"net/http"

	"github.com/lipandr/go-microservice-rating-app/movie/internal/controller/movie"
	metadatagateway "github.com/lipandr/go-microservice-rating-app/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/lipandr/go-microservice-rating-app/movie/internal/gateway/rating/http"
	httpHandler "github.com/lipandr/go-microservice-rating-app/movie/internal/handler/http"
)

func main() {
	log.Println("Starting the movie service")
	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httpHandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
