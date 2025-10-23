package main

import (
	"context"

	"github.com/agumiroff/BigTechProject/inventory/v1/server"
)

func main() {
	ctx := context.Background()
	server.StartGRPCServer(ctx)
}
