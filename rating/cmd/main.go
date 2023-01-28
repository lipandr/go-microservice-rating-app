package main

import (
	"log"
	"net/http"

	"github.com/lipandr/go-microservice-rating-app/rating/internal/controller/rating"
	httpHandler "github.com/lipandr/go-microservice-rating-app/rating/internal/handler/http"
	"github.com/lipandr/go-microservice-rating-app/rating/internal/repository/memory"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httpHandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
