package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/nitin1259/grpc-go-learn/calculator/calcpb"
	"google.golang.org/grpc"
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
