package lister

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/trek10inc/awsets/option"
	"github.com/trek10inc/awsets/resource"
)

type AWSGlueConnection struct {
}

func init() {
	i := AWSGlueConnection{}
	listers = append(listers, i)
}

func (l AWSGlueConnection) Types() []resource.ResourceType {
	return []resource.ResourceType{
		resource.GlueConnection,
	}
}

func (l AWSGlueConnection) List(cfg option.AWSetsConfig) (*resource.Group, error) {
	svc := glue.NewFromConfig(cfg.AWSCfg)

	rg := resource.NewGroup()
	err := Paginator(func(nt *string) (*string, error) {
		res, err := svc.GetConnections(cfg.Context, &glue.GetConnectionsInput{
			HidePassword: aws.Bool(true),
			MaxResults:   aws.Int32(100),
			NextToken:    nt,
		})
		if err != nil {
			return nil, err
		}
		for _, v := range res.ConnectionList {
			r := resource.New(cfg, resource.GlueConnection, v.Name, v.Name, v)
			rg.AddResource(r)
		}
		return res.NextToken, nil
	})
	return rg, err
}
