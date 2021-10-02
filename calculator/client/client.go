package main

import (
	"context"
	"fmt"
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

	resp, err := client.Sum(context.Background(), &calculatorpb.SumRequest{Num1: 10, Num2: 20})
	if err != nil {
		log.Fatalf("error while calling sum api %v", err)
	}

	fmt.Printf("Result is %v", resp.GetResult())
}
