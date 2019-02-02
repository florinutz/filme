package subcommands

import (
	"os"

	log "github.com/sirupsen/logrus"

	"io"

	"github.com/spf13/cobra"
)

const (
	Bash = "bash"
	Zsh  = "zsh"
)

var (
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
			if f, ok := completionFuncs[args[0]]; ok {
				if err := f(os.Stdout); err != nil {
					log.WithError(err).Fatal("error while generating shell completion")
				}
			} else {
				log.Fatal("invalid shell for completion")
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(completionCmd)
}
