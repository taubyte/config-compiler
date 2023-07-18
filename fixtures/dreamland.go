package fixtures

import (
	dreamlandRegistry "github.com/taubyte/dreamland/core/registry"
)

func init() {
	dreamlandRegistry.Fixture("fakeProject", fakeProject)
	dreamlandRegistry.Fixture("injectProject", injectProject)
}
