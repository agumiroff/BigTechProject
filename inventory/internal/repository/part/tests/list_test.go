package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/mocks"
	repomodel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
)

type RepoSuite struct {
	suite.Suite
	repo *mocks.InvRepository
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}

func (s *RepoSuite) SetupTest() {
	s.repo = mocks.NewInvRepository(s.T())
}

func (s *RepoSuite) TestListPartsFilterSuccess() {
	ctx := s.createContext()
	expectedList := []*model.Part{
		{
			Uuid: "123",
			Name: "Engine",
		},
		{
			Uuid: "222",
			Name: "Wing",
		},
	}

	filter := &repomodel.PartsFilter{
		Names: []string{"Engine"},
	}

	s.repo.EXPECT().ListParts(ctx, filter).Return(expectedList[:1], nil)

	res, err := s.repo.ListParts(ctx, filter)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Equal(1, len(res))
	s.Equal(expectedList[0], res[0])
	s.Equal("123", res[0].Uuid)
}

func (s *RepoSuite) TestListPartsFilterNoMatch() {
	ctx := s.createContext()
	filter := &repomodel.PartsFilter{
		Names: []string{"NonExistentPart"},
	}

	s.repo.EXPECT().ListParts(ctx, filter).Return(nil, nil)

	res, err := s.repo.ListParts(ctx, filter)

	s.Require().NoError(err)
	s.Require().Empty(res, "Result should be empty when no parts match the filter")
}

func (s *RepoSuite) TestListPartsEmptyFilter() {
	ctx := s.createContext()
	allParts := []*model.Part{
		{
			Uuid: "123",
			Name: "Engine",
		},
		{
			Uuid: "222",
			Name: "Wing",
		},
		{
			Uuid: "333",
			Name: "Porthole",
		},
	}

	emptyFilter := &repomodel.PartsFilter{}

	s.repo.EXPECT().ListParts(ctx, emptyFilter).Return(allParts, nil)

	res, err := s.repo.ListParts(ctx, emptyFilter)

	s.Require().NoError(err)
	s.Require().NotEmpty(res)
	s.Equal(3, len(res))
	s.Equal("123", res[0].Uuid)
	s.Equal("222", res[1].Uuid)
	s.Equal("333", res[2].Uuid)
}

func (s *RepoSuite) createContext() context.Context {
	return context.Background()
}
