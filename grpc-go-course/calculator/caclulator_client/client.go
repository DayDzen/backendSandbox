package main

import (
	"context"
	"fmt"
	"io"
	"log"

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
	decompositeNumber(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Srting to do Unary RPC")
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
