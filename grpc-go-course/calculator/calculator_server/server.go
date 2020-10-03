package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

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

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("The PrimeNumberDecomposition function was invoked with %v\n", req)
	var k int32
	k = 2
	number := req.GetNumber()
	for number > 1 {
		if number%k == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeNumber: k,
			}
			stream.Send(res)
			number = number / k

			time.Sleep(500 * time.Millisecond)
		} else {
			k++
		}
	}
	return nil
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
