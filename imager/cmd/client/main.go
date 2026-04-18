package main

import (
	backend "imager/internal/client/backend"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	client, conn := backend.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	backend.ListImages(client)

	//tui.RunTui()
}
