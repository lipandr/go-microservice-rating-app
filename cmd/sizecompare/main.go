package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/lipandr/go-microservice-rating-app/gen"
	"github.com/lipandr/go-microservice-rating-app/metadata/pkg/model"
)

var metadata = &model.Metadata{
	ID:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary The Movie",
	Director:    "Foo Bars",
}

var genMetadata = &gen.Metadata{
	Id:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary The Movie",
	Director:    "Foo Bars",
}

func main() {
	jsonBytes, err := serializeToJSON(metadata)
	if err != nil {
		panic(err)
	}
	xmlBytes, err := serializeToXML(metadata)
	if err != nil {
		panic(err)
	}
	protoBytes, err := serializeToProto(genMetadata)
	if err != nil {
		panic(err)
	}
	fmt.Printf("JSON size:\t%d Bytes\n", len(jsonBytes))
	fmt.Printf("XML size:\t%d Bytes\n", len(xmlBytes))
	fmt.Printf("Proto size:\t%d Bytes\n", len(protoBytes))
}

func serializeToJSON(m *model.Metadata) ([]byte, error) {
	return json.Marshal(m)
}

func serializeToXML(m *model.Metadata) ([]byte, error) {
	return xml.Marshal(m)
}

func serializeToProto(m *gen.Metadata) ([]byte, error) {
	return proto.Marshal(m)
}
