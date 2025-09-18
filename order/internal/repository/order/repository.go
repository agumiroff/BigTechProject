package order

import (
	"github.com/agumiroff/BigTechProject/order/v1/internal/repository/db"
)

type repository struct {
	db db.DB
}

func NewRepository(db db.DB) *repository {
	return &repository{
		db: db,
	}
}
