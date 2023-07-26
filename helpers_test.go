package compiler

import (
	commonTest "github.com/taubyte/dreamland/helpers"
)

var (
	fakeMeta = commonTest.ConfigRepo.HookInfo
)

func init() {
	fakeMeta.Repository.Provider = "github"
	fakeMeta.Repository.Branch = "master"
}
