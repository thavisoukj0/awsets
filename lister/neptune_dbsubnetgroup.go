package lister

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/trek10inc/awsets/arn"
	"github.com/trek10inc/awsets/option"
	"github.com/trek10inc/awsets/resource"
)

type AWSNeptuneDbSubnetGroup struct {
}

func init() {
	i := AWSNeptuneDbSubnetGroup{}
	listers = append(listers, i)
}

func (l AWSNeptuneDbSubnetGroup) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.NeptuneDbSubnetGroup}
}

func (l AWSNeptuneDbSubnetGroup) List(cfg option.AWSetsConfig) (*resource.Group, error) {
	svc := neptune.NewFromConfig(cfg.AWSCfg)

	rg := resource.NewGroup()
	err := Paginator(func(nt *string) (*string, error) {
		res, err := svc.DescribeDBSubnetGroups(cfg.Context, &neptune.DescribeDBSubnetGroupsInput{
			MaxRecords: aws.Int32(100),
			Marker:     nt,
		})
		if err != nil {
			return nil, err
		}
		for _, v := range res.DBSubnetGroups {
			subnetArn := arn.ParseP(v.DBSubnetGroupArn)
			r := resource.New(cfg, resource.NeptuneDbSubnetGroup, subnetArn.ResourceId, "", v)
			r.AddRelation(resource.Ec2Vpc, v.VpcId, "")
			for _, subnet := range v.Subnets {
				r.AddRelation(resource.Ec2Subnet, subnet.SubnetIdentifier, "")
			}
			rg.AddResource(r)
		}
		return res.Marker, nil
	})
	return rg, err
}
