package main

import (
	"context"
	"fmt"
	"log"
	"net"

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
