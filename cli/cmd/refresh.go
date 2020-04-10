/*
Copyright 2020 Cortex Labs, Inc.

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

package cmd

import (
	"github.com/cortexlabs/cortex/cli/cluster"
	"github.com/cortexlabs/cortex/cli/local"
	"github.com/cortexlabs/cortex/pkg/lib/exit"
	"github.com/cortexlabs/cortex/pkg/lib/print"
	"github.com/cortexlabs/cortex/pkg/lib/telemetry"
	"github.com/cortexlabs/cortex/pkg/operator/schema"
	"github.com/cortexlabs/cortex/pkg/types"
	"github.com/spf13/cobra"
)

var _flagRefreshForce bool

func init() {
	_refreshCmd.PersistentFlags().BoolVarP(&_flagRefreshForce, "force", "f", false, "override the in-progress api update")
	addEnvFlag(_refreshCmd, types.LocalProviderType.String())
}

var _refreshCmd = &cobra.Command{
	Use:   "refresh API_NAME",
	Short: "restart all replicas for an api (witout downtime)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		telemetry.Event("cli.refresh")

		env := MustReadOrConfigureEnv(_flagEnv)
		var refreshResponse schema.RefreshResponse
		var err error
		if env.Provider == types.AWSProviderType {
			refreshResponse, err = cluster.Refresh(MustGetOperatorConfig(env.Name), args[0], _flagRefreshForce)
			if err != nil {
				exit.Error(err)
			}
		} else {
			// TODO show that flags are being ignored?
			refreshResponse, err = local.Refresh(args[0])
			if err != nil {
				exit.Error(err)
			}
		}

		print.BoldFirstLine(refreshResponse.Message)
	},
}
