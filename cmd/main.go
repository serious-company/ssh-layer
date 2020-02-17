package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "ssh-layer",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(2)
	},
}

type Config struct {
}

func main() {
	var c Config
	cobra.OnInitialize(func() {
		initConfig(&c)
	})
	if err := commandRoot(&c).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}

func commandRoot(c *Config) *cobra.Command {
	//
	return rootCmd
}

func initConfig(c *Config) {
	v := viper.New()
	err := v.Unmarshal(c)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
