package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ilyahub/backendSandbox/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client started")
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed ClientConnection creation: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewSumServiceClient(cc)
	doUnary(c)
}

func doUnary(c calculatorpb.SumServiceClient) {
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
