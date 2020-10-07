package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

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

	// doClientStreaming(c)

	// doBiDiStreaming(c)

	// doUnaryCallWithDeadline(c, 5*time.Second) // should complete
	// doUnaryCallWithDeadline(c, 3*time.Second) // should complete
	doUnaryCallWithDeadline(c, 1*time.Second) //should timeout
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

	log.Print(lgResp.GetResult())
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting a BiDi Streaming RPC...")
	request := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ilya",
				LastName:  "Izmailov",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lidiya",
				LastName:  "Izmailova",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "James",
				LastName:  "Bond",
			},
		},
	}
	geClient, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Failed while getting client from GreetEveryone: %v", err)
		return
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range request {
			log.Printf("Sending %v", req.GetGreeting().GetFirstName())
			err := geClient.Send(req)
			time.Sleep(1000 * time.Millisecond)
			if err != nil {
				log.Fatalf("Failed sending request: %v", err)
				return
			}
		}
		err := geClient.CloseSend()
		if err != nil {
			log.Fatalf("Failed closing sending channel from client: %v", err)
			return
		}
	}()

	go func() {
		for {
			res, err := geClient.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed while getting response from server: %v", err)
				break
			}
			fmt.Printf("Recieved: %v\n", res.GetResult())
		}
		close(waitc)
	}()
	<-waitc
}

func doUnaryCallWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a UnaryWithDeadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ilya",
			LastName:  "Izmailov",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("Unexpected error!: %v\n", statusErr)
			}
		} else {
			log.Fatalf("Error while calling GreetWtihDeadline RPC: %v", err)
		}
		return
	}
	log.Printf("Response from GreetWtihDeadline RPC: %v", res.Result)
}
