package fixtures

import (
	dreamlandRegistry "bitbucket.org/taubyte/dreamland/registry"
)

func init() {
	dreamlandRegistry.Fixture("fakeProject", fakeProject)
	dreamlandRegistry.Fixture("injectProject", injectProject)
}
