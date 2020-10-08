package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"golang.org/x/net/context"

	"github.com/DayDzen/backendSandbox/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	createBlog(c)
}

func createBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Creating blog")

	blog := &blogpb.Blog{
		AuthorId: "Ilya",
		Title:    "My First Blog",
		Content:  "Some content for first blog",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		status.Errorf(codes.Internal, fmt.Sprintf("Failed creating blog: %v\n", err))
	}

	fmt.Printf("Blog that we created: %v", res.GetBlog())
}
