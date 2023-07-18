package fixtures

import (
	"testing"

	commonDreamland "github.com/taubyte/dreamland/core/common"
	dreamland "github.com/taubyte/dreamland/core/services"
	commonIface "github.com/taubyte/go-interfaces/common"
	_ "github.com/taubyte/odo/clients/p2p/tns"
	_ "github.com/taubyte/odo/protocols/tns/service"
)

func TestFakeProject(t *testing.T) {
	u := dreamland.Multiverse("TestFakeProject")
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

	err = u.RunFixture("fakeProject")
	if err != nil {
		t.Error(err)
		return
	}
}
