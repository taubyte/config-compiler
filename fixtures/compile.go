package fixtures

import (
	commonTest "bitbucket.org/taubyte/dreamland-test/common"
	"bitbucket.org/taubyte/dreamland/common"
	"github.com/taubyte/config-compiler/compile"
	"github.com/taubyte/go-project-schema/project"
)

func inject(project project.Project, simple common.Simple) error {
	fakeMeta := commonTest.ConfigRepo.HookInfo
	fakeMeta.Repository.Provider = "github"
	fakeMeta.Repository.Branch = "master"
	fakeMeta.HeadCommit.ID = "testCommit"
	rc, err := compile.CompilerConfig(project, fakeMeta)
	if err != nil {
		return err
	}

	compiler, err := compile.New(rc, compile.Dev())
	if err != nil {
		return err
	}
	defer compiler.Close()

	err = compiler.Build()
	if err != nil {
		return err
	}

	// publish ( compile & send to TNS )
	err = compiler.Publish(simple.TNS())
	if err != nil {
		return err
	}

	return nil
}
