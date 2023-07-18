package indexer

import (
	"sync"

	"github.com/taubyte/go-project-schema/domains"
	projectSchema "github.com/taubyte/go-project-schema/project"
)

func getDomain(name string, app string, project projectSchema.Project) (domObj domains.Domain, err error) {
	domObj, err = project.Domain(name, app)
	if err != nil || (len(domObj.Get().Id()) == 0) {
		domObj, err = project.Domain(name, "")
	}
	return
}

type IndexContext struct {
	AppId     string
	AppName   string
	ProjectId string
	Branch    string
	Commit    string
	Obj       map[string]interface{}

	Dev bool

	validDomainsLock sync.Mutex
	ValidDomains     []string
	DVPublicKey      []byte
}
