package cmd

import (
	"github.com/florinutz/filme/pkg/config"
	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

func BuildCompletionCmd(f *filme.Filme) *cobra.Command {
	return &cobra.Command{
		Use:   "completion <bash|zsh>",
		Short: "Generates shell completion scripts",
		Long: `
To load completion for bash run
	. <(filme completion bash)

To load for zsh run 
	. <(filme completion zsh)

To configure your shell to load completions for each session add this to your bash's .bashrc or zsh's .zshrc:
	. <(filme completion bash)

To configure your zsh shell to load completions for each session add this to your .zshrc:
	. <(filme completion zsh)

And then source them:
	source ~/.zshrc

or restart the terminal.
`,
		ValidArgs: []string{string(config.ShellBash), string(config.ShellZsh)},
		Args:      cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return f.AutoComplete(cmd.Root(), args[0])
		},
	}
}
