package compile

import (
	"errors"
	"fmt"

	tnsIface "github.com/taubyte/go-interfaces/services/tns"
	specsCommon "github.com/taubyte/go-specs/common"
	specs "github.com/taubyte/go-specs/methods"
)

func (c *compiler) Publish(tns tnsIface.Client) (err error) {
	if c.index == nil || c.ctx.Obj == nil {
		return errors.New("Build first")
	}

	err = tns.Push([]string{}, c.index)
	if err != nil {
		return fmt.Errorf("Publish index failed with: %w", err)
	}

	project := c.config.project.Get().Id()

	prefix := specs.ProjectPrefix(project, c.ctx.Branch, c.ctx.Commit)
	err = tns.Push(prefix.Slice(), c.ctx.Obj)
	if err != nil {
		return fmt.Errorf("Publish project failed with: %w", err)
	}

	//TODO: DO THIS CLEANER
	err = tns.Push(
		specsCommon.Current(
			project,
			c.ctx.Branch).Slice(),
		map[string]string{
			specsCommon.CurrentCommitPathVariable.String(): c.ctx.Commit,
		},
	)
	if err != nil {
		return fmt.Errorf("Publishing current commit for project `%s` on branch `%s` failed with: %w", project, c.ctx.Branch, err)
	}

	return
}
