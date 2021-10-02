package main

import (
	"context"
	"fmt"
	"io"
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
	callPND(client)
}

func callSum(client calculatorpb.CalculatorServiceClient) {
	resp, err := client.Sum(context.Background(), &calculatorpb.SumRequest{Num1: 10, Num2: 20})
	if err != nil {
		log.Fatalf("error while calling sum api %v", err)
	}

	fmt.Printf("Result is %v\n", resp.GetResult())
}

func callPND(client calculatorpb.CalculatorServiceClient) {
	stream, err := client.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("server finish streaming")
			return
		}

		log.Printf("Receive %v\n", resp.GetResult())
	}
}
