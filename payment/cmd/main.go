package main

import (
	"context"

	"github.com/agumiroff/BigTechProject/payment/v1/server"
)

const configPath = "../../deploy/compose/inventory/.env"

func main() {
	ctx := context.Background()
	server.StartGRPCServer(ctx)
}
