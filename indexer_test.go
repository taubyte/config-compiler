package compiler

import (
	"testing"

	commonDreamland "bitbucket.org/taubyte/dreamland/common"
	dreamland "bitbucket.org/taubyte/dreamland/services"
	"github.com/taubyte/config-compiler/compile"
	TestFixtures "github.com/taubyte/config-compiler/internal/fixtures"
	commonIface "github.com/taubyte/go-interfaces/common"
	"github.com/taubyte/go-interfaces/services/tns"
	projectSchema "github.com/taubyte/go-project-schema/project"
	maps "github.com/taubyte/utils/maps"
)

func TestIndexer(t *testing.T) {
	u := dreamland.Multiverse("indexer")
	defer u.Stop()
	err := u.StartWithConfig(&commonDreamland.Config{
		Services: map[string]commonIface.ServiceConfig{
			"tns": {},
		},
		Simples: map[string]commonDreamland.SimpleConfig{
			"client": {
				Clients: commonDreamland.SimpleConfigClients{
					TNS: &commonIface.ClientConfig{},
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	simple, err := u.Simple("client")
	if err != nil {
		t.Error(err)
		return
	}

	fs, err := TestFixtures.VirtualFSWithBuiltProject()
	if err != nil {
		t.Error(err)
		return
	}

	project, err := projectSchema.Open(projectSchema.VirtualFS(fs, "/test_project/config"))
	if err != nil {
		t.Error(err)
		return
	}

	rc, err := compile.CompilerConfig(project, fakeMeta)
	if err != nil {
		t.Error(err)
		return
	}

	compiler, err := compile.New(rc, compile.Dev())
	if err != nil {
		t.Error(err)
		return
	}

	err = compiler.Build()
	if err != nil {
		t.Error(err)
		return
	}

	err = compiler.Publish(simple.TNS())
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := simple.TNS().Lookup(tns.Query{Prefix: []string{"domains"}, RegEx: false})
	_map := make(map[string]interface{})
	_map["test"] = resp
	list, err := maps.StringArray(_map, "test")
	if err != nil {
		t.Error(err)
		return
	}

	if len(list) != 2 { // local/global domains index, project,branch
		t.Errorf("Expected 2 got %d", len(list))
	}
}
