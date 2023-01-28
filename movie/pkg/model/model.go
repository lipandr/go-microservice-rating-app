package model

import "github.com/lipandr/go-microservice-rating-app/metadata/pkg/model"

// MovieDetails includes movie metadata its aggregated rating.
type MovieDetails struct {
	Rating   float64        `json:"rating,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}
