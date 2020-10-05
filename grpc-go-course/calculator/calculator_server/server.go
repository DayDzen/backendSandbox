package main

import (
	"context"
	"fmt"
	"io"
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
	k := int32(2)
	number := req.GetNumber()
	for number > 1 {
		if number%k == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeNumber: k,
			}
			err := stream.Send(res)
			if err != nil {
				log.Fatalf("Failed while sending response from PrimeNumberDecomposition: %v", err)
			}
			number = number / k
			time.Sleep(500 * time.Millisecond)
		} else {
			k++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("The ComputeAverage function was invoked with stream of request\n")
	sum := 0.0
	i := 0.0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			avarage := sum / i
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Avarage: avarage,
			})
		}
		if err != nil {
			log.Fatalf("Failed while getting stream request: %v", err)
		}
		sum += req.GetNumber()
		i++
	}
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
