package main

import (
	"context"
	"fmt"
	"greet/greetpb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	greet := fmt.Sprintf("Hello %s", firstName)
	res := &greetpb.GreetResponse{
		Result: greet,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManytimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)
	for i := 0; i < 10; i++ {
		firstName := req.GetGreeting().GetFirstName()
		greet := fmt.Sprintf("Hello %s #%d", firstName, i)
		res := &greetpb.GreetManyTimesResponse{
			Result: greet,
		}
		stream.Send(res)
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
	greetpb.RegisterGreetServiceServer(s, &server{})

	fmt.Println("serving gRPC at 0.0.0.0:50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
