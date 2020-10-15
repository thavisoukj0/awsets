package lister

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/trek10inc/awsets/arn"
	"github.com/trek10inc/awsets/option"
	"github.com/trek10inc/awsets/resource"
)

type AWSSagemakerEndpoint struct {
}

func init() {
	i := AWSSagemakerEndpoint{}
	listers = append(listers, i)
}

func (l AWSSagemakerEndpoint) Types() []resource.ResourceType {
	return []resource.ResourceType{
		resource.SagemakerEndpoint,
	}
}

func (l AWSSagemakerEndpoint) List(cfg option.AWSetsConfig) (*resource.Group, error) {
	svc := sagemaker.NewFromConfig(cfg.AWSCfg)

	rg := resource.NewGroup()
	err := Paginator(func(nt *string) (*string, error) {
		res, err := svc.ListEndpoints(cfg.Context, &sagemaker.ListEndpointsInput{
			MaxResults: aws.Int32(100),
			NextToken:  nt,
		})
		if err != nil {
			return nil, err
		}
		for _, ep := range res.Endpoints {
			v, err := svc.DescribeEndpoint(cfg.Context, &sagemaker.DescribeEndpointInput{
				EndpointName: ep.EndpointName,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to describe endpoint %s: %w", *ep.EndpointName, err)
			}
			epArn := arn.ParseP(v.EndpointArn)
			r := resource.New(cfg, resource.SagemakerEndpoint, epArn.ResourceId, v.EndpointName, v)
			if v.DataCaptureConfig != nil {
				r.AddARNRelation(resource.KmsKey, v.DataCaptureConfig.KmsKeyId)
			}
			rg.AddResource(r)
		}
		return res.NextToken, nil
	})
	return rg, err
}
