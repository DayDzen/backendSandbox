package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/DayDzen/backendSandbox/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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

func sendMaxNumber(stream calculatorpb.CalculatorService_FindMaximumServer, max int64) {
	err := stream.Send(&calculatorpb.FindMaximumResponse{
		MaxNumber: max,
	})
	if err != nil {
		log.Fatalf("Failed while sending response to Client: %v", err)
		return
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("The FindMaximum function was invoked with stream of request\n")
	var max int64
	fTime := true
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Failed while getting stream request from Client: %v", err)
			return nil
		}
		number := req.GetNumber()

		if fTime {
			max = number
			fTime = false
			sendMaxNumber(stream, max)
		} else {
			if max < number {
				max = number
				sendMaxNumber(stream, max)
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("The SquareRoot function was invoked with request\n")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Recieved a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
