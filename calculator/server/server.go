package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/stewie1520/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Max(stream calculatorpb.CalculatorService_MaxServer) error {
	var max int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error while receiving client max number %v\n", err)
			return err
		}

		num := req.GetNumber()
		if num > max {
			max = num
		}

		err = stream.Send(&calculatorpb.MaxResponse{Result: max})
		if err != nil {
			log.Fatalf("error while send response %v\n", err)
			return err
		}
	}
}

func (s *server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	var total float32
	var count int32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}

			stream.SendAndClose(resp)
			return nil
		}

		if err != nil {
			log.Fatalf("error while receiving client number %v", err)
		}

		log.Printf("receive %v\n", req.GetNumber())
		count++
		total += req.GetNumber()
	}
}

//Sum return sum of 2 numbers
func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PNDRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {

	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			stream.Send(&calculatorpb.PNDResponse{Result: k})
			time.Sleep(500 * time.Millisecond)
		} else {
			k++
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		log.Fatalf("error while creating lisen %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("calculator is running...")
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("error while serving %v", err)
	}
}
