package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/nitin1259/grpc-go-learn/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Calculator(ctx context.Context, req *calcpb.CalcRequest) (*calcpb.CalcResponse, error) {

	log.Printf("Coming in server to calucate req: %v \n", req)

	num1 := req.Calcate.GetNum1()
	num2 := req.Calcate.GetNum2()

	result := num1 + num2

	res := &calcpb.CalcResponse{
		Result: result,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest, stream calcpb.CalcService_PrimeNumberDecompositionServer) error {
	log.Printf("server streaming rpc to find prime decompsition req: %v \n", req)
	num := req.GetNumber()
	k := int64(2)
	for num > 1 {
		if num%k == 0 {
			fmt.Printf("Prime factor %v", k)
			res := &calcpb.PrimeNumberDecompositionResponse{
				Result: k,
			}
			stream.Send(res)
			num = num / k
			time.Sleep(1000 * time.Millisecond)
		} else {
			k++
		}
	}

	return nil
}

func (*server) ComputeAverage(stream calcpb.CalcService_ComputeAverageServer) error {
	log.Println("Invoking Compute Average server for clinet streaming rpc")
	count := int64(0)
	sum := int64(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we are done with getting the stream from client
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Result: sum / count,
			})

		}
		count++
		num := req.GetNumber()
		sum += num
	}

	return nil

}

func (*server) FindMaximum(stream calcpb.CalcService_FindMaximumServer) error {
	log.Println("Invoking FindMaximum on server for BiDi streaming rpc")
	maxNum := int64(0)
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			// done with recieving req from client
			return nil
		}
		if err != nil {
			log.Fatalf("Error while receivng the BiDi Stream from client: %v", err)
			return nil
		}
		num := req.GetNumber()

		if num > maxNum {
			maxNum = num
			errSend := stream.Send(&calcpb.FindMaximumResponse{
				Result: maxNum,
			})
			if errSend != nil {
				log.Fatalf("Error while sending BiDi stream to client : %v", errSend)
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {
	log.Println("Invoking SquareRoot on server for Error handling in rpc")
	num := req.GetNumber()

	if num < 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("Recieved a negative number: %v", num))
	}
	result := int64(math.Sqrt(float64(num)))
	fmt.Println("Square root of the number: ", result)
	return &calcpb.SquareRootResponse{
		Result: result,
	}, nil
}

func main() {
	fmt.Println("running in calc server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Error while listening: %v", err)
	}

	s := grpc.NewServer()

	calcpb.RegisterCalcServiceServer(s, &server{})

	s.Serve(lis)

	log.Printf("Server started for calc...")
}
