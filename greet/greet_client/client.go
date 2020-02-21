package main

import "fmt"

import "google.golang.org/grpc"

import "log"

import "github.com/nitin1259/grpc-go-learn/greet/greetpb"

import "context"

func main() {
	fmt.Println("Welcome to grpc go client")

	clientConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("unable to create connection: %v", err)
	}

	defer clientConn.Close()

	client := greetpb.NewGreetServiceClient(clientConn)

	fmt.Printf("Created cleint %f \n", client)

	doUnary(client)
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
