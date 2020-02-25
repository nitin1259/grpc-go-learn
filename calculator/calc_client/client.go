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

	// doClientStreamComputeAverageRPC(client)

	doBiDiFindMaxRPC(client)

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

func doBiDiFindMaxRPC(client calcpb.CalcServiceClient) {
	fmt.Println("BiDi streaming request for Find Max rpc")

	// get the client stream for request
	stream, err := client.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while getting client stream for Find Maximun: %v \n", err)
	}

	requests := []*calcpb.FindMaximumRequest{
		&calcpb.FindMaximumRequest{
			Number: 5,
		},
		&calcpb.FindMaximumRequest{
			Number: 3,
		},
		&calcpb.FindMaximumRequest{
			Number: 6,
		},
		&calcpb.FindMaximumRequest{
			Number: 1,
		}, &calcpb.FindMaximumRequest{
			Number: 9,
		},
		&calcpb.FindMaximumRequest{
			Number: 18,
		},
		&calcpb.FindMaximumRequest{
			Number: 4,
		},
	}

	waitCh := make(chan struct{})

	// send bunch of request to the server from client
	go func() {
		// sending request one by one in stream
		for _, req := range requests {
			log.Printf("Sending request: %v \n", req)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("Error while sending request as stream: %v \n", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// recieve multiple response from server
	go func() {

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				// done with receiving res from the server
				close(waitCh)
			}
			if err != nil {
				log.Fatalf("Error while receiving the response from server in client: %v \n", err)
				close(waitCh)
			}

			fmt.Printf("Receiving response from server : %v \n", res.GetResult())

		}

	}()

	// block until done with the process.
	<-waitCh

}
