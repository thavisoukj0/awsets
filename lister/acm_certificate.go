package lister

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/trek10inc/awsets/arn"
	"github.com/trek10inc/awsets/option"
	"github.com/trek10inc/awsets/resource"
)

type AWSAcmCertificate struct {
}

func init() {
	i := AWSAcmCertificate{}
	listers = append(listers, i)
}

func (l AWSAcmCertificate) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.AcmCertificate}
}

func (l AWSAcmCertificate) List(cfg option.AWSetsConfig) (*resource.Group, error) {
	svc := acm.NewFromConfig(cfg.AWSCfg)

	rg := resource.NewGroup()
	err := Paginator(func(nt *string) (*string, error) {
		req, err := svc.ListCertificates(cfg.Context, &acm.ListCertificatesInput{
			MaxItems:  aws.Int32(100),
			NextToken: nt,
		})
		if err != nil {
			return nil, err
		}
		for _, cert := range req.CertificateSummaryList {

			res, err := svc.DescribeCertificate(cfg.Context, &acm.DescribeCertificateInput{CertificateArn: cert.CertificateArn})
			if err != nil {
				return nil, fmt.Errorf("unable to describe certificate %s: %w", *cert.CertificateArn, err)
			}
			//if arn.IsArnP(res.Certificate.CertificateArn) {
			certArn := arn.ParseP(res.Certificate.CertificateArn)
			r := resource.New(cfg, resource.AcmCertificate, certArn.ResourceId, certArn.ResourceId, res.Certificate)
			//}
			tagRes, err := svc.ListTagsForCertificate(cfg.Context, &acm.ListTagsForCertificateInput{
				CertificateArn: cert.CertificateArn,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to list tags for cert %s: %w", *cert.CertificateArn, err)
			}
			for _, tag := range tagRes.Tags {
				r.Tags[*tag.Key] = *tag.Value
			}
			rg.AddResource(r)
		}
		return req.NextToken, nil
	})
	return rg, err
}
