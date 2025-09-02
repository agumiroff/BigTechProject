package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/mocks"
)

type RepositorySuite struct {
	suite.Suite
	repo *mocks.InvRepository
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) SetupTest() {
	s.repo = mocks.NewInvRepository(s.T())
}

func (s *RepositorySuite) createContext() context.Context {
	return context.Background()
}
