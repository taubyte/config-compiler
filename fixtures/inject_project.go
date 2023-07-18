package fixtures

import (
	"fmt"

	"github.com/taubyte/dreamland/core/common"
	"github.com/taubyte/go-project-schema/project"
)

func injectProject(u common.Universe, params ...interface{}) error {
	simple, err := u.Simple("client")
	if err != nil {
		return fmt.Errorf("failed getting simple with error: %v", err)
	}

	err = simple.Provides("tns")
	if err != nil {
		return err
	}

	project, ok := params[0].(project.Project)
	if ok == false {
		return fmt.Errorf("Param 0 not a valid project to inject got %T", params[0])
	}

	return inject(project, simple)
}
