package model

// All domain-specific errors have been moved to shared/apperrors
// This file is kept for backward compatibility until all references are updated
import "github.com/agumiroff/BigTechProject/shared/apperrors"

var (
	ErrOrderNotFound         = apperrors.ErrNotFound
	ErrOrderAlreadyCancelled = apperrors.ErrForbidden
	ErrOrderAlreadyPaid      = apperrors.ErrForbidden
	ErrCreateOrderFailed     = apperrors.ErrInternal
	ErrUpdateOrderFailed     = apperrors.ErrInternal
	ErrInvalidOrderUUID      = apperrors.ErrInvalidRequest
	ErrDatabase              = apperrors.ErrInternal
)
