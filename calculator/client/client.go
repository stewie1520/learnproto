package main

import (
	"log"

	"github.com/stewie1520/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error while dialing %v", err)
	}
	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	log.Printf("service client %f", client)
}
