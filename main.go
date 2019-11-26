package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/florinutz/filme/commands"
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
			f.Log.Level = f.DebugLevel.Level
		}}

	defaultDebugLevel := logrus.ErrorLevel
	_ = f.DebugLevel.Set(defaultDebugLevel.String())
	cmd.PersistentFlags().VarP(&f.DebugLevel, "debug-level", "", fmt.Sprintf("one of: %s",
		strings.Join(value.GetAllLevels(), ", ")))

	crawlCmd := commands.BuildCrawlCmd(f)
	crawlCmd.AddCommand(
		commands.Build1337xDetailPageCmd(f),
		commands.BuildImdbDetailPageCmd(f),
		// commands.BuildGoogleCmd(f),
	)

	cmd.AddCommand(
		crawlCmd,
		commands.BuildSearchCmd(f),
		commands.BuildCompletionCmd(f),
		commands.BuildVersionCmd(f),
	)

	return cmd
}
