package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/DayDzen/backendSandbox/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client started")
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed ClientConnection creation: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)
	//decompositeNumber(c)
	// computeAvarage(c)
	// findMaximum(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Unary RPC")
	req := &calculatorpb.SumRequest{
		FirstNumber:  32,
		SecondNumber: 54,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed Sum: %v", err)
	}
	log.Printf("Result is: %v", res.SumResult)
}

func decompositeNumber(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Decomposition number from client")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 240,
	}
	decClient, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed PrimeNumberDecomposition RPC: %v", err)
	}
	log.Printf("%v decomposing on: ", req.Number)
	for {
		decResp, err := decClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed while getting responses from PrimeNumberDecomposition: %v", err)
		}
		log.Printf("Response from PrimeNumberDecomposition: %v", decResp.GetPrimeNumber())
	}
}

func computeAvarage(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Computing avarage of numbers from client")
	caClient, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Failed client creating in computeAvarage function: %v", err)
	}
	numbers := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 32,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 42,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 52,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 48,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 64,
		},
	}
	for _, number := range numbers {
		err := caClient.Send(number)
		time.Sleep(500 * time.Millisecond)
		log.Printf("Sended number: %v\n", number.GetNumber())
		if err != nil {
			log.Fatalf("Failed while sending numbers: %v", err)
		}
	}
	res, err := caClient.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed getting response from server: %v", err)
	}
	log.Printf("Avarage of numbers is %v", res.GetAvarage())
}

func findMaximum(c calculatorpb.CalculatorServiceClient) {
	numbers := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number: 3,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 23,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 4,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 32,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 51,
		},
		&calculatorpb.FindMaximumRequest{
			Number: -1,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 60,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 4,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 100,
		},
	}

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Failed while creating stream for client: %v", err)
	}
	waitc := make(chan struct{})

	go func() {
		for _, number := range numbers {
			err := stream.Send(number)
			log.Printf("Sended number: %v", number.GetNumber())
			time.Sleep(500 * time.Millisecond)
			if err != nil {
				log.Fatalf("Failed while sending number to Server from Client: %v", err)
				return
			}
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("Failed while closing stream for sending numuber: %v", err)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed while getting response from Server: %v", err)
				break
			}
			log.Printf("NEW MAXIMUM IS: %v", res.GetMaxNumber())
		}
		close(waitc)
	}()

	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Error Unary RPC")
	req := &calculatorpb.SquareRootRequest{
		Number: 11,
	}
	res, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We sent negative number")
			}
		} else {
			log.Fatalf("Big Errror calling SquareRoot: %v", err)
		}
	}
	fmt.Printf("Result of %v is: %v\n", req.GetNumber(), res.GetNumberRoot())
}
