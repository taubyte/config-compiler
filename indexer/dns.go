package indexer

import (
	_ "embed"

	dv "github.com/taubyte/domain-validation"
	commonIface "github.com/taubyte/go-interfaces/services/common"
	domainSpec "github.com/taubyte/go-specs/domain"
	"golang.org/x/exp/slices"
)

var (
	//go:embed domain_private.key
	domainValPrivateKeyData []byte
	//go:embed domain_public.key
	domainValPublicKeyData []byte
)

func (ctx *IndexContext) validateDomain(fqdn string) error {
	ctx.validDomainsLock.Lock()
	defer ctx.validDomainsLock.Unlock()

	if ctx.ValidDomains == nil {
		ctx.ValidDomains = []string{}
	}

	if slices.Contains(ctx.ValidDomains, fqdn) == true {
		return nil
	}

	var err error
	if commonIface.Deployment == commonIface.Odo {
		err = domainSpec.ValidateDNS(ctx.ProjectId, fqdn, ctx.Dev, dv.PublicKey(ctx.DVPublicKey))
	} else {
		err = domainSpec.ValidateDNS(ctx.ProjectId, fqdn, ctx.Dev, dv.PublicKey(domainValPublicKeyData))
	}
	if err != nil {
		return err
	}

	ctx.ValidDomains = append(ctx.ValidDomains, fqdn)
	return nil
}
