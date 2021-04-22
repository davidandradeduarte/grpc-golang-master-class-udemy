package main

import (
	"calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with a streaming request\n")
	sum := int32(0)
	count := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: float64(sum) / float64(count),
			})
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request\n")

	var max int32
	firstTime := true
	for {
		send := false
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		number := req.GetNumber()
		if firstTime {
			firstTime = false
			max = number
			send = true
		} else {
			if number > max {
				max = number
				send = true
			}
		}

		if send {
			err = stream.Send(&calculatorpb.FindMaximumResponse{
				Result: max,
			})
			if err != nil {
				log.Fatalf("error while sending data to client: %v", err)
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot function was invoked with %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("received a negative number %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	fmt.Println("serving gRPC at 0.0.0.0:50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
