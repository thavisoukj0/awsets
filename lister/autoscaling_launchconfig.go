package lister

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/trek10inc/awsets/arn"
	"github.com/trek10inc/awsets/option"
	"github.com/trek10inc/awsets/resource"
)

type AWSAutoscalingLaunchConfiguration struct {
}

func init() {
	i := AWSAutoscalingLaunchConfiguration{}
	listers = append(listers, i)
}

func (l AWSAutoscalingLaunchConfiguration) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.AutoscalingLaunchConfig}
}

func (l AWSAutoscalingLaunchConfiguration) List(cfg option.AWSetsConfig) (*resource.Group, error) {
	svc := autoscaling.NewFromConfig(cfg.AWSCfg)

	rg := resource.NewGroup()
	err := Paginator(func(nt *string) (*string, error) {
		res, err := svc.DescribeLaunchConfigurations(cfg.Context, &autoscaling.DescribeLaunchConfigurationsInput{
			MaxRecords: aws.Int32(100),
			NextToken:  nt,
		})
		if err != nil {
			return nil, err
		}
		for _, v := range res.LaunchConfigurations {

			configArn := arn.ParseP(v.LaunchConfigurationARN)
			r := resource.New(cfg, resource.AutoscalingLaunchConfig, configArn.ResourceId, v.LaunchConfigurationName, v)

			for _, sg := range v.SecurityGroups {
				r.AddRelation(resource.Ec2SecurityGroup, sg, "")
			}
			rg.AddResource(r)
		}

		return res.NextToken, nil
	})
	return rg, err
}
