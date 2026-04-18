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

func NewImageClient(addr string, opts grpc.DialOption) (pb.ImageServiceClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient(addr, opts)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	return pb.NewImageServiceClient(conn), conn
}

func ListImages(client pb.ImageServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ListImages(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("ListImages failed: %v", err)
	}

	fmt.Println("Available images:")
	for _, img := range resp.Images {
		fmt.Printf("ID: %d | %s | %d bytes\n", img.Id, img.Name, img.Size)
	}
}

func PullImage(client pb.ImageServiceClient, id int32) {
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
