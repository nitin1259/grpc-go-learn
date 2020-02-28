package main

import "fmt"

import "google.golang.org/grpc"

import "log"

import "github.com/nitin1259/grpc-go-learn/greet/greetpb"

import "context"

import "io"

import "time"

import "google.golang.org/grpc/status"

import "google.golang.org/grpc/codes"

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

	// doClientStreamingRPC(client)

	// doBiDiStreamingRPC(client)

	doUnaryWithDeadLines(client)
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

func doBiDiStreamingRPC(client greetpb.GreetServiceClient) {
	fmt.Println("Starting do do BiDi Streaming from client")

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nitin",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Vipin",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Sachin",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Kapil",
			},
		},
	}

	// create the stream by invoking the client
	stream, err := client.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while calling BiDi stream rpc: %v \n", err)
	}

	waitc := make(chan struct{})

	// send bunch of message to the server (go routine)
	go func() {
		// method to send lot of message
		for _, req := range requests {
			fmt.Printf("Sending request : %v \n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()

	}()

	// recieve bunch of message from the server (go routine)
	go func() {

		for {

			res, err := stream.Recv()
			if err == io.EOF {
				// done with receiving response in BiDi
				close(waitc)
			}
			if err != nil {
				log.Fatalf("Error while receiveing response from BiDi : %v \n", err)
				close(waitc)
			}

			fmt.Printf("Recieving response for BiDi : %v", res.GetResult())
		}

	}()

	// block until everything is done
	<-waitc

}

func doUnaryWithDeadLines(client greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do Unary With Deadlines RPC \n")

	doDeadlinesUnary(client, 5*time.Second)

	doDeadlinesUnary(client, 1*time.Second)

}

func doDeadlinesUnary(client greetpb.GreetServiceClient, timeout time.Duration) {

	req := &greetpb.GreetWithDeadlinesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Nitin",
			LastName:  "Singh",
		},
	}

	ctx := context.Background()
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := client.GreetWithDeadlines(c, req)

	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			if resErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Time out err! request deadline exceeded")
				fmt.Println(resErr.Message())
			} else {
				log.Fatalf("Some Error while calling dealines server, %v", err)
			}
		}
		log.Fatalf("Error while calling server GreetWithDeadlines err: %v", err)
		return
	}

	log.Printf("Response from the GreetWithDeadlines method res: %v", res)
}
