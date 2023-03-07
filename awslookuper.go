// awslookuper is a package used to lookup values from AWS SSM Parameter Store. It implements the Lookuper interface.
package awslookuper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/go-playground/errors/v5"
)

//go:generate mockgen -package $GOPACKAGE -destination mock_awsssmapi_test.go github.com/AscendiumApps/cpp-aggregator-service/config/awslookuper AwsSsmAPI

type AwsSsmLookuper struct {
	ssm AwsSsmAPI
	ctx context.Context
}

type AwsSsmAPI interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

// New creates a new AWSSSMLookuper
func New(ctx context.Context) (*AwsSsmLookuper, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "config.LoadDefaultConfig()")
	}
	ssmClient := ssm.NewFromConfig(cfg)

	return &AwsSsmLookuper{ssm: ssmClient, ctx: ctx}, nil
}

// Lookup uses the AWS SSM Parameter Store to lookup a value for a given key
func (a *AwsSsmLookuper) Lookup(key string) (string, bool) {
	param, err := a.ssm.GetParameter(a.ctx, &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		var ssmErr *types.ParameterNotFound
		if errors.As(err, &ssmErr) {
			return "", false
		}

		panic(errors.Wrap(err, "ssm.GetParameter()"))
	}

	if param.Parameter == nil || param.Parameter.Value == nil {
		return "", false
	}

	return *param.Parameter.Value, true
}
