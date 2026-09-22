package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/syncgo/syncgo/internal/processor"
	"github.com/syncgo/syncgo/pkg/config"
	_ "github.com/syncgo/syncgo/pkg/logger"
	"github.com/syncgo/syncgo/pkg/metrics"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "version":
		err = runVersion(os.Args[2:])
	case "prepare":
		err = runPrepare(os.Args[2:])
	case "run":
		err = run(os.Args[2:])
	case "-h", "--help", "help":
		printUsage()
		return
	default:
		slog.Error("unknown command", slog.String("cmd", os.Args[1]))
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		slog.Error("command failed", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func runVersion(args []string) error {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)

	if err := fs.Parse(args); err != nil {
		return err
	}

	fmt.Printf("syncgo %s (commit %s, built %s)\n", version, commit, date)
	return nil
}

func runPrepare(args []string) error {
	fmt.Printf("preparing slots: %s", args)
	return nil
}

func run(args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfgPath, err := parseConfigFlag(args)
	if err != nil {
		return err
	}

	cfg, err := config.LoadFromYAML(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	monitoring := metrics.New(nil)

	p, err := processor.New(ctx, cfg, monitoring)
	if err != nil {
		return fmt.Errorf("failed initialize processor: %w", err)
	}

	if err := p.Run(ctx); err != nil {
		return fmt.Errorf("processor run failed: %w", err)
	}

	slog.Info("syncgo stopped successfully")
	return nil
}

func parseConfigFlag(args []string) (string, error) {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)

	var configPath string

	fs.StringVar(&configPath, "config", "", "Path to configuration file (YAML)")
	fs.StringVar(&configPath, "cfg", "", "Path to configuration file (YAML) (shorthand for --config)")

	if err := fs.Parse(args); err != nil {
		return "", err
	}
	return configPath, nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `syncgo is binary with few commands
	Usage:
	  syncgo version			prints version
	  syncgo prepare [flags]	run prepare commands for postgres (creates publication, slot, etc)
	  syncgo run [flags]		run sync process

	Use "syncgo <command> -h" for flags on specific command.`+"\n\n")
}
