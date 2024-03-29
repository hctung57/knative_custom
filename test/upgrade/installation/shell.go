/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package installation

import (
	"testing"

	pkgupgrade "knative.dev/pkg/test/upgrade"
	"knative.dev/pkg/test/upgrade/shell"
)

// Head installs Knative Serving from the HEAD of the default branch.
func Head() pkgupgrade.Operation {
	return install("ServingHead", "install_head_reuse_ingress")
}

// LatestRelease installs Knative Serving from the latest stable release.
func LatestRelease() pkgupgrade.Operation {
	return install("ServingLatestRelease", "install_latest_release")
}

func install(installName, shellFunc string) pkgupgrade.Operation {
	return pkgupgrade.NewOperation(installName, func(c pkgupgrade.Context) {
		c.Log.Info("Running shell function: ", shellFunc)
		if err := callShellFunction(shellFunc, c.T); err != nil {
			c.T.Error(err)
		}
	})
}

func callShellFunction(funcName string, t *testing.T) error {
	loc, err := shell.NewProjectLocation("../../..")
	if err != nil {
		return err
	}
	exec := shell.NewExecutor(t, loc)
	fn := shell.Function{
		Script: shell.Script{
			Label:      funcName,
			ScriptPath: "test/e2e-common.sh",
		},
		FunctionName: funcName,
	}
	return exec.RunFunction(fn)
}
