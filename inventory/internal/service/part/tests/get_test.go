package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/agumiroff/BigTechProject/inventory/v1/internal/model"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/repository/mocks"
	"github.com/agumiroff/BigTechProject/inventory/v1/internal/service/part"
)

func TestService_GetPart(t *testing.T) {
	tests := []struct {
		name    string
		uuid    string
		setup   func(*mocks.InvRepository)
		want    *model.Part
		wantErr error
	}{
		{
			name: "successful get",
			uuid: "test-uuid",
			setup: func(m *mocks.InvRepository) {
				m.On("GetPart", mock.Anything, "test-uuid").Return(&model.Part{
					Uuid:          "test-uuid",
					Name:          "Test Part",
					Description:   "Test Description",
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
				}, nil)
			},
			want: &model.Part{
				Uuid:          "test-uuid",
				Name:          "Test Part",
				Description:   "Test Description",
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
			wantErr: nil,
		},
		{
			name: "not found error",
			uuid: "non-existent-uuid",
			setup: func(m *mocks.InvRepository) {
				m.On("GetPart", mock.Anything, "non-existent-uuid").Return(&model.Part{}, errors.New("part not found"))
			},
			want:    &model.Part{},
			wantErr: errors.New("part not found"),
		},
		{
			name: "empty uuid",
			uuid: "",
			setup: func(m *mocks.InvRepository) {
				// No repository call is expected because we check for empty UUID first
			},
			want:    &model.Part{},
			wantErr: errors.New("uuid is empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewInvRepository(t)
			tt.setup(mockRepo)
			svc := part.NewService(mockRepo)

			got, err := svc.GetPart(context.Background(), tt.uuid)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
