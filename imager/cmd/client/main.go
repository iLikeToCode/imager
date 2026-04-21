package main

import (
	"fmt"
	backend "imager/internal/client/backend"
	"imager/internal/client/tui"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	var ip string
	var port string

	for {
		fmt.Print("Server IP: ")
		fmt.Scan(&ip)
		fmt.Print("Server Port: ")
		fmt.Scan(&port)

		client := backend.NewClient(fmt.Sprintf("%s:%s", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		defer client.Conn.Close()

		_err := client.Alive.Alive()

		if _err != nil {
			log.Println(_err)
		}

		err := status.Code(_err)
		if err == codes.Unavailable {
			log.Println("Unable to connect to server")
		}

		if _err != nil {
			continue
		}

		tui.RunTui(client)

		break
	}

	//fmt.Println("Press ENTER to shutdown...")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	//exec.Command("shutdown", "now").Run()
}
