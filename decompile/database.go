package decompile

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	projectLib "github.com/taubyte/go-project-schema/project"
	structureSpec "github.com/taubyte/go-specs/structure"
)

// TODO generic
// type test[T any] interface {
// 	SetWithStruct(sync bool, db T) error
// }

// var _ test[structureSpec.Database] = &databases.Database{}

func database(project projectLib.Project, _id string, obj interface{}, appName string) error {
	resource := &structureSpec.Database{}
	mapstructure.Decode(obj, resource)

	iFace, err := project.Database(resource.Name, appName)
	if err != nil {
		return fmt.Errorf("Open database `%s/%s` failed: %s", appName, resource.Name, err)
	}

	resource.SetId(_id)
	return iFace.SetWithStruct(false, resource)
}
