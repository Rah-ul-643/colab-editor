package main

import (
	"fmt"
	"os"
	"strconv"

	"colab/client"
	"colab/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <server|client> [address:port]")
		return
	}

	mode := os.Args[1]

	switch mode {
	case "server":
		port := 8080 // default port
		protocol := "tcp"
		if len(os.Args) > 2 {
			// Basic argument parsing for the port.
			p, err := strconv.Atoi(os.Args[2])
			if err == nil {
				port = p
			}
		}
		// Create and start the server.
		serv := server.NewServer(port, protocol)
		if err := serv.ListenAndServe(); err != nil {
			fmt.Println("Server error:", err)
		}

	case "client":
		address := "localhost:8080" // default address
		if len(os.Args) > 2 {
			address = os.Args[2]
		}
		// Start the client editor.
		if err := client.StartEditorAndSync(address); err != nil {
			fmt.Println("Client error:", err)
		}

	default:
		fmt.Println("Unknown mode:", mode)
		fmt.Println("Usage: go run main.go <server|client> [address:port]")
	}
}
