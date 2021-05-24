package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/florinutz/filme/pkg/config/value"
	"github.com/florinutz/filme/pkg/filme"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	if err := buildRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func buildRootCommand() *cobra.Command {
	f := filme.New()

	cmd := &cobra.Command{
		Use:   "filme",
		Short: "film torrenting cli utility",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			f.Out = cmd.OutOrStdout()
			f.Err = cmd.OutOrStderr()
			f.Log.Out = f.Out
			f.Log.Level = f.DebugLevel.Level
			f.Log.ReportCaller = f.DebugReportCaller
		}}

	defaultDebugLevel := logrus.PanicLevel
	_ = f.DebugLevel.Set(defaultDebugLevel.String())
	cmd.PersistentFlags().Var(&f.DebugLevel, "debug-level", fmt.Sprintf("one of: %s",
		strings.Join(value.GetAllLevels(), ", ")))

	cmd.PersistentFlags().BoolVar(&f.DebugReportCaller, "debug-report-caller", false, "show debug callers")

	var crawlGroupCmd, rarbgGroupCmd *cobra.Command

	{
		crawlGroupCmd = &cobra.Command{
			Use:   "crawl",
			Short: "crawl commands",
		}
		crawlGroupCmd.AddCommand(
			cmd.BuildImdbDetailPageCmd(f),
			// commands.BuildGoogleCmd(f),
		)
		rarbgGroupCmd = &cobra.Command{
			Use:   "rarbg",
			Short: "rarbg commands",
		}
	}

	cmd.AddCommand(
		crawlGroupCmd,
		rarbgGroupCmd,
		cmd.BuildSearchCmd(f),
		cmd.Build1337xDetailPageCmd(f),
		cmd.BuildServeCmd(f),
		cmd.BuildCompletionCmd(f),
		cmd.BuildVersionCmd(f),
		cmd.BuildScreenshotCmd(f),
	)

	return cmd
}
