package main

import (
	"os"

	"github.com/florinutz/filme/pkg/filme"
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
		Short: "movies utility",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			f.Out = cmd.OutOrStdout()
			f.Err = cmd.OutOrStderr()
		}}

	cmd.PersistentFlags().BoolVar(&f.Debug, "debug", false, "Enable debug logging")

	cmd.AddCommand(
		buildListCmd(f),
		buildCompletionCmd(f),
		buildVersionCmd(f),
	)

	return cmd
}
