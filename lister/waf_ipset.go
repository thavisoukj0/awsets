package lister

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/waf"
	"github.com/trek10inc/awsets/option"
	"github.com/trek10inc/awsets/resource"
)

var listWafIpSetsOnce sync.Once

type AWSWafIpSet struct {
}

func init() {
	i := AWSWafIpSet{}
	listers = append(listers, i)
}

func (l AWSWafIpSet) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.WafIpSet}
}

func (l AWSWafIpSet) List(cfg option.AWSetsConfig) (*resource.Group, error) {
	svc := waf.NewFromConfig(cfg.AWSCfg)
	rg := resource.NewGroup()

	var outerErr error

	listWafIpSetsOnce.Do(func() {
		outerErr = Paginator(func(nt *string) (*string, error) {
			res, err := svc.ListIPSets(cfg.Context, &waf.ListIPSetsInput{
				Limit:      aws.Int32(100),
				NextMarker: nt,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to list ip sets: %w", err)
			}
			for _, ipsetId := range res.IPSets {
				ipset, err := svc.GetIPSet(cfg.Context, &waf.GetIPSetInput{IPSetId: ipsetId.IPSetId})
				if err != nil {
					return nil, fmt.Errorf("failed to get ipset %s: %w", *ipsetId.IPSetId, err)
				}
				if ipset.IPSet == nil {
					continue
				}
				r := resource.NewGlobal(cfg, resource.WafIpSet, ipset.IPSet.IPSetId, ipset.IPSet.Name, ipset.IPSet)
				rg.AddResource(r)
			}
			return res.NextMarker, nil
		})
	})
	return rg, outerErr
}
