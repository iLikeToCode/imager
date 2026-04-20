package backend

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "imager/gen/pb"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ImageClient struct {
	client pb.ImageServiceClient
}

func newImageClient(conn *grpc.ClientConn) *ImageClient {
	return &ImageClient{
		client: pb.NewImageServiceClient(conn),
	}
}

func (c *ImageClient) ListImages() ([]*pb.Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := c.client.ListImages(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}

	return resp.Images, nil
}

func PullImage(client pb.ImageServiceClient, id uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.PullImage(ctx, &pb.PullImageRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("PullImage failed: %v", err)
	}

	out, err := os.Create("downloaded.wim")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer out.Close()

	var total int64

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream error: %v", err)
		}

		n, err := out.Write(chunk.Data)
		if err != nil {
			log.Fatalf("write error: %v", err)
		}

		total += int64(n)
	}

	fmt.Printf("Download complete: %d bytes written\n", total)
}
