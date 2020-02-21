package main

import "fmt"

import "google.golang.org/grpc"

import "log"
import "github.com/nitin1259/grpc-go-learn/calculator/calcpb"

import "context"

func main() {
	fmt.Println("Starting the client main")

	clientconn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error while creating client connection %v", err)
	}

	client := calcpb.NewCalcServiceClient(clientconn)

	fmt.Printf("client created with calc pb: %v \n", client)

	req := &calcpb.CalcRequest{
		Calcate: &calcpb.Calc{
			Num1: 10,
			Num2: 3,
		},
	}

	res, err := client.Calculator(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while connection server and caluclate: %v", err)
	}

	log.Printf("Response from the calc server %v", res)

}
