package decompile

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	projectLib "github.com/taubyte/go-project-schema/project"
	structureSpec "github.com/taubyte/go-specs/structure"
)

func storage(project projectLib.Project, _id string, obj interface{}, appName string) error {
	resource := &structureSpec.Storage{}
	mapstructure.Decode(obj, resource)

	iFace, err := project.Storage(resource.Name, appName)
	if err != nil {
		return fmt.Errorf("Open storage `%s/%s` failed: %s", appName, resource.Name, err)
	}

	resource.SetId(_id)
	return iFace.SetWithStruct(false, resource)
}