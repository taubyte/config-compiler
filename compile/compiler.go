package compile

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"bitbucket.org/taubyte/mycelium/logfile"
	"github.com/spf13/afero"
	"github.com/taubyte/config-compiler/decompile"
	"github.com/taubyte/config-compiler/ifaces"
	"github.com/taubyte/config-compiler/indexer"
	"github.com/taubyte/go-interfaces/services/patrick"
	projectSchema "github.com/taubyte/go-project-schema/project"
)

type compiler struct {
	ctx    indexer.IndexContext
	index  map[string]interface{}
	post   []func() (err error)
	log    *logfile.File
	config *Config

	dev bool
}

type result struct {
	data map[string]interface{}
}

type Option func(*compiler) error

type Config struct {
	Branch       string
	Commit       string
	Provider     string
	RepositoryId string
	project      projectSchema.Project
}

func CompilerConfig(project projectSchema.Project, meta patrick.Meta) (*Config, error) {
	if project == nil {
		return nil, errors.New("Project is nil")
	}

	if meta.Repository.ID == 0 {
		return nil, errors.New("RepoId is nil")
	}

	for idx, s := range map[string]string{"branch": meta.Repository.Branch, "commit": meta.HeadCommit.ID, "provider": meta.Repository.Provider} {
		if len(s) < 1 {
			return nil, fmt.Errorf("Metadata %s is empty", idx)
		}
	}

	return &Config{project: project, Branch: meta.Repository.Branch, Commit: meta.HeadCommit.ID, Provider: meta.Repository.Provider, RepositoryId: strconv.Itoa(meta.Repository.ID)}, nil
}

func Dev() Option {
	return func(c *compiler) error {
		c.dev = true
		return nil
	}
}

func DVKey(publicKey []byte) Option {
	return func(c *compiler) error {
		c.ctx.DVPublicKey = publicKey
		return nil
	}
}

func New(config *Config, options ...Option) (ifaces.Compiler, error) {
	log, err := logfile.New()
	if err != nil {
		return nil, err
	}

	compiler := &compiler{
		post:   make([]func() (err error), 0),
		log:    log,
		config: config,
	}

	for _, opt := range options {
		err := opt(compiler)
		if err != nil {
			compiler.Close()

			return nil, err
		}
	}

	return compiler, nil
}

func (c *compiler) Object() map[string]interface{} {
	return c.ctx.Obj
}

func (c *compiler) Indexes() map[string]interface{} {
	return c.index
}

func (c *compiler) Commit() string {
	return c.config.Commit
}

func (c *compiler) Branch() string {
	return c.config.Branch
}

func (c *compiler) Close() error {
	return c.log.Close()
}

func (c *compiler) Logs() io.ReadSeeker {
	return c.log
}

func (c *compiler) Load(object map[string]interface{}) error {
	decompiler, err := decompile.New(afero.NewMemMapFs(), object)
	if err != nil {
		return err
	}

	decompiler.Build()
	prj, err := decompiler.Build()
	if err != nil {
		fmt.Fprintf(c.log, "[LOAD|OBJECT] failed with %s\n", err)
		return err
	}

	c.config.project = prj

	return nil
}

func (c *compiler) Build() error {
	var err error
	tee := func(prefix string, err error) error {
		fmt.Fprintln(c.log, prefix, " ", err.Error())
		return err
	}

	if c.config.project == nil {
		return tee("[Build]", errors.New("No project found"))
	}

	getter := c.config.project.Get()

	c.ctx.Branch = c.config.Branch
	c.ctx.Commit = c.config.Commit
	c.ctx.ProjectId = getter.Id()
	c.index = make(map[string]interface{})
	c.ctx.Obj = map[string]interface{}{
		"id":          getter.Id(),
		"name":        getter.Name(),
		"description": getter.Description(),
		"email":       getter.Email(),
	}
	c.ctx.Dev = c.dev

	for _type, iFace := range compilationGroup(c.config.project) {
		_, global := iFace.Get("")
		if len(global) > 0 {
			c.ctx.Obj[_type], err = c.magic(global, "", iFace.Compile)
			if err != nil {
				return tee("[Build] project="+getter.Id(), err)
			}

			if iFace.Indexer != nil {
				err = c.indexer(&c.ctx, iFace.Indexer)
				if err != nil {
					return tee("[Build] project="+getter.Id(), err)
				}
			}
		}
	}

	// Get all applications and their resources
	apps := getter.Applications()
	if len(apps) > 0 {
		applications := make(map[string]interface{})
		for _, app := range apps {
			_id, appObject, err := c.application(app)
			if err != nil {
				return tee("[Build] app= "+app+" project= "+getter.Id(), err)
			}

			applications[_id] = appObject
		}
		if len(applications) > 0 {
			c.ctx.Obj["applications"] = applications
		}
	}

	for _, post := range c.post {
		err = post()
		if err != nil {
			return tee("[POST] project = "+getter.Id(), err)
		}
	}

	return nil
}
