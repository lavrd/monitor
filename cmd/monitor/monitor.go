package main

import (
	"monitor/pkg/monitor"
	"monitor/pkg/types"
	"monitor/pkg/utils/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	verbose    bool
	cInterval  int
	cmInterval int
	port       int

	// CLI main command
	CLI = &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetVerbose(verbose)
		},
		Run: func(cmd *cobra.Command, args []string) {
			monitor.Daemon()
		},
	}
)

func init() {
	CLI.Flags().IntVarP(&port, types.FPort, types.FSPort, 2000, "set http port")
	CLI.Flags().IntVarP(&cInterval, types.FCInterval, types.FSCInterval, 3, "set update containers interval")
	CLI.Flags().IntVarP(&cmInterval, types.FCMInterval, types.FSCMInterval, 1, "set update containers metrics interval")
	CLI.Flags().BoolVarP(&verbose, types.FVerbose, types.FSVerbose, false, "set verbose output")

	err := viper.BindPFlag(types.FPort, CLI.Flags().Lookup(types.FPort))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag(types.FCInterval, CLI.Flags().Lookup(types.FCInterval))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag(types.FCMInterval, CLI.Flags().Lookup(types.FCMInterval))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := CLI.Execute(); err != nil {
		log.Fatal(err)
	}
}
