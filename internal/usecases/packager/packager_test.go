package packager

import (
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"homework/internal/domain"
	"homework/internal/usecases/packager/mocks"
	"testing"
)

func TestOrderPackager_PackageOrder(t *testing.T) {
	t.Parallel()

	type fields struct {
		strategies map[domain.PackagingType]OrderPackagerStrategy
	}

	type args struct {
		order     domain.PVZOrder
		packaging domain.PackagingType
	}

	ctrl := minimock.NewController(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "box packaging",
			fields: fields{
				strategies: map[domain.PackagingType]OrderPackagerStrategy{
					domain.PackagingTypeBox: mocks.NewOrderPackagerStrategyMock(ctrl).PackageOrderMock.Return(domain.PVZOrder{}, nil),
				},
			},
			args: args{
				order:     domain.PVZOrder{},
				packaging: domain.PackagingTypeBox,
			},
			wantErr: assert.NoError,
		},
		{
			name: "unknown packaging type",
			fields: fields{
				strategies: map[domain.PackagingType]OrderPackagerStrategy{},
			},
			args: args{
				order:     domain.PVZOrder{},
				packaging: domain.PackagingTypeUnknown,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err) && errors.Is(err, domain.ErrInvalidArgument)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderPackager{
				strategies: tt.fields.strategies,
			}
			_, err := o.PackageOrder(tt.args.order, tt.args.packaging)
			tt.wantErr(t, err)
		})
	}
}
