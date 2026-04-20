package backend

import (
	"log"

	"google.golang.org/grpc"
)

type Client struct {
	Conn   *grpc.ClientConn
	Images *ImageClient
}

func NewClient(addr string, opts grpc.DialOption) Client {
	conn, err := grpc.NewClient(addr, opts)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	return Client{
		Conn:   conn,
		Images: newImageClient(conn),
	}
}
