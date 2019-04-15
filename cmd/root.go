package cmd

import (
	"fmt"
	"os"

	"github.com/docker/dhe-deploy/gocode/dtr/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/gocode/dtr/ipc/settings/drivers/kv"
	"github.com/docker/dhe-deploy/gocode/dtr/shared/dtrutil/kvutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var cfgFile string
var replicaID string
var httpPort int
var httpsPort int
var rethinkCache int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dtr-global-change",
	Short: "Command to change configuration of all replicas.",
	Long: `This command will change the configuration of all replicas.

/!\ Please do a backup before using it !

You will need to do a "dtr reconfigure" to apply globally these changes.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if httpPort == 0 && httpsPort == 0 && rethinkCache == -1 {
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
		for rID, element := range haConfig.ReplicaConfig {
			if httpPort != 0 && element.HTTPPort != uint16(httpPort) {
				element.HTTPPort = uint16(httpPort)
				log.Printf("Changing HTTP port to %v\n", httpPort)
			}
			if httpsPort != 0 && element.HTTPSPort != uint16(httpsPort) {
				element.HTTPSPort = uint16(httpsPort)
				log.Printf("Changing HTTPS port to %v\n", httpsPort)
			}
			if rethinkCache != -1 && element.RethinkdbCacheMB != rethinkCache {
				element.RethinkdbCacheMB = rethinkCache
				log.Printf("Changing rethinkdb cache to %vmb\n", rethinkCache)
			}
			haConfig.ReplicaConfig[rID] = element
		}
		err = settingsStore.SetHAConfig(haConfig)
		if err != nil {
			log.Fatal("Failed to set ha config %s\n", err)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	viper.BindEnv("dtr_replica_id")

	if viper.IsSet("dtr_replica_id") {
		rootCmd.Flags().StringVar(&replicaID, "replica-id", viper.Get("dtr_replica_id").(string), "Replica-id to connect")
	} else {
		rootCmd.Flags().StringVar(&replicaID, "replica-id", "", "Replica-id to connect")
		rootCmd.MarkFlagRequired("replica-id")
	}
	rootCmd.Flags().IntVarP(&httpPort, "http-port", "", 0, "Http port that will use all replicas")
	rootCmd.Flags().IntVarP(&httpsPort, "https-port", "", 0, "Https port that will use all replicas")
	rootCmd.Flags().IntVarP(&rethinkCache, "rethinkdb-cache-mb", "", -1, "Max rethinkdb memory cache that will use all replicas | 0 = auto")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

}
