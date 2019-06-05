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
	"github.com/docker/dhe-deploy/gocode/dtr/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/gocode/dtr/ipc/settings/drivers/kv"
	"github.com/docker/dhe-deploy/gocode/dtr/shared/dtrutil/kvutil"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var replicaIDToDelete string

// removeReplicaCmd represents the removeReplica command
var removeReplicaCmd = &cobra.Command{
	Use:   "removeReplica",
	Short: "Command to remove a replica from DTR internal configuration",
	Long: `This command will remove a replica from the DTR internal configuration.

/!\ Please do a backup before using it !
/!\ To use only if asked by Docker Support !

You will need to make other actions to bring back the database (rethinkops, etc).`,
	Run: func(cmd *cobra.Command, args []string) {
		if replicaIDToDelete == "" {
			cmd.Help()
			os.Exit(1)
		}
		kvStore, err := kvutil.NewDefaultKVStore(replicaID)
		if err != nil {
			log.Fatal("Failed to get kv store: %s\n", err)
		}
		settingsStore := sanitizers.Wrap(kv.NewKVSettingsStore(kvStore))
		haConfig, err := settingsStore.HAConfig()
		if err != nil {
			log.Fatal("Failed to get ha config %s\n", err)
		}
		delete(haConfig.ReplicaConfig, replicaIDToDelete)
		log.Println("Removing replica " + replicaIDToDelete)
		err = settingsStore.SetHAConfig(haConfig)
		if err != nil {
			log.Fatal("Failed to set ha config %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeReplicaCmd)
	removeReplicaCmd.Flags().StringVar(&replicaIDToDelete, "replica-id-to-remove", "", "Replica-id to remove")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeReplicaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeReplicaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
