package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/DayDzen/backendSandbox/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("The Greet function was invoked with %v\n", req)
	fn := req.GetGreeting().GetFirstName()
	result := "Hello " + fn
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("The GreetManyTimes function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		err := stream.Send(res)
		if err != nil {
			log.Fatalf("Failed while sending response from GreetManyTimes: %v", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("The LongGreet function was invoked with steaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//end of streaming
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Failed stream reading: %v", err)
		}
		fName := req.GetGreeting().GetFirstName()
		result += "\nHello " + fName + "!!!"
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("The GreetEveryone function was invoked with steaming request\n")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Failed while getting streaming request in GreetEveryone: %v", err)
			return err
		}
		fName := req.GetGreeting().GetFirstName()
		result := "Hello " + fName + "!!!\n"
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Failed while sending response: %v", err)
			return err
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("The GreetWithDeadline function was invoked with %v\n", req)
	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("the client canceled the request")
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		default:
			time.Sleep(1 * time.Second)
			log.Printf("Pinged: %v\n", i)
		}
		//return nil, status.Error(codes.Internal, "pizda proizoshla")
	}
	fn := req.GetGreeting().GetFirstName()
	result := "Hello " + fn
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	log.Println(res)
	return res, nil
}

func main() {
	fmt.Println("Hello there")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
