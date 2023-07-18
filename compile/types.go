package compile

import (
	"github.com/taubyte/config-compiler/indexer"
	projectSchema "github.com/taubyte/go-project-schema/project"
)

type compileObject struct {
	Get     func(string) (local []string, global []string)
	Compile magicFunc
	Indexer indexerFunc
}

type indexerFunc func(
	ctx *indexer.IndexContext,
	project projectSchema.Project,
	urlIndex map[string]interface{},
) error

type magicFunc func(
	name,
	app string,
	p projectSchema.Project,
) (
	_id string,
	ReturnMap map[string]interface{},
	err error,
)

// TODO: Something along these lines
/* package compile

import (
	projectSchema "github.com/taubyte/go-project-schema/project"
)

type compileObject struct {
	Get     func(string) (local []string, global []string)
	Compile magicFunc
	Indexer indexerFunc
}

type app struct {
	id   string
	name string
}

type appContext struct {
	compiler compiler
	app      *app
}

type indexable interface {
	Index(f indexerFunc, appObject map[string]interface{}) error
}

type subCompiler interface {
	indexable
	// compilable
}

func (c compiler) Application(id, name string) indexable {
	return &appContext{
		compiler: c,
		app: &app{
			id:   id,
			name: name,
		},
	}
}

func (a *appContext) Index(f indexerFunc, appObject map[string]interface{}) error {
	return f(a.app.id, a.app.name, a.compiler.project, appObject, a.compiler.index)
}

func (c *compiler) Index(f indexerFunc, appObject map[string]interface{}) error {
	return f("", "", c.project, appObject, c.index)
}

type indexerFunc func(appID string, appName string, project projectSchema.Project, obj map[string]interface{}, urlIndex map[string]interface{}) error

type magicFunc func(name string, app string, p projectSchema.Project) (_id string, ReturnMap map[string]interface{}, err error)
*/
