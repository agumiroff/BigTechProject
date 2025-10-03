package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	rConverter "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/converter"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/mocks"
	repomodel "github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/model"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/service/part"
)

func TestServiceList(t *testing.T) {
	mockRepo := mocks.NewInvRepository(t)
	svc := part.NewService(mockRepo)
	testFilter := &model.PartsFilter{
		Uuids: []string{"test-uuid-1"},
	}

	tests := []struct {
		name    string
		filter  *model.PartsFilter
		setup   func(*mocks.InvRepository)
		want    []*model.Part
		wantErr error
	}{
		{
			name:   "successful get part",
			filter: testFilter, setup: func(m *mocks.InvRepository) {
				m.On("ListParts", mock.Anything, rConverter.FilterToRepo(testFilter)).Return([]*model.Part{
					{
						Uuid:          "test-uuid-1",
						Name:          "Test Part 1",
						Description:   "Test Description 1",
						Price:         100.0,
						StockQuantity: 5,
						Category:      int32(model.CategoryEngine),
						Dimensions: model.Dimensions{
							Length: 10.0,
							Width:  20.0,
							Height: 30.0,
							Weight: 5.0,
						},
						Manufacturer: model.Manufacturer{
							Name:    "Test Manufacturer",
							Country: "Test Country",
							Website: "http://test.com",
						},
					},
				}, nil)
			},
			want: []*model.Part{
				{
					Uuid:          "test-uuid-1",
					Name:          "Test Part 1",
					Description:   "Test Description 1",
					Price:         100.0,
					StockQuantity: 5,
					Category:      int32(model.CategoryEngine),
					Dimensions: model.Dimensions{
						Length: 10.0,
						Width:  20.0,
						Height: 30.0,
						Weight: 5.0,
					},
					Manufacturer: model.Manufacturer{
						Name:    "Test Manufacturer",
						Country: "Test Country",
						Website: "http://test.com",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "empty list",
			filter: &model.PartsFilter{
				Uuids: []string{"test-uuid"},
			},
			setup: func(m *mocks.InvRepository) {
				m.On("ListParts", mock.Anything, mock.MatchedBy(func(f *repomodel.PartsFilter) bool {
					if len(f.UUIDs) != 1 {
						return false
					}
					return f.UUIDs[0] == "test-uuid"
				})).Return([]*model.Part{}, nil)
			},
			want:    []*model.Part{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockRepo)

			got, err := svc.ListParts(context.Background(), tt.filter)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
