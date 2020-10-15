package lister

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/trek10inc/awsets/option"
	"github.com/trek10inc/awsets/resource"
)

type AWSEc2Subnet struct {
}

func init() {
	i := AWSEc2Subnet{}
	listers = append(listers, i)
}

func (l AWSEc2Subnet) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.Ec2Subnet}
}

func (l AWSEc2Subnet) List(cfg option.AWSetsConfig) (*resource.Group, error) {
	svc := ec2.NewFromConfig(cfg.AWSCfg)

	rg := resource.NewGroup()
	err := Paginator(func(nt *string) (*string, error) {
		res, err := svc.DescribeSubnets(cfg.Context, &ec2.DescribeSubnetsInput{
			MaxResults: aws.Int32(100),
			NextToken:  nt,
		})
		if err != nil {
			return nil, err
		}
		for _, v := range res.Subnets {
			r := resource.New(cfg, resource.Ec2Subnet, v.SubnetId, v.SubnetArn, v)
			r.AddRelation(resource.Ec2Vpc, v.VpcId, "")
			rg.AddResource(r)
		}
		return res.NextToken, nil
	})
	return rg, err
}
