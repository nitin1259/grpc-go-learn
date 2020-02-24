package main

import "fmt"

import "google.golang.org/grpc"

import "log"
import "github.com/nitin1259/grpc-go-learn/calculator/calcpb"

import "context"

import "io"

import "time"

func main() {
	fmt.Println("Starting the client main")

	clientconn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error while creating client connection %v", err)
	}

	client := calcpb.NewCalcServiceClient(clientconn)

	// doUnaryRPC(client)

	// doServerStreamPrimeDecompsitonRPC(client)

	doClientStreamComputeAverageRPC(client)

}

func doUnaryRPC(client calcpb.CalcServiceClient) {
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

func doServerStreamPrimeDecompsitonRPC(client calcpb.CalcServiceClient) {
	fmt.Println("client request for prime with stream rpc")

	req := &calcpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	resStream, err := client.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while doing steaming rpc call from client: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			log.Println("reached End of stream")
			break
		}

		log.Printf("Response from the calc server %v", msg.GetResult())

	}

}

func doClientStreamComputeAverageRPC(client calcpb.CalcServiceClient) {
	fmt.Println("client streaming request for Computer Average rpc")

	requests := []*calcpb.ComputeAverageRequest{
		&calcpb.ComputeAverageRequest{
			Number: 20,
		},
		&calcpb.ComputeAverageRequest{
			Number: 21,
		},
		&calcpb.ComputeAverageRequest{
			Number: 22,
		},
		&calcpb.ComputeAverageRequest{
			Number: 23,
		},
		&calcpb.ComputeAverageRequest{
			Number: 24,
		},
		&calcpb.ComputeAverageRequest{
			Number: 25,
		},
	}

	stream, err := client.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error while doing Client steaming rpc call : %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request from client stream: %v \n", req)
		stream.Send(req)

		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receving response from server: %v", err)
	}

	log.Printf("Response from client Compute Average res: %v", res.GetResult())

}
