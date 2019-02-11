package main

import (
	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

func buildVersionCmd(f *filme.Filme) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the application version",
		Run: func(cmd *cobra.Command, args []string) {
			f.PrintVersion()
		},
	}
}
