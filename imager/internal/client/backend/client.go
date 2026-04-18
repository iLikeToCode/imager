package backend

import (
	pb "imager/gen/pb"
	"log"

	"google.golang.org/grpc"
)

func NewClient(addr string, opts grpc.DialOption) (pb.ImageServiceClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient(addr, opts)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	return pb.NewImageServiceClient(conn), conn
}
