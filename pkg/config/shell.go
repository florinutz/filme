package config

import "fmt"

type Shell string

const (
	ShellBash Shell = "bash"
	ShellZsh  Shell = "zsh"
)

func ParseShell(shellName string) (Shell, error) {
	shell := Shell(shellName)
	switch shell {
	case ShellBash, ShellZsh:
		return shell, nil
	default:
		return "", fmt.Errorf("invalid shell %q", shellName)
	}
}
