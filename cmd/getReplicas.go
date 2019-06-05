// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"github.com/docker/dhe-deploy/gocode/dtr/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/gocode/dtr/ipc/settings/drivers/kv"
	"github.com/docker/dhe-deploy/gocode/dtr/shared/dtrutil/kvutil"
	"github.com/spf13/cobra"
)

// getReplicasCmd represents the getReplicas command
var getReplicasCmd = &cobra.Command{
	Use:   "getReplicas",
	Short: "Command to retrieve all replicas",
	Long: `This command will retrieve all replicas visible inside the database of DTR.

It can help to see the id before doing other commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		kvStore, err := kvutil.NewDefaultKVStore(replicaID)
		if err != nil {
			log.Fatal("Failed to get kv store: %s\n", err)
		}
		settingsStore := sanitizers.Wrap(kv.NewKVSettingsStore(kvStore))
		haConfig, err := settingsStore.HAConfig()
		if err != nil {
			log.Fatal("Failed to get ha config %s\n", err)
		}
		i := 1
		for rID, element := range haConfig.ReplicaConfig {
			fmt.Printf("Replica %v: ReplicaId: %s | NodeId: %s | Version: %s\n", i, rID, element.Node, element.Version)
			i = i + 1
		}
	},
}

func init() {
	rootCmd.AddCommand(getReplicasCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getReplicasCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getReplicasCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
