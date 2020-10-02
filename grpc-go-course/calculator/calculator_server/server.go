package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/DayDzen/backendSandbox/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("The Sum function was invoked with %v\n", req)
	fn := req.GetFirstNumber()
	sn := req.GetSecondNumber()
	sum := fn + sn
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}
	return res, nil
}

func main() {
	fmt.Println("Server started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Listener server problems: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed Serve: %v", err)
	}
}
