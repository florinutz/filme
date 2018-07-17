package subcommands

import (
	"os"

	"io"

	"strings"

	"github.com/spf13/cobra"
)

const (
	Bash = "bash"
	Zsh  = "zsh"
)

var (
	shell string

	completionFuncs = map[string]func(io.Writer) error{
		Bash: RootCmd.GenBashCompletion,
		Zsh:  RootCmd.GenZshCompletion,
	}

	completionCmd = &cobra.Command{
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
		ValidArgs: []string{Bash, Zsh},
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			shell = args[0]
			completionFuncs[shell](os.Stdout)
		},
	}
)

func init() {
	RootCmd.AddCommand(completionCmd)
}

// shellIsValid makes sure the shell is among zsh or bash, case insensitive
func shellIsValid(shell string) bool {
	return strings.EqualFold(shell, Bash) || strings.EqualFold(shell, Zsh)
}
