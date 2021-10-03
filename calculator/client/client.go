package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	callMax(client)
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

func callAverate(client calculatorpb.CalculatorServiceClient) {
	stream, err := client.Average(context.Background())
	if err != nil {
		log.Fatalf("error while calling average %v", err)
	}

	numbers := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, number := range numbers {
		err := stream.Send(&calculatorpb.AverageRequest{Number: number})
		time.Sleep(500 * time.Millisecond)
		if err != nil {
			log.Fatalf("error while sending average request %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving avarage response %v", err)
	}

	log.Printf("result %v\n", resp.GetResult())
}

func callMax(client calculatorpb.CalculatorServiceClient) {
	stream, err := client.Max(context.Background())
	if err != nil {
		log.Fatalf("error while calling max %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		numbers := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9}
		for _, number := range numbers {
			err := stream.Send(&calculatorpb.MaxRequest{Number: number})
			if err != nil {
				log.Fatalf("error while sending max request %v", err)
			}
			time.Sleep(500 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}

			if err != nil {
				log.Fatalf("error while receiving max response %v\n", err)
				break
			}

			log.Printf("result of max is %v\n\n", resp.GetResult())
		}
	}()

	<-waitc
}
