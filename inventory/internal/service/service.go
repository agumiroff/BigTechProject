package service

import (
	"context"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
)

type InvService interface {
	GetPart(ctx context.Context, uuid string) (res *model.Part, err error)
	ListParts(ctx context.Context, f *model.PartsFilter) (res []*model.Part, err error)
	CreatePart(ctx context.Context, f *model.Part) (res string, err error)
}
