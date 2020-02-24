package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/nitin1259/grpc-go-learn/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	fmt.Printf("invoking the server Greet wiht paramater req: %v \n", req)

	fname := req.GetGreeting().GetFirstName()

	result := "Welcome: " + fname

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("invoking server streaming rpc with req: %v", req)
	firstname := req.GetGreeting().GetFirstName()
	lastname := req.GetGreeting().GetLastName()

	for i := 0; i < 10; i++ {

		result := "Welcome " + firstname + " " + lastname + " number " + strconv.Itoa(i)

		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(res)
	}

	return nil

}

func main() {
	fmt.Println("Welcome to grpc world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v", err)
	}

}
