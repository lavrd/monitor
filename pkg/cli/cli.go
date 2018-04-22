package cli

//
// import (
// 	"os"
//
// 	"github.com/spacelavr/monitor/pkg/cli/cmd"
// 	"github.com/spacelavr/monitor/pkg/context"
// 	"github.com/spacelavr/monitor/pkg/logger"
// 	"github.com/urfave/cli"
// )
//
// // Run start client
// func Run() {
// 	const (
// 		// flags
// 		verbose = "v"
// 		addr    = "a"
//
// 		// commands
// 		stopped  = "stopped"
// 		launched = "launched"
// 		metrics  = "metrics"
// 		logs     = "logs"
// 		status   = "status"
// 	)
//
// 	cli.VersionFlag = cli.BoolFlag{Name: "version, V", Usage: "print the version"}
//
// 	app := cli.NewApp()
// 	app.Name = "monitor"
// 	app.Usage = "Docker load monitor"
// 	app.Version = "0.1.0"
//
// 	app.Flags = []cli.Flag{
// 		cli.BoolFlag{
// 			Name:  verbose + ", verbose",
// 			Usage: "enable verbose output",
// 		},
// 		cli.StringFlag{
// 			Name:  addr + ", addr",
// 			Value: "http://localhost:4222",
// 			Usage: "set API address",
// 		},
// 	}
//
// 	app.Commands = []cli.Command{
// 		// command for view stopped containers
// 		{
// 			Name:  stopped,
// 			Usage: "view stopped containers",
// 			Action: func(c *cli.Context) error {
// 				cmd.StoppedContainersCmd()
// 				return nil
// 			},
// 		},
// 		// command for view launched containers
// 		{
// 			Name:  launched,
// 			Usage: "view launched containers",
// 			Action: func(c *cli.Context) error {
// 				cmd.LaunchedContainersCmd()
// 				return nil
// 			},
// 		},
// 		// command for view container logs
// 		{
// 			Name:  logs,
// 			Usage: "view container logs",
// 			Action: func(c *cli.Context) error {
// 				cmd.ContainerLogsCmd(c.Args().First())
// 				return nil
// 			},
// 		},
// 		// command for view container(s) metrics
// 		{
// 			Name:  metrics,
// 			Usage: "view containers metrics",
// 			Action: func(c *cli.Context) error {
// 				cmd.ContainersMetricsCmd(c.Args())
// 				return nil
// 			},
// 		},
// 		// command for view API status
// 		{
// 			Name:  status,
// 			Usage: "view API status",
// 			Action: func(c *cli.Context) error {
// 				cmd.APIStatusCmd()
// 				return nil
// 			},
// 		},
// 	}
//
// 	app.Before = func(c *cli.Context) error {
// 		// set verbose mode if use flag "v"
// 		context.Get().Verbose = c.Bool(verbose)
// 		// set API address
// 		context.Get().Address = c.String(addr)
// 		return nil
// 	}
//
// 	// for show help msg if use incorrect cmd
// 	app.Action = func(c *cli.Context) error {
// 		return cli.ShowAppHelp(c)
// 	}
//
// 	err := app.Run(os.Args)
// 	if err != nil {
// 		logger.Panic(err)
// 	}
// }
