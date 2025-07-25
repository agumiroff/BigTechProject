package service

import (
	"context"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	rModel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
)

type InvService interface {
	GetPart(ctx context.Context, uuid string) (res *model.Part, err error)
	LitParts(ctx context.Context, f *rModel.PartsFilter) (res []*model.Part, err error)
}
