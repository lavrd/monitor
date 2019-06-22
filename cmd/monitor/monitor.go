package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"monitor/pkg/monitor"
)

func main() {
	metricsInterval := flag.Duration("cmi", time.Second, "set metrics interval")
	containersInterval := flag.Duration("ci ", time.Second*5, "set container interval")
	port := flag.Int("port", 2000, "set api port")
	verbose := flag.Bool("v", false, "set verbose output")
	flag.Parse()

	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().
		Caller().
		Logger().
		Level(zerolog.InfoLevel)
	if *verbose {
		log.Logger = log.Level(zerolog.DebugLevel)
	}

	if err := monitor.Daemon(*metricsInterval, *containersInterval, *port); err != nil {
		log.Fatal().Err(err).Msg("daemon error")
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-interrupt
	log.Info().Msg("handle SIGINT, SIGTERM, SIGQUIT, SIGKILL")
}
