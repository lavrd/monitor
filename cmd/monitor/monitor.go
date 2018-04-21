package main

import (
	"github.com/spacelavr/monitor/pkg/monitor"
	"github.com/spacelavr/pandora/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	verbose                      bool
	updContainersInterval        int
	updCOntainersMetricsInterval int
	changesFlushInterval         int
	port                         int

	// CLI main command
	CLI = &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				err := cmd.Help()
				if err != nil {
					log.Fatal(err)
				}

				return
			}

			log.SetVerbose(verbose)
		},
		Run: func(cmd *cobra.Command, args []string) {
			monitor.Daemon()
		},
	}
)

func init() {
	CLI.Flags().IntVarP(&port, "port", "p", 4222, "set api port")
	CLI.Flags().IntVarP(&updContainersInterval, "CInterval", "c", 3, "set update containers interval")
	CLI.Flags().IntVarP(&updCOntainersMetricsInterval, "CMInterval", "m", 1, "set update containers metrics interval")
	CLI.Flags().IntVarP(&changesFlushInterval, "FInterval", "f", 1, "set changes flush interval")
	CLI.Flags().BoolVarP(&verbose, "verbose", "v", false, "set verbose output")

	err := viper.BindPFlag("port", CLI.Flags().Lookup("port"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("CInterval", CLI.Flags().Lookup("CInterval"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("CMInterval", CLI.Flags().Lookup("CMInterval"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("FInterval", CLI.Flags().Lookup("FInterval"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("verbose", CLI.Flags().Lookup("verbose"))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := CLI.Execute(); err != nil {
		log.Fatal(err)
	}
}
