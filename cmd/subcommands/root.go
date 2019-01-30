package subcommands

import (
	log "github.com/sirupsen/logrus"

	"github.com/florinutz/filme/cmd/subcommands/coll33tx"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var CmdRoot = "filme"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   CmdRoot,
	Short: "Parses torrent sites",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	RootCmd.AddCommand(coll33tx.L33txRootCmd)
	if err := RootCmd.Execute(); err != nil {
		log.Debug(err)
	}
}

func init() {
	log.SetLevel(log.InfoLevel)
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.filme.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.WithError(err).Fatal()
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".filme")
	}
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithField("config_file", viper.ConfigFileUsed()).Debug("using config file")
	}
}
