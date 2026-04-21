package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"

	pb "imager/gen/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedImageServiceServer
	images []*pb.Image
}

const IMAGE_PATH = "./images"

func NewServer() *Server {
	dir := "./images"

	log.Println("Scanning images directory:", dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read images dir: %v", err)
	}

	var images []*pb.Image

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		path := filepath.Join(dir, name)

		info, err := e.Info()
		if err != nil {
			log.Fatalf("failed to stat %s: %v", path, err)
		}

		img := &pb.Image{
			Id:   uint32(len(images) + 1), // simple sequential ID
			Name: name,
			Size: uint64(info.Size()),
		}

		log.Printf("Loaded image: %s size=%d", name, img.Size)

		images = append(images, img)
	}

	return &Server{
		images: images,
	}
}

func (s *Server) Alive(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Server) ListImages(ctx context.Context, req *emptypb.Empty) (*pb.ListImagesResponse, error) {
	return &pb.ListImagesResponse{
		Images: s.images,
	}, nil
}

func getImagePath(name string) (string, error) {
	return path.Join(IMAGE_PATH, name), nil
}

func (s *Server) getImageName(id uint32) (string, error) {
	for _, img := range s.images {
		if img.Id == id {
			return "./" + img.Name, nil
		}
	}
	return "", fmt.Errorf("image not found")
}

func (s *Server) PullImage(req *pb.PullImageRequest, stream pb.ImageService_PullImageServer) error {
	name, err := s.getImageName(req.Id)
	if err != nil {
		return err
	}
	path, err := getImagePath(name)
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 1024*1024) // 1MB chunks

	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = stream.Send(&pb.ImageChunk{
			Data: buf[:n],
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := NewServer()
	pb.RegisterImageServiceServer(grpcServer, server)
	pb.RegisterAliveServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
