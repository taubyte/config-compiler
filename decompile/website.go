package decompile

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	projectLib "github.com/taubyte/go-project-schema/project"
	lib "github.com/taubyte/go-project-schema/website"
	structureSpec "github.com/taubyte/go-specs/structure"
)

func website(project projectLib.Project, _id string, obj interface{}, appName string) error {
	resource := &structureSpec.Website{}
	mapstructure.Decode(obj, resource)

	iFace, err := project.Website(resource.Name, appName)
	if err != nil {
		return fmt.Errorf("Open website `%s/%s` failed: %s", appName, resource.Name, err)
	}

	resource.SetId(_id)
	return iFace.SetWithStruct(false, resource)
}

func website_clean(project projectLib.Project, name, app string) (err error) {
	website, err := project.Website(name, app)
	if err != nil {
		return fmt.Errorf("Couldn't open website `%s/%s` to clean: %v", app, name, err)
	}

	old_domains := website.Get().Domains()
	new_domains, err := cleanDoms(project, old_domains, app)
	if err != nil {
		return fmt.Errorf("Clean domains of website `%s/%s` failed with: %v", app, name, err)
	}

	if len(new_domains) > 0 {
		err = website.Set(false, lib.Domains(new_domains))
		if err != nil {
			return fmt.Errorf("Set domains of website `%s/%s` failed with: %v", app, name, err)
		}
	}

	return
}
