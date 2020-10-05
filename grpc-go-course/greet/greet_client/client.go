package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/DayDzen/backendSandbox/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, I'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)
	// doUnary(c)

	// doServerStreaming(c)

	doClientStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ilya",
			LastName:  "Izmailov",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet RPC: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ilya",
			LastName:  "Izmailov",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we have reached end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading the stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a Client Streaming RPC...")

	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ilya",
				LastName:  "Izmailov",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lidiya",
				LastName:  "Izmailova",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "James",
				LastName:  "Bond",
			},
		},
	}

	lgClient, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Failed getting LongGreetClient: %v", err)
	}

	for _, req := range request {
		err := lgClient.Send(req)
		log.Printf("Sended %v", req.GetGreeting().GetFirstName())
		time.Sleep(2000 * time.Millisecond)
		if err != nil {
			log.Fatalf("Failed sending request: %v", err)
		}
	}

	lgResp, err := lgClient.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed getting response: %v", err)
	}

	log.Printf(lgResp.GetResult())
}
