package indexer

import (
	"sync"

	"github.com/taubyte/go-project-schema/domains"
	"github.com/taubyte/go-project-schema/libraries"
	projectSchema "github.com/taubyte/go-project-schema/project"
)

func getDomain(name string, app string, project projectSchema.Project) (domObj domains.Domain, err error) {
	domObj, err = project.Domain(name, app)
	if err != nil || (len(domObj.Get().Id()) == 0) {
		domObj, err = project.Domain(name, "")
	}
	return
}

func getLibraries(name string, app string, project projectSchema.Project) (libObj libraries.Library, global bool, err error) {
	if len(app) == 0 {
		global = true
	}

	libObj, err = project.Library(name, app)
	if err != nil || (len(libObj.Get().Id()) == 0) {
		libObj, err = project.Library(name, "")
		global = true
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
