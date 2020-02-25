package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

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
		time.Sleep(1000 * time.Millisecond)
	}

	return nil

}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("invoking client streaming rpc with incoming stream")
	result := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have done with reading client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while processing client streaming rpc: %v", err)
		}

		firstname := req.GetGreeting().GetFirstName()
		result += "Hello " + firstname + " ! \n"
	}

	return nil

}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while receiving clinet stream in server %v \n", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + " !\n"

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if sendErr != nil {
			log.Fatalf("Error while send stream back to client err: %v", sendErr)
			return sendErr
		}
	}

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
