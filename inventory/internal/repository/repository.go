package repository

import (
	"context"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	repomodel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
)

type InvRepository interface {
	GetPart(ctx context.Context, uuid string) (res *model.Part, err error)
	ListParts(ctx context.Context, filter *repomodel.PartsFilter) ([]*model.Part, error)
}
