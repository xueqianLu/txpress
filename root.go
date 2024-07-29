package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xueqianLu/txpress/chains"
	"github.com/xueqianLu/txpress/config"
	"github.com/xueqianLu/txpress/types"
	"github.com/xueqianLu/txpress/workflow"
	"os"
	"runtime/pprof"
	"time"
)

var (
	cpuProfile   bool
	configpath   string
	startCommand bool
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&startCommand, "start", false, "Start after initializing the account")
	rootCmd.PersistentFlags().BoolVar(&cpuProfile, "cpuProfile", false, "Statistics cpu profile")
	rootCmd.PersistentFlags().StringVar(&configpath, "config", "app.json", "config file path")

	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Program execute error: %s", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "txpress",
	Short: "Stress test tools",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("check start and ", "start is", startCommand)
		cfg, err := config.ParseConfig(configpath)
		if err != nil {
			os.Exit(1)
		}
		if startCommand {
			allchain := chains.NewChains(cfg)
			if len(allchain) == 0 {
				log.Error("have no chain to start")
				return
			}
			flow := workflow.NewWorkFlow(allchain, types.RunConfig{
				BaseCount: cfg.BaseCount,
				Batch:     cfg.Batch,
				Interval:  time.Duration(cfg.Interval) * time.Second,
			})

			if cpuProfile {
				f, err := os.Create("cpuprofile.log")
				if err != nil {
					log.Fatal(err)
				}
				err = pprof.StartCPUProfile(f)
				if err != nil {
					log.Error("Start cpu profile err:", err)
					return
				}
			}
			flow.Start()

			if cpuProfile {
				pprof.StopCPUProfile()
			}
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of txpress",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Version: ", Version)
		log.Info("Git Commit: ", Commit)
	},
}
