package main

import "fmt"

import "google.golang.org/grpc"

import "log"

import "github.com/nitin1259/grpc-go-learn/greet/greetpb"

import "context"

import "io"

import "time"

func main() {
	fmt.Println("Welcome to grpc go client")

	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("unable to create connection: %v", err)
	}

	defer clientConn.Close()

	client := greetpb.NewGreetServiceClient(clientConn)

	fmt.Printf("Created cleint %f \n", client)

	// doUnary(client)

	// doServerStreamingRPC(client)

	doClientStreamingRPC(client)
}

func doUnary(client greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do Unary RPC \n")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nitin",
			LastName:  "Singh",
		},
	}

	res, err := client.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Greet server method err: %v", err)
	}

	log.Printf("Response from the Greet method res: %v", res)
}

func doServerStreamingRPC(client greetpb.GreetServiceClient) {
	fmt.Println("Starting do do Server Streaming from client")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Kapil",
			LastName:  "Gill",
		},
	}

	resStream, err := client.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling GreetingManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we are done with the stream and reached EOF
		}
		if err != nil {
			log.Fatalf("error while reading stream : %v", err)
		}

		log.Printf("Response from Greet Many Times %v", msg.GetResult())
	}

}

func doClientStreamingRPC(client greetpb.GreetServiceClient) {
	fmt.Println("Starting do do Client Streaming from client")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nitin",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Vipin",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sachin",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Kapil",
			},
		},
	}

	stream, err := client.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling LongGreet rpc: %v", err)
	}

	// we iterate over slice and send each message one by one
	for _, req := range requests {
		fmt.Printf("Sending the request : %v \n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving LongGreet rpc: %v", err)
	}

	log.Printf("Response getting back from Long Greet : %v", res.GetResult())

}
