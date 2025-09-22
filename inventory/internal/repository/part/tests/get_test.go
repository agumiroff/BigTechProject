package tests

import (
	"errors"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
)

func (s *RepositorySuite) TestGetSuccess() {
	ctx := s.createContext()
	expectedPart := &model.Part{
		Uuid: "123",
		Name: "Engine",
	}

	s.repo.EXPECT().GetPart(ctx, "123").Return(expectedPart, nil)

	res, err := s.repo.GetPart(ctx, "123")

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Equal("123", res.Uuid)
	s.Equal("Engine", res.Name)
}

func (s *RepositorySuite) TestGetError() {
	ctx := s.createContext()
	s.repo.EXPECT().GetPart(ctx, "123").Return(nil, errors.New("not found"))

	res, err := s.repo.GetPart(ctx, "123")

	s.Require().Error(err)
	s.Require().Nil(res)
}
