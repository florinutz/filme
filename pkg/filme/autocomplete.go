package filme

import (
	"io"

	"github.com/florinutz/filme/pkg/config"

	"github.com/spf13/cobra"
)

func (f *Filme) AutoComplete(rootCmd *cobra.Command, shellName string) error {
	var shell config.Shell
	shell, err := config.ParseShell(shellName)
	if err != nil {
		return err
	}
	completionFuncs := map[config.Shell]func(io.Writer) error{
		config.ShellBash: rootCmd.GenBashCompletion,
		config.ShellZsh:  rootCmd.GenZshCompletion,
	}

	return completionFuncs[shell](f.Out)
}
