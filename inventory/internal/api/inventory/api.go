package api

import (
	"context"

	invV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
)

type InvAPI interface {
	GetPart(ctx context.Context, req *invV1.GetPartRequest) (res *invV1.GetPartResponse, err error)
	ListParts(ctx context.Context, req *invV1.ListPartsRequest) (*invV1.ListPartsResponse, error)
}
