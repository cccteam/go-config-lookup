package awslookuper

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	gomock "github.com/golang/mock/gomock"
)

func TestAwsSsmLookuper_Lookup(t *testing.T) {
	t.Parallel()

	var notFoundErr *types.ParameterNotFound
	errors.As(errors.New("Not found"), &notFoundErr)

	type fields struct {
		ctx context.Context
	}
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      string
		want1     bool
		wantPanic bool
		prepare   func(ssm *MockAwsSsmAPI)
	}{
		{
			name: "fails to find value",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				key: "test",
			},
			wantPanic: false,
			prepare: func(ssmMock *MockAwsSsmAPI) {
				ssmMock.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(nil, notFoundErr)
			},
		},
		{
			name: "fails fetching value and panics",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				key: "test",
			},
			want:      "",
			want1:     false,
			wantPanic: true,
			prepare: func(ssmMock *MockAwsSsmAPI) {
				ssmMock.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(nil, errors.New("Not found"))
			},
		},
		{
			name: "success finding value",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				key: "test",
			},
			want:      "test",
			want1:     true,
			wantPanic: false,
			prepare: func(ssmMock *MockAwsSsmAPI) {
				ssmMock.EXPECT().GetParameter(context.Background(), &ssm.GetParameterInput{
					Name:           aws.String("test"),
					WithDecryption: aws.Bool(true),
				}).Return(&ssm.GetParameterOutput{
					Parameter: &types.Parameter{
						Value: aws.String("test"),
					},
				}, nil)
			},
		},
		{
			name: "success finding value but it is nil",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				key: "test",
			},
			want:      "",
			want1:     false,
			wantPanic: false,
			prepare: func(ssmMock *MockAwsSsmAPI) {
				ssmMock.EXPECT().GetParameter(context.Background(), &ssm.GetParameterInput{
					Name:           aws.String("test"),
					WithDecryption: aws.Bool(true),
				}).Return(&ssm.GetParameterOutput{
					Parameter: nil,
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockapi := NewMockAwsSsmAPI(ctrl)
			tt.prepare(mockapi)
			a := &AwsSsmLookuper{
				ssm: mockapi,
				ctx: tt.fields.ctx,
			}
			defer func() {
				if err := recover(); (err != nil) != tt.wantPanic {
					t.Errorf("AwsSsmLookuper.Lookup() panic = %v", err)
				}
			}()
			got, got1 := a.Lookup(tt.args.key)
			if got != tt.want {
				t.Errorf("AwsSsmLookuper.Lookup() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AwsSsmLookuper.Lookup() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
