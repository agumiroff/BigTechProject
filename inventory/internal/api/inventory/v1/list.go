package inventory

import (
	"context"
	"log"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/converter"
	InvV1 "github.com/agumiroff/BigTechProject/shared/pkg/proto/inventory/v1"
	"github.com/agumiroff/BigTechProject/shared/v1/apperrors"
)

func (a *api) ListParts(ctx context.Context, req *InvV1.ListPartsRequest) (*InvV1.ListPartsResponse, error) {
	list, err := a.service.ListParts(ctx, converter.FilterToModel(req.Filter))
	if err != nil {
		log.Printf("There is no any part, %v", err)
		return &InvV1.ListPartsResponse{
			Parts: []*InvV1.Part{},
		}, apperrors.Map(err)
	}
	return &InvV1.ListPartsResponse{
		Parts: converter.ModelsToProto(list),
	}, nil
}
