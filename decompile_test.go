package compiler

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/taubyte/config-compiler/decompile"
)

// Change this variable to false to output to ./ rather than temp dir
var runInTemp = true

func TestDecompileBasic(t *testing.T) {
	var err error
	var gitRoot string

	if runInTemp {
		gitRoot, err = ioutil.TempDir("", "gitTestRoot")
		if err != nil {
			t.Error(err)
			return
		}
		defer os.RemoveAll(gitRoot)
	} else {
		gitRoot = "./testdata"
		os.RemoveAll(gitRoot)
	}

	gitRootConfig := gitRoot + "/config"
	os.MkdirAll(gitRootConfig, 0750)

	decompiler, err := decompile.New(afero.NewBasePathFs(afero.NewOsFs(), gitRootConfig), createdProjectObject)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = decompiler.Build()
	if err != nil {
		t.Error(err)
		return
	}
}
