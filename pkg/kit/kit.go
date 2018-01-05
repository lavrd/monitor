package kit

import (
	"net/http"
	"os"
	"time"

	"github.com/spacelavr/dlm/pkg/context"
	m "github.com/spacelavr/dlm/pkg/kit/metrics"
	"github.com/spacelavr/dlm/pkg/kit/router"
	"github.com/spacelavr/dlm/pkg/logger"
	"github.com/urfave/cli"
)

// Run start kit
func Run() {
	const (
		// flags
		verbose       = "v"
		updContsInt   = "uci"
		updContMetInt = "ucmi"
		chFLushInt    = "f"
		port          = "p"
	)

	cli.VersionFlag = cli.BoolFlag{Name: "version, V", Usage: "print the version"}

	app := cli.NewApp()
	app.Name = "dlm"
	app.Usage = "Docker load monitor"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  port + ", port",
			Value: "4222",
			Usage: "set API port",
		},
		cli.IntFlag{
			Name:  updContsInt + ", upd-conts-int",
			Value: 3,
			Usage: "set update containers interval",
		},
		cli.IntFlag{
			Name:  updContMetInt + ", upd-cont-met-int",
			Value: 1,
			Usage: "set update container metrics interval",
		},
		cli.IntFlag{
			Name:  chFLushInt + ", ch-flush-int",
			Value: 10,
			Usage: "set changes flush interval",
		},
		cli.BoolFlag{
			Name:  verbose + ", verbose",
			Usage: "enable verbose output",
		},
	}

	app.Action = func(c *cli.Context) error {
		// if args > 0 -> error
		if c.NArg() > 0 {
			err := cli.ShowAppHelp(c)
			if err != nil {
				return err
			}
			return nil
		}

		// set verbose mode if use flag "v"
		context.Get().Verbose = c.Bool(verbose)
		// set api address
		context.Get().Address = ":" + c.String(port)

		// set intervals
		m.Get().SetUContsInterval(time.Duration(c.Int(updContsInt)) * time.Second)
		m.Get().SetUCMetricsInterval(time.Duration(c.Int(updContMetInt)) * time.Second)
		m.Get().SetChangesFlushInterval(time.Duration(c.Int(chFLushInt)) * time.Second)
		// start collect metrics
		go m.Get().Collect()

		// listen and serve
		fsrv := &http.Server{
			Handler: router.Router(),
			Addr:    context.Get().Address,
		}
		return fsrv.ListenAndServe()
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Panic(err)
	}
}
