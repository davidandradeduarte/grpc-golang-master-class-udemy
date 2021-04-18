package main

import (
	"calculator/calculatorpb"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	numbers := req.GetNumbers()
	sum := numbers.FirstNumber + numbers.SecondNumber
	res := &calculatorpb.SumResponse{
		Result: sum,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)

	number := req.GetNumber()

	for i := int32(2); number > 1; {
		if number%i == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: i,
			}
			stream.Send(res)
			number = number / i
		} else {
			i = i + 1
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("serving gRPC at 0.0.0.0:50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
