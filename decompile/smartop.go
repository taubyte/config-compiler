package decompile

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	projectLib "github.com/taubyte/go-project-schema/project"
	structureSpec "github.com/taubyte/go-specs/structure"
)

func smartop(project projectLib.Project, _id string, obj interface{}, appName string) error {
	resource := &structureSpec.SmartOp{}
	mapstructure.Decode(obj, resource)

	iFace, err := project.SmartOps(resource.Name, appName)
	if err != nil {
		return fmt.Errorf("Open smart-op `%s/%s` failed: %s", appName, resource.Name, err)
	}

	resource.SetId(_id)
	return iFace.SetWithStruct(false, resource)
}
