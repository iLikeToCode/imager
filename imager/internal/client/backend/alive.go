package backend

import (
	"context"
	"time"

	pb "imager/gen/pb"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AliveClient struct {
	client pb.AliveClient
}

func newAliveClient(conn *grpc.ClientConn) *AliveClient {
	return &AliveClient{
		client: pb.NewAliveClient(conn),
	}
}

func (c *AliveClient) Alive() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.client.Alive(ctx, &emptypb.Empty{})

	if err != nil {
		return err
	}

	return nil
}
