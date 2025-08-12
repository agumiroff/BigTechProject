module github.com/agumiroff/BigTechProject/order/v1

go 1.24.4

replace github.com/agumiroff/BigTechProject/shared => ./../shared

require (
	github.com/brianvoe/gofakeit/v6 v6.28.0
	github.com/go-chi/chi/v5 v5.2.2
	github.com/go-faster/errors v0.7.1
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.74.0
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250721164621-a45f3dfb1074 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
