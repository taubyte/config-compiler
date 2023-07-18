package compiler

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	commonTest "bitbucket.org/taubyte/dreamland-test/common"
	gitTest "bitbucket.org/taubyte/dreamland-test/git"
	"github.com/spf13/afero"
	"github.com/taubyte/config-compiler/compile"
	"github.com/taubyte/config-compiler/internal/fixtures"
	projectLib "github.com/taubyte/go-project-schema/project"
	"github.com/taubyte/utils/maps"
)

func TestCompile(t *testing.T) {
	project, err := fixtures.Project()
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

	jsonBytes, err := json.Marshal(compiler.Object())
	if err != nil {
		t.Error(err)
		return
	}

	osFS := afero.NewOsFs()
	// Uncomment the below to refresh the file if changes to the yaml written in internal
	// err = afero.WriteFile(osFS, "./compile_test2.json", jsonBytes, 0644)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	expectedBytes, err := afero.ReadFile(osFS, "./compile_test.json")
	if err != nil {
		t.Error(err)
		return
	}

	if string(expectedBytes) != string(jsonBytes) {
		t.Error("Bytes don't match")
		return
	}

}

func TestFromCloneCompile(t *testing.T) {
	testCtx, testCtxC := context.WithCancel(context.Background())
	defer func() {
		s := (3 * time.Second)
		go func() {
			time.Sleep(s)
			testCtxC()
		}()
		time.Sleep(s)
	}()

	gitRoot := "./testGIT"
	defer os.RemoveAll(gitRoot)
	gitRootConfig := gitRoot + "/config"
	os.MkdirAll(gitRootConfig, 0755)

	// clone repo
	err := gitTest.CloneToDirSSH(testCtx, gitRootConfig, commonTest.ConfigRepo)
	if err != nil {
		t.Error(err)
		return
	}

	// read with seer
	project, err := projectLib.Open(projectLib.SystemFS(gitRootConfig))
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

	maps.Display("", compiler.Object())
}
