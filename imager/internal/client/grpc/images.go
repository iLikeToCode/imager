package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "imager/gen/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func ListImages() {}

func PullImage() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewImageServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ---- List images ----
	resp, err := client.ListImages(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("ListImages failed: %v", err)
	}

	fmt.Println("Available images:")
	for _, img := range resp.Images {
		fmt.Printf("ID: %d | %s | %d bytes\n", img.Id, img.Name, img.Size)
	}

	// ---- Pull first image (example) ----
	if len(resp.Images) == 0 {
		log.Println("no images available")
		return
	}

	imageID := resp.Images[0].Id

	stream, err := client.PullImage(context.Background(), &pb.PullImageRequest{
		Id: imageID,
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
