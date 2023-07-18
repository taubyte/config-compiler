package decompile

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/taubyte/go-project-schema/project"
	structureSpec "github.com/taubyte/go-specs/structure"
	"github.com/taubyte/utils/id"
)

type buildContext struct {
	projectId string
	dir       string
	project   project.Project
}

func generateId(_id string) string {
	if len(_id) > 0 {
		return _id
	} else {
		return id.Generate(_id)
	}
}

// Takes a slice of structureSpec Structures and converts them into the project
func MockBuild(projectId string, dir string, ifaces ...interface{}) (project.Project, error) {
	ctx := &buildContext{projectId: projectId, dir: dir}

	err := ctx.newProject()
	if err != nil {
		return nil, err
	}

	if err := ctx.newStructs(ifaces...); err != nil {
		return nil, err
	}

	return ctx.project, nil
}

func (ctx *buildContext) newProject() (err error) {
	if ctx.dir == "" {
		ctx.dir, err = ioutil.TempDir(os.TempDir(), "project-*")
		if err != nil {
			return
		}
	}

	err = os.MkdirAll(ctx.dir, 0777)
	if err != nil {
		return fmt.Errorf("Creating tx.dir %s failed with: %v", ctx.dir, err)
	}

	ctx.project, err = project.Open(project.SystemFS(ctx.dir))
	if err != nil {
		return fmt.Errorf("project.Open failed with: %v", err)
	}

	err = ctx.project.Set(
		true,
		project.Id(ctx.projectId),
		project.Name("builtProject"),
	)
	if err != nil {
		return fmt.Errorf("p.set failed with: %v", err)
	}

	return
}

func (ctx *buildContext) newStructs(ifaces ...interface{}) (err error) {
	for _, iface := range ifaces {
		if err = ctx.newStruct(iface); err != nil {
			return
		}
	}
	return
}

func (ctx *buildContext) newStruct(iface interface{}) (err error) {
	switch iface.(type) {
	case *structureSpec.Function:
		_iface := iface.(*structureSpec.Function)
		return function(ctx.project, generateId(_iface.Id), iface, "")
	case *structureSpec.Messaging:
		_iface := iface.(*structureSpec.Messaging)
		return messaging(ctx.project, generateId(_iface.Id), iface, "")
	case *structureSpec.Domain:
		_iface := iface.(*structureSpec.Domain)
		return domain(ctx.project, generateId(_iface.Id), iface, "")
	case *structureSpec.Database:
		_iface := iface.(*structureSpec.Database)
		return database(ctx.project, generateId(_iface.Id), iface, "")
	case *structureSpec.Storage:
		_iface := iface.(*structureSpec.Storage)
		return storage(ctx.project, generateId(_iface.Id), iface, "")
	case *structureSpec.Service:
		_iface := iface.(*structureSpec.Service)
		return service(ctx.project, generateId(_iface.Id), iface, "")
	case *structureSpec.Library:
		_iface := iface.(*structureSpec.Library)
		return library(ctx.project, generateId(_iface.Id), iface, "")
	case *structureSpec.SmartOp:
		_iface := iface.(*structureSpec.SmartOp)
		return smartop(ctx.project, generateId(_iface.Id), iface, "")
	case *structureSpec.Website:
		_iface := iface.(*structureSpec.Website)
		return website(ctx.project, generateId(_iface.Id), iface, "")
	case []interface{}:
		for _, _iface := range iface.([]interface{}) {
			err = ctx.newStruct(_iface)
			if err != nil {
				return
			}
		}
	default:
		err = fmt.Errorf("struct `%T` not yet supported", iface)
	}
	return
}
