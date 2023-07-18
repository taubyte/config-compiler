package compiler

import (
	"os"
	"testing"

	commonTest "bitbucket.org/taubyte/dreamland-test/common"
	gitTest "bitbucket.org/taubyte/dreamland-test/git"
	commonDreamland "bitbucket.org/taubyte/dreamland/common"
	dreamland "bitbucket.org/taubyte/dreamland/services"
	_ "bitbucket.org/taubyte/tns-p2p-client"
	_ "bitbucket.org/taubyte/tns/service"
	"github.com/spf13/afero"
	"github.com/taubyte/config-compiler/compile"
	"github.com/taubyte/config-compiler/decompile"
	commonIface "github.com/taubyte/go-interfaces/common"
	projectLib "github.com/taubyte/go-project-schema/project"
	specs "github.com/taubyte/go-specs/methods"
	"github.com/taubyte/utils/maps"
)

func TestDecompileProd(t *testing.T) {
	u := dreamland.Multiverse("scratch")
	defer u.Stop()
	err := u.StartWithConfig(&commonDreamland.Config{
		Services: map[string]commonIface.ServiceConfig{
			"tns": {},
		},
		Simples: map[string]commonDreamland.SimpleConfig{
			"me": {
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

	simple, err := u.Simple("me")
	if err != nil {
		t.Error(err)
		return
	}
	tns := simple.TNS()

	gitRoot := "./testGIT"
	gitRootConfig := gitRoot + "/prodConfigDreamland"
	os.MkdirAll(gitRootConfig, 0755)

	fakeMeta := commonTest.ConfigRepo.HookInfo
	fakeMeta.Repository.SSHURL = "git@github.com:taubyte-test/tb_prodproject.git"
	fakeMeta.Repository.Branch = "dreamland"
	fakeMeta.Repository.Provider = "github"

	err = gitTest.CloneToDirSSH(u.Context(), gitRootConfig, commonTest.Repository{
		ID:       517160737,
		Name:     "tb_prodproject",
		HookInfo: fakeMeta,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// read with seer
	projectIface, err := projectLib.Open(projectLib.SystemFS(gitRootConfig))
	if err != nil {
		t.Error(err)
		return
	}

	rc, err := compile.CompilerConfig(projectIface, fakeMeta)
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

	err = compiler.Publish(tns)
	if err != nil {
		t.Error(err)
		return
	}

	test_obj, err := tns.Fetch(specs.ProjectPrefix(projectIface.Get().Id(), fakeMeta.Repository.Branch, fakeMeta.HeadCommit.ID))
	if test_obj.Interface() == nil {
		t.Error("NO OBject found", err)
		return
	}

	maps.Display("", test_obj)

	testProjectDir := "./testGIT/testDecompileProd"
	os.RemoveAll(testProjectDir)
	os.Mkdir(testProjectDir, 0777)

	decompiler, err := decompile.New(afero.NewBasePathFs(afero.NewOsFs(), testProjectDir), test_obj.Interface())
	if err != nil {
		t.Error(err)
		return
	}

	_, err = decompiler.Build()
	if err != nil {
		t.Error(err)
	}

}
